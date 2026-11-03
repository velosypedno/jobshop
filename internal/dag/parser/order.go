package parser

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/velosypedno/jobshop/pkg/dag/core"
)

type orderConfig struct {
	Orders []orderDTO `json:"orders"`
}

type orderDTO struct {
	JobTemplateID int `json:"job_template_id"`
	Amount        int `json:"amount"`
}

func ParseOrders(file io.Reader) (core.Order, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read orders file: %w", err)
	}

	var config orderConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal orders json: %w", err)
	}

	orderEntries := make([]core.OrderEntry, 0, len(config.Orders))

	for _, orderDTO := range config.Orders {
		orderEntries = append(orderEntries,
			core.OrderEntry{
				JobTemplateID: core.JobTemplateID(orderDTO.JobTemplateID),
				Amount:        orderDTO.Amount,
			})
	}

	return orderEntries, nil
}
