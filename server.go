package main

import (
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
)

var (
	validate       = validator.New()
	NotDigitsRegex = regexp.MustCompile("\\D+")
)

type Mock struct {
	ContactsSuccess bool              `json:"contacts_success"`
	MessagesSuccess bool              `json:"messages_success"`
	MessagesStatus  string            `json:"messages_success_status" validate:"oneof=sent read failed"`
	Webhook         string            `json:"webhook" validate:"url,startswith=http"`
	WebhookHeaders  map[string]string `json:"webhook_headers"`
}

type Server struct {
	g       *gin.Engine
	shooter *Shooter
	mock    Mock
}

func NewServer() (s *Server) {
	s = &Server{
		g: gin.New(),
		mock: Mock{
			ContactsSuccess: true,
			MessagesSuccess: true,
			MessagesStatus:  "sent",
			Webhook:         "",
			WebhookHeaders:  map[string]string{},
		},
	}
	s.updateShooter()
	s.g.GET("/mock", s.mockData)
	s.g.POST("/mock", s.updateMockData)
	api := s.g.Group("/v1")
	{
		api.POST("/contacts", s.contactsHandler)
		api.POST("/messages", s.messagesHandler)
	}
	return s
}

func (s *Server) Run(addr ...string) error {
	return s.g.Run(addr...)
}

func (s *Server) updateShooter() {
	if s.shooter == nil {
		s.shooter = NewShooter(s.mock.Webhook, s.mock.WebhookHeaders)
		return
	}
	s.shooter.Webhook = s.mock.Webhook
	s.shooter.Headers = s.mock.WebhookHeaders
}

func (s *Server) baseResponseOk() BaseResponse {
	return BaseResponse{
		Meta: &Metadata{
			Success:   true,
			APIStatus: "stable",
			Version:   "v2.31.5",
		},
	}
}

func (s *Server) bindRequest(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("error: %s\n", err)
		return err
	}
	if err := validate.Struct(req); err != nil {
		log.Printf("error: %s\n", err)
		return err
	}
	return nil
}

func (s *Server) mockData(c *gin.Context) {
	c.JSON(http.StatusOK, s.mock)
}

func (s *Server) updateMockData(c *gin.Context) {
	var mock Mock
	if err := c.ShouldBindJSON(&mock); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	current := s.mock
	current.ContactsSuccess = mock.ContactsSuccess
	current.MessagesSuccess = mock.MessagesSuccess

	if mock.MessagesStatus != "" {
		current.MessagesStatus = mock.MessagesStatus
	}

	if mock.Webhook != "" {
		current.Webhook = mock.Webhook
	}

	if mock.WebhookHeaders != nil {
		current.WebhookHeaders = mock.WebhookHeaders
	}

	if err := validate.Struct(current); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.mock = current
	s.updateShooter()
	c.JSON(http.StatusOK, s.mock)
}

func (s *Server) contactsHandler(c *gin.Context) {
	var req ContactsRequest
	if err := s.bindRequest(c, &req); err != nil || !s.mock.ContactsSuccess {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res := ContactsResponse{
		BaseResponse: s.baseResponseOk(),
		Contacts:     make([]Contact, len(req.Contacts)),
	}

	for i, contact := range req.Contacts {
		res.Contacts[i] = Contact{
			WaID:   NotDigitsRegex.ReplaceAllString(contact, ""),
			Input:  contact,
			Status: ContactStatusValid,
		}
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) messagesHandler(c *gin.Context) {
	var req Message
	if err := s.bindRequest(c, &req); err != nil || !s.mock.MessagesSuccess {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	messageID := RandomString(27)
	text := ""
	if req.Text != nil {
		text = req.Text.Body
	}

	log.Printf("Received new message: %#v\n", req)

	if s.mock.Webhook != "" {
		defer func(msgID, text string, to string) {
			go func(msgID, text string, to string) {
				time.Sleep(time.Millisecond * 500)

				code, err := s.shooter.SendStatus(InboundStatus{
					Type:        "message",
					ID:          messageID,
					RecipientID: to,
					Status:      s.mock.MessagesStatus,
				})
				if err != nil {
					log.Printf("error: %s\n", err)
					return
				}
				log.Printf("status webhook code: %d\n", code)

				if text == "reply" {
					code, err := s.shooter.SendText("Replying to the message", to)
					if err != nil {
						log.Printf("error: %s\n", err)
						return
					}
					log.Printf("reply webhook code: %d\n", code)
				}
			}(msgID, text, req.To)
		}(messageID, text, req.To)
	}

	c.JSON(http.StatusOK, MessagesResponse{
		BaseResponse: s.baseResponseOk(),
		Messages: []IDModel{{
			ID: messageID,
		}},
	})
}
