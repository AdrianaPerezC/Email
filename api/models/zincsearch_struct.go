package models

// Estructura de la respuesta de ZincSearch
type ZincSearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source struct {
				Email
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// Estructura para solicitudes de búsqueda
type SearchRequest struct {
	Query struct {
		Term string `json:"term"`
	} `json:"query"`
	From      int    `json:"from"`
	Size      int    `json:"size"`
	SortField string `json:"sort_field"`
}
