package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type CheckRequest struct {
	Links []string `json:"links"`
}

type CheckResponse struct {
	Links    map[string]string `json:"links"`
	LinksNum int 			   `json:"links_num"`
}

var (
	id 		int
	mu      sync.Mutex
)

func main() {
	http.HandleFunc("/check", handleCheck)
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}


func handleCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	results := make(map[string]string)
	var wg sync.WaitGroup
	for _, link := range req.Links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			results[link] = checkLink(link)
		}(link)
	}

	wg.Wait()

	mu.Lock()
	id ++
	mu.Unlock()

	resp := CheckResponse{
		Links:    results,
		LinksNum: id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func checkLink(link string) string {
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "https://" + link
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(link)
	if err != nil || resp.StatusCode >= 400 {
		return "not available"
	}
	defer resp.Body.Close()
	return "available"
}
