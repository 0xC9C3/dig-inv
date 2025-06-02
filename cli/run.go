package cli

import (
	"dig-inv/log"
	"dig-inv/services"
)

// https://entgo.io/

const ErrorExitCode = 0xF1

type Entrypoint struct {
	serverHandler func() error
	workerHandler func() error
}

func (e *Entrypoint) Run() int {
	log.S.Info("Starting dig-inv")
	if err := NewCLI(e.serverHandler, e.workerHandler).Run(); err != nil {
		log.S.Errorw("Failed to run CLI", "error", err)

		return ErrorExitCode
	}

	log.S.Info("Exiting gracefully")
	return 0
}

func NewEntrypoint(
	serverHandler, workerHandler func() error,
) *Entrypoint {
	return &Entrypoint{
		serverHandler: serverHandler,
		workerHandler: workerHandler,
	}
}

func Run() int {
	return NewEntrypoint(
		services.NewGatewayServer().Run,
		worker,
	).Run()
}

func worker() error {
	log.S.Info("Running as worker")

	return nil
}
