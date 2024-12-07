package email

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AdrianaPerezC/Email/models"
	"github.com/AdrianaPerezC/Email/zincsearch"
)

func GetEmails(w http.ResponseWriter, r *http.Request) {
	var request models.SearchRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error":"La estructura de la petición debe ser la indicada es obligatorio"}`, http.StatusBadRequest)
		return
	}

	// Validar el campo `term`
	if request.Query.Term == "" {
		http.Error(w, `{"error":"El campo 'term' no puede estar vacío"}`, http.StatusBadRequest)
		return
	}

	// Llamar a ZincSearch
	searchResp, err := zincsearch.SearchZinc(request)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	// Devolver los resultados
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(searchResp)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
