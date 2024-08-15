package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "os"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Estrutura para a resposta do Confluent
type TopicResponse struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func createTopic(ctx *pulumi.Context, name string) (string, error) {
    apiKey := os.Getenv("CONFLUENT_API_KEY")
    apiSecret := os.Getenv("CONFLUENT_API_SECRET")

    // Configuração do payload para criação do tópico
    topicData := fmt.Sprintf(`{
        "name": "%s",
        "partitions": 3,
        "replication_factor": 1
    }`, name)

    req, err := http.NewRequest("POST", "https://api.confluent.cloud/v1/kafka/v1/topics", strings.NewReader(topicData))
    if err != nil {
        return "", err
    }

    req.SetBasicAuth(apiKey, apiSecret)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var topicResponse TopicResponse
    if err := json.NewDecoder(resp.Body).Decode(&topicResponse); err != nil {
        return "", err
    }

    return topicResponse.ID, nil
}

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        topicID, err := createTopic(ctx, "my-pulumi-go-topic")
        if err != nil {
            return err
        }

        ctx.Export("topicId", pulumi.String(topicID))

        return nil
    })
}



