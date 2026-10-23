package parser

type factoryConfig struct {
	Machines     []machineConfigDTO `json:"machines"`
	JobTemplates []jobTemplateDTO   `json:"job_templates"`
}

type machineConfigDTO struct {
	TypeID   int    `json:"type_id"`
	TypeName string `json:"type_name"`
	Count    int    `json:"count"`
}

type (
	jobTemplateDTO struct {
		ID int `json:"id"`

		Name       string                 `json:"name"`
		Operations []operationTemplateDTO `json:"operations"`
		Edges      []edge                 `json:"edges"`
	}

	edge struct {
		FromOpID int `json:"from_op_id"`
		ToOpID   int `json:"to_op_id"`
	}

	operationTemplateDTO struct {
		ID int `json:"id"`

		Name           string `json:"name"`
		MachineType    int    `json:"machine_type"`
		ProcessingTime string `json:"processing_time"`
	}
)
