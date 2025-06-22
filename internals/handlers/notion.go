package handlers

import "net/http"

func HandleNotion(w http.ResponseWriter, r *http.Request) {
	//https://api.notion.com/v1/databases/:id make a request to Notion API
	// This handler is a placeholder for the Notion API integration.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notion handler is not implemented yet"))
}
