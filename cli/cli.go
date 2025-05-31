package cli

import (
	"dig-inv/log"
	"fmt"
	"github.com/alecthomas/kong"
)

type Handler struct {
	ServerHandler func() error
	WorkerHandler func() error
}

type ServerParams struct {
}

func (r *ServerParams) Run(ctx *Handler) error {
	if ctx.ServerHandler == nil {
		log.S.Warn("Server handler is not set, skipping server execution")
		return nil
	}
	return ctx.ServerHandler()
}

type WorkerParams struct {
}

func (l *WorkerParams) Run(ctx *Handler) error {
	if ctx.WorkerHandler == nil {
		log.S.Warn("Worker handler is not set, skipping worker execution")
		return nil
	}
	return ctx.WorkerHandler()
}

var Wrapper struct {
	Server ServerParams `cmd:"" help:"Run the server"`
	Worker WorkerParams `cmd:"" help:"Run the worker"`
}

const (
	CommandServer = "server"
	CommandWorker = "worker"
)

func (cli *Handler) Run() error {
	return cli.RunWithOptions()
}

func (cli *Handler) RunWithOptions(opts ...kong.Option) error {
	ctx := kong.Parse(&Wrapper, opts...)

	log.S.Debugw("Parsed CLI context", "command", ctx.Command(), "args", ctx.Args)

	if err := ctx.Run(cli); err != nil {
		log.S.Errorw("Failed to run CLI command", "error", err)
		return fmt.Errorf("failed to run CLI command: %w", err)
	}

	log.S.Debugw("CLI command executed successfully", "command", ctx.Command())
	return nil
}

func NewCLI(serverHandler, workerHandler func() error) *Handler {
	return &Handler{
		ServerHandler: serverHandler,
		WorkerHandler: workerHandler,
	}
}
