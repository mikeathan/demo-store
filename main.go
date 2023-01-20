package main

import (
	"demo-store/server"
	"demo-store/utils"
	"flag"
	"os"
)

type Args struct {
	Port  int
	Depth int
}

func main() {
	args := readArgs()
	defer utils.CloseLoggers()

	listen(args)
}

func listen(args Args) {
	err := server.Listen(args.Port, args.Depth)
	if err != nil {
		utils.ApplicationTracer().LogError(err)
		os.Exit(-2)
	}
}

func readArgs() Args {
	var port int
	var depth int

	flag.IntVar(&port, "port", -1, "port to listen on")
	flag.IntVar(&depth, "depth", 0, "LRU depth")
	flag.Parse()

	if port == -1 {
		utils.ApplicationTracer().LogError("Error: port not specified")
		os.Exit(-1)
	}

	return Args{Port: port, Depth: depth}
}
