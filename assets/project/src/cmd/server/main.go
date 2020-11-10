package main

import (
	"flag"
	"fmt"
)

var configFile = flag.String("c", "server", "set config file")

func main() {
	flag.Parse()

	app, err := CreateApp(*configFile)
	if err != nil {
		fmt.Println("create fail:", err)
		panic(err)
	}

	if err := app.Start(); err != nil {
		fmt.Println("Start fail:", err)
		panic(err)
	}
	fmt.Println("AwaitSignal:", err)

	app.AwaitSignal()
}
