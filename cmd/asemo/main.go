package main

import "github.com/lusingander/asemo/asemo"

func run() error {
	return asemo.NewServer().Start()
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
