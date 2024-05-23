package structs

type ServiceRequest struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ServiceVersionRequest struct {
	Name        string `json:"name"`
	ServiceName string `json:"service_name"`
	Description string `json:"description"`
}
