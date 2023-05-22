package asemo

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Server struct {
	E *echo.Echo

	port             uint16
	sendEmailHandler SendEmailHandler
}

func NewServer() *Server {
	e := echo.New()
	e.HideBanner = true
	server := &Server{
		E:                e,
		port:             8080,
		sendEmailHandler: defaultSendEmailHandler,
	}
	e.POST("/v2/email/outbound-emails", server.sendEmail)
	return server
}

func (s *Server) Start() error {
	return s.E.Start(fmt.Sprintf("localhost:%d", s.port))
}

func (s *Server) SetPort(port uint16) {
	s.port = port
}

func (s *Server) SetSendEmailHandler(handler SendEmailHandler) {
	s.sendEmailHandler = handler
}
