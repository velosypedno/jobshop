package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/velosypedno/jobshop/internal/tree/engine"
	"github.com/velosypedno/jobshop/internal/tree/report"
	"github.com/velosypedno/jobshop/pkg/tree/core"
	"github.com/velosypedno/jobshop/pkg/tree/digitaltwin"
	"go.uber.org/zap"
)

type App struct {
	ManufactoryTwin *digitaltwin.Manufactory
	Engine          *engine.Engine
}

func New(machinesConfig []core.MachineConfigEntry, templates []core.JobTemplate, strategies []core.Strategy) *App {
	return &App{
		ManufactoryTwin: digitaltwin.New(machinesConfig, templates),
		Engine:          engine.New(strategies...),
	}
}

func (a *App) Run(startTime time.Time, orders []core.OrderEntry, customName string) error {
	// Step 1: prepare output directories and paths
	outputDir := "results"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating directory: %v", err)
	}
	timestamp := time.Now().Format("20060102-150405")
	baseName := fmt.Sprintf("plan_%s%s", timestamp, customName)

	// Step 2: setup scheduler engine that will run problem against each sceduling strategy
	logPath := filepath.Join(outputDir, baseName+".log")
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logPath}

	logger, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	a.Engine.SetLogger(logger)

	logger.Info("Starting application run",
		zap.Time("start_time", startTime),
		zap.Int("orders_count", len(orders)),
	)

	// Step 3
	problem := a.ManufactoryTwin.NewProblem(orders)

	// Step 4
	schedulingReports, err := a.Engine.Solve(problem)
	if err != nil {
		logger.Error("Planning failed", zap.Error(err))
		return fmt.Errorf("during planning: %v", err)
	}

	// Step 5: create reports
	stdoutTable := report.NewSimpleTable(os.Stdout)
	if err := stdoutTable.Report(schedulingReports); err != nil {
		logger.Warn("Could not generate text report", zap.Error(err))
	}

	chartPath := filepath.Join(outputDir, baseName+".html")
	outputFile, err := os.Create(chartPath)
	if err != nil {
		return fmt.Errorf("creating output file: %v", err)
	}
	defer outputFile.Close()

	ganttCharts := report.NewGanttCharts(outputFile)
	if err := ganttCharts.Report(problem, startTime, schedulingReports); err != nil {
		logger.Warn("Could not generate gantt charts", zap.Error(err))
	}

	fmt.Printf("Successfully generated chart: %s\n", chartPath)
	fmt.Printf("Log file created: %s\n", logPath)

	logger.Info("Run completed successfully", zap.String("chart_path", chartPath))

	return nil
}
