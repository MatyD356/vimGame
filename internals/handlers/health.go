package handlers

import (
	"fmt"
	"net/http"
)

func HandleHealt(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server is healthy")
}
