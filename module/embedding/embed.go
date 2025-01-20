package embedding

import (
	"math"
	"net/http"

	"github.com/guixu633/base-server/module/config"
)

type EmbedEngine interface {
	Embed(text string, query bool) ([]float32, error)
	EmbedBatch(texts []string, query bool) ([][]float32, error)
}

func NewEmbedEngine(cfg *config.Embedding, client *http.Client) EmbedEngine {
	switch cfg.Engine {
	case "openai":
		return &OpenaiEmbedEngine{
			cfg:    cfg,
			client: client,
		}
	case "qwen":
		return &QwenEmbedEngine{
			cfg:    cfg,
			client: client,
		}
	}
	return nil
}

func ConsineSimilarity(embedding1, embedding2 []float32) float32 {
	dotProduct := float32(0)
	for i := range embedding1 {
		dotProduct += embedding1[i] * embedding2[i]
	}
	return dotProduct / (L2Norm(embedding1) * L2Norm(embedding2))
}

func L2Norm(embedding []float32) float32 {
	sum := float32(0)
	for _, value := range embedding {
		sum += value * value
	}
	return float32(math.Sqrt(float64(sum)))
}
