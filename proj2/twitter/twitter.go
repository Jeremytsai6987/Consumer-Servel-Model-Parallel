package main

import (
	"encoding/json"
	"os"
	"proj2/server"
	"strconv"
)

func main() {
	var mode string
	var consumersCount int

	if len(os.Args) > 1 {
		count, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic("Invalid number of consumers")
		}
		mode = "p"
		consumersCount = count
	} else {
		mode = "s"
		consumersCount = 0
	}

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	config := server.Config{
		Encoder:       encoder,
		Decoder:       decoder,
		Mode:          mode,
		ConsumersCount: consumersCount,
	}

	server.Run(config)
}
