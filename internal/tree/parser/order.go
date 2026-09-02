package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/velosypedno/jobshop/pkg/tree/core"
)

type orderConfig struct {
	Orders []orderDTO `json:"orders"`
}

type orderDTO struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

func ParseOrders(filePath string) ([]core.OrderEntry, error) {
	data, err := os.ReadFile(filePath)
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
				Name:   orderDTO.Name,
				Amount: orderDTO.Amount,
			})
	}

	return orderEntries, nil
}
