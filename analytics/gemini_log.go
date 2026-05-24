package analytics

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"log"
)

type ClientModel struct {
	ctx    context.Context
	client *genai.Client
}

func NewAIClientModel() *ClientModel {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &ClientModel{
		ctx:    ctx,
		client: client,
	}
}

func (cm *ClientModel) QueryModel(query string, schema *genai.Schema) string {
	result, err := cm.client.Models.GenerateContent(
		cm.ctx,
		"gemini-2.5-flash",
		genai.Text(query),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema:   schema,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Text())
	return result.Text()
}
