package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var form ContactForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	inbox := os.Getenv("TEAM_INBOX")
	if inbox == "" {
		inbox = "team@example.com"
	}
	result, err := routeContact(form, inbox)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/contact", contactHandler)
	log.Println("contact service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
