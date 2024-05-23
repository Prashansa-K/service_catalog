package structs

import "time"

// single response structures
type ServiceResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	VersionCount int       `json:"version_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type ServiceVersion struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// paginated response structures
type ServicePaginationResponse struct {
	Services     []ServiceResponse `json:"services"`
	TotalPages   int               `json:"total_pages"`
	CurrentPage  int               `json:"current_page"`
	TotalRecords int64             `json:"total_records"`
}

type ServiceResponseWithVersionPagination struct {
	ID                  uint             `json:"id"`
	Name                string           `json:"name"`
	Description         string           `json:"description"`
	VersionCount        int              `json:"version_count"`
	CreatedAt           time.Time        `json:"created_at"`
	Versions            []ServiceVersion `json:"versions"`
	TotalPages          int              `json:"total_pages"`
	CurrentPage         int              `json:"current_page"`
	TotalVersionRecords int64            `json:"total_version_records"`
}
