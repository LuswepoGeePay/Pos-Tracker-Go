package models

type PocketBaseAuthResponse struct {
	Token  string                 `json:"token"`
	Record map[string]interface{} `json:"record"`
}

type GetRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func (p *GetRequest) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}

	if p.PageSize == 0 {
		p.PageSize = 10
	}
}

type SearchRequest struct {
	GetRequest
	SearchQuery string `json:"searchQuery,omitempty"`
}
