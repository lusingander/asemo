package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/lusingander/asemo/asemo"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	args, err := parseArgs()
	if err != nil {
		return err
	}

	app := &application{
		server:            asemo.NewServer(),
		messageRepository: newMessageRepository(),
	}

	app.server.E.HidePort = true

	app.server.SetPort(args.port)
	app.server.SetSendEmailHandler(app.sendEmailHandler)

	app.server.E.GET("/api/messages", app.listMessageHandler)
	app.server.E.GET("/api/messages/:id", app.getMessageHandler)

	app.server.E.Use(loggingRequest)

	startLog(app, args)

	return app.server.Start()
}

func loggingRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		log.Printf("Request: %s %s\n", req.Method, req.RequestURI)
		err := next(c)
		log.Printf("Response: %s %s %d %d\n", req.Method, req.RequestURI, res.Status, res.Size)
		return err
	}
}

func startLog(app *application, args *args) {
	baseUrl := fmt.Sprintf("http://localhost:%d", args.port)
	routes := app.server.E.Routes()
	max := methodMaxLen(routes)

	log.Printf("Start mock server on %s\n", baseUrl)
	for _, r := range routes {
		// fixme: should be sorted because order of routes is random...
		log.Printf("* %-*s %s%s\n", max, r.Method, baseUrl, r.Path)
	}
}

func methodMaxLen(rs []*echo.Route) int {
	max := 0
	for _, r := range rs {
		l := len(r.Method)
		if max < l {
			max = l
		}
	}
	return max
}

type application struct {
	server *asemo.Server
	*messageRepository
}

type args struct {
	port uint16
}

func parseArgs() (*args, error) {
	port := flag.Uint("port", 8080, "port number")
	flag.Parse()
	return &args{
		port: uint16(*port),
	}, nil
}
