package asemo

import (
	"github.com/labstack/echo/v4"
)

type Server struct {
	E                *echo.Echo
	SendEmailHandler func(SendEmailRequest) SendEmailResponse
}

func NewServer() *Server {
	e := echo.New()
	server := &Server{
		E:                e,
		SendEmailHandler: defaultSendEmailHandler,
	}
	e.POST("/v2/email/outbound-emails", server.sendEmail)
	return server
}

func (s *Server) Run() error {
	return s.E.Start("localhost:8080")
}
