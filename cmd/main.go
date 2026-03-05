package main

import (
	"fmt"
	"os"
	"time"

	"github.com/velosypedno/resource-allocation/internal/app"
	"github.com/velosypedno/resource-allocation/internal/parser"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <factory_config_path> <orders_path> [optional_name]")
		os.Exit(1)
	}

	factoryConfigPath := os.Args[1]
	ordersPath := os.Args[2]
	customName := ""
	if len(os.Args) > 3 {
		customName = "_" + os.Args[3]
	}

	machinesConfig, templates, err := parser.ParseFactoryConfig(factoryConfigPath)
	if err != nil {
		fmt.Printf("Error parsing factory config: %v\n", err)
		os.Exit(1)
	}

	orders, err := parser.ParseOrders(ordersPath)
	if err != nil {
		fmt.Printf("Error parsing orders: %v\n", err)
		os.Exit(1)
	}

	a := app.New(machinesConfig, templates)

	startTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)

	err = a.Run(startTime, orders, customName)
}
