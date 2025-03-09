package credentialsdb

import (
	"fmt"
	"net/http"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create New Entity endpoint hit")
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete Entity endpoint hit")
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Update Entity endpoint hit")
}

func FetchWithIdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get Entity endpoint hit")
}
