# Helix VCON Telco Bot

A Helix-powered application that provides natural language interface for call analysis using VCON standard.
This project helps you deploy a Helix app that interfaces with a VCON (Voice Conversation) server, allowing you to analyze call data using natural language queries like:

1. Show me calls between Frank Smith and Obi Johnson
2. Display conversations about project timelines
3. Find calls from specific dates or times

## Features
- Natural language queries for call data
- Real-time call analysis
- Structured conversation data using VCON format
- Integration with Together AI

## Prerequisites
- Docker and Docker Compose
- Go 1.19 or later
- Together AI API key
- Helix CLI

## Quick Start

1. **Install Helix**
```
cd /opt
sudo mkdir helix
sudo chown $USER:$USER helix
cd helix
curl -O https://raw.githubusercontent.com/helixml/helix/main/install.sh
chmod +x install.sh
./install.sh

```

2. **Start Helix**

```
cd /opt/helix
docker compose up
```

3. **Create VCON Server**
```
mkdir -p ~/vcon-server
cd vcon-server
```

4. **Create main.go**
```
nano main.go
```
5. **Copy code and paste**

```
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
```
```
go run main.go
```

6. **Configure Helix UI**

Access Helix UI at http://localhost:8080

7. **Create new app configuration with YAML:**

```
apiVersion: app.aispec.org/v1alpha1
kind: AIApp
metadata:
  name: telco-bot
spec:
  assistants:
    - name: Default Assistant
      provider: togetherai
      model: mistralai/Mistral-7B-Instruct-v0.2
      type: text
      system_prompt: |
        1. Present call information in a clear, structured format
        2. Include relevant details like participants, timestamps, and call subjects
        3. If multiple calls are found, summarize them chronologically
        4. When showing dialog, present it in a conversational format
        5. Always mention if no calls are found matching the query criteria
      apis:
        - name: VCON Search API
          description: Search phone call records by participant name or phone number
          url: "http://host.docker.internal:8005"
          schema: |
            openapi: 3.0.0
            info:
              title: VCON Search API
              version: 1.0.0
            paths:
              /vcon:
                get:
                  summary: Get all VCONs
                  operationId: getAllVcons
                  responses:
                    '200':
                      description: List of all VCONs
                      content:
                        application/json:
                          schema:
                            type: array
                            items:
                              $ref: '#/components/schemas/VCON'
              /vcons/search:
                get:
                  summary: Search VCONs by participant
                  operationId: searchVcons
                  parameters:
                    - name: name
                      in: query
                      schema:
                        type: string
                      description: Search by participant name
                    - name: tel
                      in: query
                      schema:
                        type: string
                      description: Search by phone number
                  responses:
                    '200':
                      description: List of matching VCONs
                      content:
                        application/json:
                          schema:
                            type: array
                            items:
                              $ref: '#/components/schemas/VCON'
            components:
              schemas:
                VCON:
                  type: object
                  properties:
                    uuid:
                      type: string
                    created_at:
                      type: string
                      format: date-time
                    subject:
                      type: string
                    parties:
                      type: array
                      items:
                        type: object
                        properties:
                          name:
                            type: string
                          tel:
                            type: string
                    dialog:
                      type: array
                      items:
                        type: object
                        properties:
                          timestamp:
                            type: string
                            format: date-time
                          speaker:
                            type: string
                          text:
                            type: string
```

8. **Test this Setup

Verify VCON server is working:

```
curl http://localhost:8005/vcon
```
## Usage

Test in Helix with queries like:

- "Show me calls with Frank Smith"
- "What was discussed in the project timeline?"
- "List all calls from February 7th"

## Important Notes:

- Use host.docker.internal:8005 in Helix config (not localhost)
- VCON server must listen on 0.0.0.0:8005
- No spaces in the URL configuration
- Respect Together AI rate limits (60 RPM)
- Keep VCON server running while using Helix

## Troubleshooting:

- Check VCON server is running
- Verify URL format in Helix config
- Ensure no trailing spaces in URL
- Check Docker network connectivity
- Monitor rate limiting

## Contributing
1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License
MIT
