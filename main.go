package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var ashMikuGoApi = "1.0.0"

func main() {
	started := make(chan bool)

	app := &cli.App{
		Name:    "Ash_Miku_Go_Api",
		Usage:   "A Web Api for Golang",
		Version: ashMikuGoApi,
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			EnvVars: []string{"MICROSERVICE_PORT"},
			Name:    "web-listen",
			Usage:   "the web port of the service.",
			Value:   ":9999",
		},
		&cli.StringFlag{
			EnvVars: []string{"MICROSERVICE_LOG_LEVEL"},
			Name:    "log-level",
			Usage:   "the log level mode.",
			Value:   "DEBUG",
		},
		&cli.StringFlag{
			EnvVars: []string{"MICROSERVICE_ADDR"},
			Name:    "addr",
			Usage:   "the web address of the service.",
			Value:   "http://127.0.0.1:9999",
		},
		&cli.StringFlag{
			EnvVars: []string{"MICROSERVICE_JWT_SECRET"},
			Name:    "jwt-secret",
			Usage:   "the secert used to encode the json web token.",
			Value:   "4Rtg8BPKwixXy2ktDPxoMMAhRzmo9mmuZjvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		},
	}

	app.Action = func(context *cli.Context) error {
		err := server(context, started)
		if err != nil {
			return err
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("cli app run error: %v", err)
		return
	}
}
