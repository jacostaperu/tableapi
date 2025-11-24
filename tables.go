package tableapi

type GenericTable struct {
	Records []GenericRecord `json:"records"`
}

type GenericRecord struct {
	ID     string            `json:"id"`
	Fields map[string]string `json:"fields"`
}
