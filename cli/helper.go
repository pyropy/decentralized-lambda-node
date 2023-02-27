package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	ufcli "github.com/urfave/cli/v2"
)

func RunApp(app *ufcli.App) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		os.Exit(1)
	}()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err) // nolint:errcheck
	}
}
