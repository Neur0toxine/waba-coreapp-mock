package main

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
)

var (
	validate       = validator.New()
	NotDigitsRegex = regexp.MustCompile("\\D+")
)

type Mock struct {
	ContactsSuccess bool `json:"contacts_success"`
	MessagesSuccess bool `json:"messages_success"`
}

type Server struct {
	g    *gin.Engine
	mock Mock
}

func NewServer() (s *Server) {
	s = &Server{g: gin.New()}
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
	if err := s.bindRequest(c, &mock); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	s.mock = mock
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

	log.Printf("Received new message: %#v\n", req)

	c.JSON(http.StatusOK, MessagesResponse{
		BaseResponse: s.baseResponseOk(),
		Messages: []IDModel{{
			ID: RandomString(27),
		}},
	})
}
