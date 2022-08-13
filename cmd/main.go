package main

import (
	"log"
	"net/http"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"

	"github.com/Reverse-Labs/webmon"
)

func main() {
	app := cli.App("webmon", "web monitoring tool")
	app.Command("run", "start monitoring", cmdRun)
	app.Run(os.Args)
}

func cmdRun(cmd *cli.Cmd) {
	var (
		configFile = cmd.StringArg("CONFIG", "", "config file")
	)

	cmd.Action = func() {
		fd, err := os.Open(*configFile)

		if err != nil {
			log.Fatal("unable to open config file: ", *configFile)
		}
		config := make([]webmon.ProbeConfig, 0)
		if err := yaml.NewDecoder(fd).Decode(&config); err != nil {
			log.Fatal(err)
		}

		webmon.Monitor(config)

		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":21122", nil); err != nil {
			log.Fatalln(err)
		}
	}
}
