package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/velosypedno/resource-allocation/internal/base"
	"github.com/velosypedno/resource-allocation/internal/parser"
)

func GenerateAllJobCharts(machineConfigs []parser.MachineConfig, templates []base.JobTemplate, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	machineMap := make(map[base.MachineType]string)
	for _, m := range machineConfigs {
		machineMap[base.MachineType(m.TypeID)] = m.TypeName
	}

	for _, job := range templates {
		safeName := strings.ReplaceAll(job.Name, " ", "_")
		dotPath := filepath.Join(outputDir, safeName+".dot")
		pngPath := filepath.Join(outputDir, safeName+".png")

		file, err := os.Create(dotPath)
		if err != nil {
			return err
		}

		if err := writeJobDot(job, machineMap, file); err != nil {
			file.Close()
			return err
		}
		file.Close()

		cmd := exec.Command("dot", "-Tpng", dotPath, "-o", pngPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to render PNG for %s: %v", job.Name, err)
		}

		os.Remove(dotPath)
	}

	return nil
}

func writeJobDot(job base.JobTemplate, machineMap map[base.MachineType]string, w io.Writer) error {
	fmt.Fprintln(w, "digraph G {")
	fmt.Fprintln(w, "  rankdir=LR;")
	fmt.Fprintln(w, "  bgcolor=\"transparent\";")
	fmt.Fprintln(w, "  nodesep=0.4;")
	fmt.Fprintln(w, "  ranksep=0.7;")
	fmt.Fprintln(w, "  splines=polyline;")

	fmt.Fprintln(w, "  node [shape=none, fontname=\"Arial\"];")
	fmt.Fprintln(w, "  edge [color=\"#718096\", arrowhead=vee, arrowsize=0.8, penwidth=1.1];")

	colors := map[base.MachineType][2]string{
		0: {"#ebf8ff", "#3182ce"},
		1: {"#faf5ff", "#805ad5"},
		2: {"#f0fff4", "#38a169"},
		3: {"#fffaf0", "#dd6b20"},
		4: {"#fff5f5", "#e53e3e"},
		5: {"#f7fafc", "#4a5568"},
	}

	nodeIdx := 0
	var traverse func(op base.OperationTemplate) string
	traverse = func(op base.OperationTemplate) string {
		nodeIdx++
		currID := fmt.Sprintf("n%d", nodeIdx)

		mName := machineMap[op.MachineType]
		if mName == "" {
			mName = "Unknown"
		}

		style, ok := colors[op.MachineType]
		if !ok {
			style = [2]string{"#ffffff", "#cbd5e0"}
		}

		label := fmt.Sprintf(`<<table border="1" cellborder="0" cellspacing="0" cellpadding="8" bgcolor="%s" color="%s" style="rounded">
			<tr><td><b>%s</b></td></tr>
			<hr/>
			<tr><td><font color="#4a5568" point-size="9">%s | %v</font></td></tr>
		</table>>`, style[0], style[1], op.Name, mName, op.ProcessingTime)

		fmt.Fprintf(w, "  %s [label=%s];\n", currID, label)

		for _, child := range op.Children {
			childID := traverse(child)
			fmt.Fprintf(w, "  %s -> %s;\n", childID, currID)
		}
		return currID
	}

	for _, root := range job.Operations {
		traverse(root)
	}

	fmt.Fprintln(w, "}")
	return nil
}

func main() {
	path := "./default/config.json"
	machineConfig, templates, _, err := parser.ParseFactoryConfig(path)
	if err != nil {
		panic(err)
	}

	err = GenerateAllJobCharts(machineConfig, templates, "./pngs/")
	if err != nil {
		fmt.Printf("Помилка генерації: %v\n", err)
	} else {
		fmt.Println("Success ./pngs/")
	}
}
