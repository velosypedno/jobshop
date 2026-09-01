package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/velosypedno/jobshop/internal/tree/app"
	"github.com/velosypedno/jobshop/internal/tree/parser"
)

func main() {
	var (
		configPath string
		ordersPath string
		customName string
	)

	flag.StringVar(&configPath, "config", "example/config.json", "path to factory configuration file")
	flag.StringVar(&configPath, "c", "example/config.json", "path to factory configuration file (shorthand)")

	flag.StringVar(&ordersPath, "orders", "example/order.json", "path to orders file")
	flag.StringVar(&ordersPath, "o", "example/order.json", "path to orders file (shorthand)")

	flag.StringVar(&customName, "name", "", "custom name for the output report")
	flag.StringVar(&customName, "n", "", "custom name for the output report (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Resource allocation scheduler for factory production.\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	machinesConfig, templates, strategies, err := parser.ParseFactoryConfig(configPath)
	if err != nil {
		fmt.Printf("Error parsing factory config: %v\n", err)
		os.Exit(1)
	}

	orders, err := parser.ParseOrders(ordersPath)
	if err != nil {
		fmt.Printf("Error parsing orders: %v\n", err)
		os.Exit(1)
	}

	a := app.New(machinesConfig, templates, strategies)
	startTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)

	err = a.Run(startTime, orders, customName)
	if err != nil {
		fmt.Printf("Application run error: %v\n", err)
		os.Exit(1)
	}
}
