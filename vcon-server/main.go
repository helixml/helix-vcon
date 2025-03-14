package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type Party struct {
    Name string `json:"name"`
    Tel  string `json:"tel"`
}

type DialogEntry struct {
    Timestamp time.Time `json:"timestamp"`
    Speaker   string    `json:"speaker"`
    Text      string    `json:"text"`
}

type VCON struct {
    UUID      string        `json:"uuid"`
    CreatedAt time.Time    `json:"created_at"`
    Subject   string       `json:"subject"`
    Parties   []Party      `json:"parties"`
    Dialog    []DialogEntry `json:"dialog"`
}

var calls = []VCON{
    {
        UUID:      "call-001",
        CreatedAt: time.Date(2024, 2, 7, 14, 30, 0, 0, time.UTC),
        Subject:   "Project Timeline Discussion",
        Parties: []Party{
            {Name: "Frank Smith", Tel: "123-456-7890"},
            {Name: "Obi Johnson", Tel: "098-765-4321"},
        },
        Dialog: []DialogEntry{
            {
                Timestamp: time.Date(2024, 2, 7, 14, 30, 0, 0, time.UTC),
                Speaker:   "Frank Smith",
                Text:     "Hi Obi, let's discuss the project timeline.",
            },
            {
                Timestamp: time.Date(2024, 2, 7, 14, 31, 0, 0, time.UTC),
                Speaker:   "Obi Johnson",
                Text:     "Sure Frank, I've reviewed the milestones.",
            },
        },
    },
}

func main() {
    http.HandleFunc("/vcon", getAllVcons)
    http.HandleFunc("/vcons/search", searchVcons)
    
    log.Printf("Starting VCON server on :8005")
    log.Fatal(http.ListenAndServe("0.0.0.0:8005", nil))
}

func getAllVcons(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(calls)
}

func searchVcons(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    name := r.URL.Query().Get("name")
    tel := r.URL.Query().Get("tel")
    
    var results []VCON
    for _, call := range calls {
        for _, party := range call.Parties {
            if (name != "" && party.Name == name) || 
               (tel != "" && party.Tel == tel) {
                results = append(results, call)
                break
            }
        }
    }
    
    json.NewEncoder(w).Encode(results)
}
