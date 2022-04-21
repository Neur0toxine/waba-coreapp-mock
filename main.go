package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mkideal/cli"
)

type CLI struct {
	cli.Helper2
	Address string `cli:"*addr,address" usage:"Address to listen"`
	Verbose bool   `cli:"v,verbose" usage:"Enable verbose logging"`
}

func main() {
	os.Exit(cli.Run(new(CLI), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CLI)

		gin.SetMode(gin.ReleaseMode)
		if argv.Verbose {
			gin.SetMode(gin.DebugMode)
		}

		http.DefaultClient.Timeout = time.Second * 30

		return NewServer().Run(argv.Address)
	}))
}
