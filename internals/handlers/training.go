package handlers

import (
	"net/http"

	"github.com/MatyD356/vimGame/internals/config"
	notionservice "github.com/MatyD356/vimGame/internals/integrations/notion/services"
)

func HandleGetTraining(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cfg := ctx.Value("config").(*config.Config) // Type assertion
	go notionservice.GetDatabase(cfg)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Training data is being fetched in the background. Please check back later."))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
