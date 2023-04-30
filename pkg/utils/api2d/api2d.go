package api2d

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type RequestBody struct {
    Model    string     `json:"model"`
    Messages []Message `json:"messages"`
}

func TextToSpeech() {
    url := "https://openai.api2d.net/v1/chat/completions"
    requestBody := RequestBody{
        Model: "gpt-3.5-turbo",
        Messages: []Message{
            {
                Role:    "user",
                Content: "你好！给我讲个笑话。",
            },
        },
    }
    requestBodyBytes, err := json.Marshal(requestBody)
    if err != nil {
        panic(err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
    if err != nil {
        panic(err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "fk186009-gCYVPTkf6aMycD4o2ZM9fRsDwp52ONdz|ck43-632713d")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("status:", resp.Status)
    var responseMap map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
        panic(err)
    }
    fmt.Println("response:", responseMap)
}
