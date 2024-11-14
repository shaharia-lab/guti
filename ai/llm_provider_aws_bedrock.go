package ai

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type BedrockLLMProvider struct {
	client *bedrockruntime.Client
	model  string
}

type BedrockProviderConfig struct {
	Client *bedrockruntime.Client
	Model  string
}

func NewBedrockLLMProvider(config BedrockProviderConfig) *BedrockLLMProvider {
	if config.Model == "" {
		config.Model = "anthropic.claude-3-5-sonnet-20240620-v1:0"
	}

	return &BedrockLLMProvider{
		client: config.Client,
		model:  config.Model,
	}
}

func (p *BedrockLLMProvider) GetResponse(messages []LLMMessage, config LLMRequestConfig) (LLMResponse, error) {
	startTime := time.Now()

	var bedrockMessages []types.Message
	for _, msg := range messages {
		role := types.ConversationRoleUser
		if msg.Role == "assistant" {
			role = types.ConversationRoleAssistant
		}

		bedrockMessages = append(bedrockMessages, types.Message{
			Role: role,
			Content: []types.ContentBlock{
				&types.ContentBlockMemberText{
					Value: msg.Text,
				},
			},
		})
	}

	input := &bedrockruntime.ConverseInput{
		ModelId:  &p.model,
		Messages: bedrockMessages,
		InferenceConfig: &types.InferenceConfiguration{
			Temperature: aws.Float32(float32(config.Temperature)),
			TopP:        aws.Float32(float32(config.TopP)),
			MaxTokens:   aws.Int32(int32(config.MaxToken)),
		},
	}

	output, err := p.client.Converse(context.Background(), input)
	if err != nil {
		return LLMResponse{}, err
	}

	var responseText string
	if msgOutput, ok := output.Output.(*types.ConverseOutputMemberMessage); ok {
		for _, block := range msgOutput.Value.Content {
			if textBlock, ok := block.(*types.ContentBlockMemberText); ok {
				responseText += textBlock.Value
			}
		}
	}

	return LLMResponse{
		Text:             responseText,
		TotalInputToken:  int(*output.Usage.InputTokens),
		TotalOutputToken: int(*output.Usage.OutputTokens),
		CompletionTime:   time.Since(startTime).Seconds(),
	}, nil
}
