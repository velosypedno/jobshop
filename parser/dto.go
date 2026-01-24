package parser

type FactoryConfig struct {
	Machines     []MachineConfig  `json:"machines"`
	JobTemplates []JobTemplateDTO `json:"job_templates"`
}

type MachineConfig struct {
	TypeID   int    `json:"type_id"`
	TypeName string `json:"type_name"`
	Count    int    `json:"count"`
}

type JobTemplateDTO struct {
	Name       string                 `json:"name"`
	Operations []OperationTemplateDTO `json:"operations"`
}

type OperationTemplateDTO struct {
	Name           string                 `json:"name"`
	MachineType    string                 `json:"machine_type"`
	ProcessingTime string                 `json:"processing_time"`
	Children       []OperationTemplateDTO `json:"children,omitempty"`
}
