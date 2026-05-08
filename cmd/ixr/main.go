package main

import (
	"log/slog"
	"os"

	ixr "github.com/ixr/ixr/pkg/ixr"
)

func main() {
	if err := ixr.Start(); err != nil {
		slog.Error("ixr exited", "err", err)
		os.Exit(1)
	}
}
