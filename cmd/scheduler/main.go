package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/velosypedno/jobshop/internal/app"
	"github.com/velosypedno/jobshop/internal/parser"
)

func main() {
	var configPath string
	var ordersPath string
	var customName string

	var rootCmd = &cobra.Command{
		Use:   "scheduler",
		Short: "Resource allocation scheduler for factory production",
		Long:  `A tool to optimize production plans using various strategies like Simulated Annealing and Genetic Algorithms.`,
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}

	rootCmd.Flags().StringVarP(&configPath, "config", "c", "example/config.json", "path to factory configuration file")
	rootCmd.Flags().StringVarP(&ordersPath, "orders", "o", "example/order.json", "path to orders file")
	rootCmd.Flags().StringVarP(&customName, "name", "n", "", "custom name for the output report")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
