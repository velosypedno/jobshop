package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/velosypedno/jobshop/internal/dag/engine"
	"github.com/velosypedno/jobshop/internal/dag/parser"
	"github.com/velosypedno/jobshop/internal/dag/report"
	"github.com/velosypedno/jobshop/pkg/dag/core"
	"github.com/velosypedno/jobshop/pkg/dag/strategy"
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

	machinesConfig, templates, err := parser.ParseFactoryConfig(configPath)
	if err != nil {
		fmt.Printf("Error parsing factory config: %v\n", err)
		os.Exit(1)
	}

	orders, err := parser.ParseOrders(ordersPath)
	if err != nil {
		fmt.Printf("Error parsing orders: %v\n", err)
		os.Exit(1)
	}

	manufactory, err := core.NewManufactory(machinesConfig, templates)
	if err != nil {
		fmt.Printf("Error creating maufactory: %v\n", err)
		os.Exit(1)
	}

	problem, err := manufactory.NewProblem(orders)
	if err != nil {
		fmt.Printf("Error creating problem: %v\n", err)
		os.Exit(1)
	}

	randomStrategy := strategy.RandomStrategy{}
	e := engine.New(&randomStrategy)

	r, err := e.Solve(problem)
	if err != nil {
		fmt.Printf("Error creating report: %v\n", err)
		os.Exit(1)
	}

	stdoutTable := report.NewSimpleTable(os.Stdout)
	if err := stdoutTable.Report([]engine.Report{r}); err != nil {
		fmt.Printf("Could not generate text report: %v", err)
		os.Exit(1)
	}

	outputDir := "results-dag"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("creating directory: %v", err)
		os.Exit(1)
	}
	timestamp := time.Now().Format("20060102-150405")
	baseName := fmt.Sprintf("plan_%s%s", timestamp, customName)

	chartPath := filepath.Join(outputDir, baseName+".html")
	outputFile, err := os.Create(chartPath)
	if err != nil {
		fmt.Printf("creating output file: %v", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	ganttCharts := report.NewGanttCharts(outputFile)
	startTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)
	if err := ganttCharts.Report(problem, startTime, []engine.Report{r}); err != nil {
		fmt.Printf("Could not generate gantt charts: %v", err)
		os.Exit(1)
	}

}
