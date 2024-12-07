package zincsearch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AdrianaPerezC/Email/models"
)

func SearchZinc(query models.SearchRequest) (models.ZincSearchResponse, error) {
	// Crear el cuerpo de la solicitud
	body := `{
	    "search_type": "match",
	    "query":
	    {
	        "term": "` + query.Query.Term + `"
	    },
		"sort_fields": [` + fmt.Sprint(query.SortField) + `],
	    "from": ` + fmt.Sprint(query.From) + `,
	    "max_results": ` + fmt.Sprint(query.Size) + `,
	    "_source": []
	}`
	// Construir la solicitud HTTP
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/api/%s/_search", os.Getenv("ZINC_HOST"), os.Getenv("ZINC_PORT"), os.Getenv("ZINC_INDEX_NAME")), strings.NewReader(body))
	if err != nil {
		return models.ZincSearchResponse{}, fmt.Errorf("error al crear la solicitud: %v", err)
	}

	// Autenticación básica
	req.SetBasicAuth(os.Getenv("ZINC_USER"), os.Getenv("ZINC_PASSWORD"))
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	client := &http.Client{
		Timeout: time.Minute * 60, // Aumentar el timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return models.ZincSearchResponse{}, fmt.Errorf("error al enviar la solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Leer y procesar la respuesta
	var searchResp models.ZincSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return models.ZincSearchResponse{}, fmt.Errorf("error al deserializar la respuesta: %v", err)
	}

	return searchResp, nil
}
