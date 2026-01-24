package chart

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/velosypedno/resource-allocation/base"
	"github.com/velosypedno/resource-allocation/factory"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const renderItemJS = `function (params, api) {
    var categoryIndex = api.value(0); 
    var start = api.coord([api.value(1), categoryIndex]); 
    var end = api.coord([api.value(2), categoryIndex]); 
    var height = api.size([0, 1])[1] * 0.6; 
    var rectShape = echarts.graphic.clipRectByRect({ x: start[0], y: start[1] - height / 2, width: Math.max(end[0] - start[0], 1), height: height }, { x: params.coordSys.x, y: params.coordSys.y, width: params.coordSys.width, height: params.coordSys.height }); 
    return rectShape && { type: 'rect', transition: ['shape'], shape: rectShape, style: api.style() }; 
}`

const renderTooltipJS = `function(p){
	var dateStart = new Date(p.value[1]);
	var dateEnd = new Date(p.value[2]);
	var timeStart = dateStart.toLocaleTimeString('uk-UA', {hour12: false});
	var timeEnd = dateEnd.toLocaleTimeString('uk-UA', {hour12: false});
	return '<b>' + p.value[4] + '</b><br/>' + 
			'Operation: ' + p.value[3] + '<br/>' + 
			timeStart + ' - ' + timeEnd;
}`

func sortMachines(machines []*base.Machine) {
	sort.Slice(machines, func(i, j int) bool {
		if machines[i].Type != machines[j].Type {
			return machines[i].Type < machines[j].Type
		}
		return machines[i].ID < machines[j].ID
	})
}

func generateMachineIndexMap(machines []*base.Machine) map[base.MachineID]int {
	mMap := make(map[base.MachineID]int)
	for i, m := range machines {
		mMap[m.ID] = i
	}
	return mMap
}

func generateYAxisCategories(machines []*base.Machine) []string {
	categories := make([]string, len(machines))
	for _, m := range machines {
		categories = append(categories, fmt.Sprintf("%s [ID: %d]", m.Type.String(), m.ID))
	}
	return categories
}

func createBaseCustomChart(machines []*base.Machine, period base.Period, description string) *charts.Custom {
	chart := charts.NewCustom()

	lineCount := strings.Count(description, "\n") + 1
	descriptionHeight := lineCount * 27
	topOffset := descriptionHeight + 50

	baseHeight := len(machines)*75 + 120
	totalHeight := baseHeight + descriptionHeight

	chart.Initialization.Width = "90%"
	chart.Initialization.Height = fmt.Sprintf("%dpx", totalHeight)

	yAxisCategories := generateYAxisCategories(machines)

	chart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Production Plan",
			Subtitle: description,
			Left:     "12%",
			TitleStyle: &opts.TextStyle{
				FontSize:   24,
				FontWeight: "bold",
				Color:      "#1a1a1a",
			},
			SubtitleStyle: &opts.TextStyle{
				Color:      "#333",
				FontSize:   16,
				FontWeight: "500",
				LineHeight: 25,
				FontFamily: "monospace, Courier New",
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      opts.Bool(true),
			Trigger:   "item",
			Formatter: opts.FuncOpts(renderTooltipJS),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type:      "time",
			Min:       period.Start.UnixMilli(),
			Max:       period.End.UnixMilli(),
			SplitLine: &opts.SplitLine{Show: opts.Bool(true)},
			AxisLabel: &opts.AxisLabel{
				Show:       opts.Bool(true),
				Formatter:  "{HH}:{mm}:{ss}",
				FontSize:   14,
				FontWeight: "bold",
				Color:      "#333",
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type:      "category",
			Data:      yAxisCategories,
			SplitLine: &opts.SplitLine{Show: opts.Bool(true)},
			AxisLabel: &opts.AxisLabel{
				Show:       opts.Bool(true),
				FontSize:   14,
				FontWeight: "600",
				Color:      "#222",
				Margin:     15,
			},
		}),
		charts.WithGridOpts(opts.Grid{
			Top:          fmt.Sprintf("%dpx", topOffset),
			Left:         "5%",
			Right:        "10%",
			ContainLabel: opts.Bool(true),
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:   opts.Bool(true),
			Orient: "vertical",
			Right:  "5px",
			Top:    fmt.Sprintf("%dpx", topOffset),
			Type:   "scroll",
		}),
	)

	return chart
}

func addSolutionSeries(chart *charts.Custom, solution base.Solution, machines []*base.Machine) {
	mMap := generateMachineIndexMap(machines)

	for _, job := range solution.Jobs {
		var seriesData []opts.CustomData
		fullJobName := fmt.Sprintf("%s [%d]", job.Job.Name, job.Job.ID)

		for _, op := range job.GetAllOperations() {
			seriesData = append(seriesData, opts.CustomData{
				Value: []interface{}{
					mMap[op.MachineID],
					op.Period.Start.UnixMilli(),
					op.Period.End.UnixMilli(),
					op.Operation.Name,
					fullJobName,
				},
			})
		}

		chart.AddSeries(fullJobName, seriesData).
			SetSeriesOptions(
				charts.WithCustomChartOpts(opts.CustomChart{
					RenderItem: opts.FuncOpts(renderItemJS),
				}),
			)
	}
}

func formatStrategyDescription(meta factory.SchedulingInfo) string {
	execTime := meta.SchedulingTime.Round(time.Millisecond).String()
	makespan := meta.MakeSpan.String()
	utilization := fmt.Sprintf("%.1f%%", meta.UtilizationLevel*100)

	line1 := fmt.Sprintf("STRATEGY: %s", strings.ToUpper(meta.StrategyName))

	line2 := fmt.Sprintf("TIME: %s  │  MAKESPAN: %s  │  UTILIZATION: %s",
		execTime, makespan, utilization)

	return fmt.Sprintf(
		"%s\n%s\n"+
			"─────────────────────────────────────────────────────────────────\n"+
			"%s",
		line1,
		line2,
		meta.StrategyDescription,
	)
}

func GenerateFromSolution(
	solution base.Solution,
	machines []*base.Machine,
	schedulingInfo factory.SchedulingInfo,
) *charts.Custom {
	sortMachines(machines)
	period := solution.GetWorkFlowPeriod()
	description := formatStrategyDescription(schedulingInfo)

	chart := createBaseCustomChart(machines, period, description)
	addSolutionSeries(chart, solution, machines)
	return chart
}
