package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const Port int = 3030

type DataItem struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"createdAt"`
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./client")))
	mux.HandleFunc("/api/data", apiDataHandler)

	log.Printf("Server listening on : http://[::1]:%d", Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), mux))
}

func apiDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := os.ReadFile("data/db.json")
	if err != nil {
		http.Error(w, "Server error reading data", http.StatusInternalServerError)
		return
	}

	var items []DataItem
	if err := json.Unmarshal(data, &items); err != nil {
		http.Error(w, "Server error parsing data", http.StatusInternalServerError)
		return
	}

	params := r.URL.Query()
	switch r.Method {
	case http.MethodGet:
		s := params.Get("s")

		if s != "" {
			filteredItems := make([]DataItem, 0)

			for _, item := range items {
				if containsSearchTerm(item, strings.ToLower(s)) {
					filteredItems = append(filteredItems, item)
				}
			}

			json.NewEncoder(w).Encode(filteredItems)
		} else {
			json.NewEncoder(w).Encode(items)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func containsSearchTerm(item DataItem, term string) bool {
	return strings.Contains(strings.ToLower(item.Title), term) ||
		strings.Contains(strings.ToLower(item.Content), term) ||
		strings.Contains(strings.ToLower(item.Details), term)
}
