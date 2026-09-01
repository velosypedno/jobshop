package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

func ParseOrders(filePath string) ([]OrderDTO, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read orders file: %w", err)
	}

	var config OrderConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal orders json: %w", err)
	}

	return config.Orders, nil
}
