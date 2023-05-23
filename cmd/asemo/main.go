package main

import (
	"flag"
	"log"

	"github.com/lusingander/asemo/asemo"
)

func main() {
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

	app.server.SetPort(args.port)
	app.server.SetSendEmailHandler(app.sendEmailHandler)

	return app.server.Start()
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
