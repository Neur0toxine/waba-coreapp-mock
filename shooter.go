package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Shooter struct {
	Webhook string
	Headers map[string]string
}

func NewShooter(webhook string, headers map[string]string) *Shooter {
	return &Shooter{
		Webhook: webhook,
		Headers: headers,
	}
}

func (s *Shooter) makeRequest(v interface{}) (*http.Request, error) {
	wh, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, s.Webhook, bytes.NewReader(wh))
	if err != nil {
		return nil, err
	}

	for h, v := range s.Headers {
		req.Header.Set(h, v)
	}

	return req, nil
}

func (s *Shooter) SendStatus(status InboundStatus) (int, error) {
	req, err := s.makeRequest(InboundWebhook{
		Statuses: []InboundStatus{status},
	})
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func (s *Shooter) SendText(text, from string) (int, error) {
	req, err := s.makeRequest(InboundWebhook{
		Contacts: []InboundContact{{
			Profile: &Profile{
				Name: from,
			},
			WaID: from,
		}},
		Messages: []InboundMessage{{
			Message: Message{
				Type: "text",
				Text: &MessageText{
					Body: text,
				},
			},
			From: from,
			ID:   RandomString(27),
		}},
	})
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}
