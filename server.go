package main

import (
	"crypto/tls"
	"errors"
	"github.com/DeanThompson/ginpprof"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"go-api/common"
	"go-api/router"
	"net/http"
	"strings"
	"time"
)

func server(c *cli.Context, started chan bool) error {
	deployed := make(chan bool)
	// Initialize the logger.
	common.LogInit(c)
	common.LogInfo("server(), runs the server.")

	if strings.ToUpper(c.String("log-level")) != "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create the Gin engine.
	g := gin.New()

	// Load config file and save in global var CFG
	//common.DBParse()
	common.ConfigParse()
	// Load pprof
	ginpprof.Wrapper(g)

	// Websocket handlers.
	//w := common.WsNew(g)

	router.Load(
		g,
		common.MiddlewareConfig(c),
		common.DatabaseConnect("miku"),
		common.RateLimit(),
		common.MiddleLogging(),
	)
	go func() {
		if err := pingServer(c); err != nil {
			common.LogFatal("The router has no response, or it might took too long to start up.")
		}
		common.LogInfo("The router has been deployed successfully.")
		close(deployed)
		close(started)
	}()

	return endless.ListenAndServe(c.String("web-listen"), g)
}

func pingServer(c *cli.Context) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	for i := 0; i < 5; i++ {
		resp, err := client.Get(c.String("addr") + "/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		common.LogInfo("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
