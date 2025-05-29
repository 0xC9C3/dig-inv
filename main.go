package main

import (
	"dig-inv/gateway"
	"dig-inv/log"
	"github.com/alecthomas/kong"
	"google.golang.org/grpc/grpclog"
)

// https://github.com/alecthomas/kong
// https://entgo.io/

var CLI struct {
	Server struct {
	} `cmd:"" help:"Run the server."`

	Worker struct {
	} `cmd:"" help:"Run the worker."`
}

func main() {
	log.S().Info("Starting dig-inv")

	ctx := kong.Parse(&CLI)

	log.S().Debugw("Parsed CLI context", "command", ctx.Command(), "args", ctx.Args)

	switch ctx.Command() {
	case "server":
		server()
	case "worker":
		worker()
	}

}

func server() {
	log.S().Info("Running server")

	if err := gateway.Run(); err != nil {
		grpclog.Fatal(err)
	}
}

func worker() {
	log.S().Info("Running worker")
}
