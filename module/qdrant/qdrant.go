package qdrant

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guixu633/base-server/module/config"
	"github.com/guixu633/base-server/module/embedding"
	"github.com/qdrant/go-client/qdrant"
	"github.com/sirupsen/logrus"
)

type Qdrant struct {
	cfg              *config.Qdrant
	cfgCryptoArticle *config.CryptoArticle
	client           *qdrant.Client
	embedEngine      *embedding.EmbedEngine
}

func NewClient(cfg *config.Qdrant, cfgCryptoArticle *config.CryptoArticle, embedEngine *embedding.EmbedEngine) (*Qdrant, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   cfg.Host,
		Port:   cfg.Port,
		APIKey: cfg.ApiKey,
		UseTLS: true,
	})
	if err != nil {
		return nil, err
	}
	return &Qdrant{
		cfg:              cfg,
		cfgCryptoArticle: cfgCryptoArticle,
		client:           client,
		embedEngine:      embedEngine,
	}, nil
}

func (q *Qdrant) InitVdb(ctx context.Context) error {
	err := q.InitCryptoArticleVDB()
	if err != nil {
		return err
	}
	return nil
}

func (q *Qdrant) ListCollection(ctx context.Context) ([]string, error) {
	return q.client.ListCollections(ctx)
}

func (q *Qdrant) CreateCollection(ctx context.Context, collectionName string) error {
	return q.client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(q.cfg.Dimension),
			Distance: qdrant.Distance_Cosine,
		}),
	})
}

func (q *Qdrant) ExistsCollection(ctx context.Context, collectionName string) (bool, error) {
	return q.client.CollectionExists(ctx, collectionName)
}

func CreateIfNotExists(q *Qdrant, collectionName string) error {
	exist, err := q.ExistsCollection(context.Background(), collectionName)
	if err != nil {
		return err
	}
	if exist {
		logrus.WithField("collection", collectionName).Info("collection already exists")
		return nil
	}
	logrus.WithField("collection", collectionName).Info("collection not exists, creating...")
	err = q.CreateCollection(context.Background(), collectionName)
	if err != nil {
		return err
	}
	logrus.WithField("collection", collectionName).Info("collection created")
	return nil
}

func (q *Qdrant) DeleteCollection(ctx context.Context, collectionName string) error {
	return q.client.DeleteCollection(ctx, collectionName)
}

func (q *Qdrant) Upsert(ctx context.Context, collectionName string, id, data string, params map[string]any) error {
	embeddings, err := q.embedEngine.Embed(data, false)
	if err != nil {
		logrus.WithField("error", err).Error("embed data failed")
		return err
	}

	if params == nil {
		params = make(map[string]any)
	}
	params["content"] = data

	points := make([]*qdrant.PointStruct, 0)
	if id == "" {
		id = uuid.New().String()
	}
	points = append(points, &qdrant.PointStruct{
		Id:      qdrant.NewIDUUID(id),
		Vectors: qdrant.NewVectors(embeddings...),
		Payload: qdrant.NewValueMap(params),
	})

	logrus.WithField("points", points).Debug("preparing to upsert points")

	operationInfo, err := q.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         points,
	})
	if err != nil {
		logrus.WithField("collection_name", collectionName).WithField("error", err).Error("upsert data failed")
		return err
	}

	logrus.WithField("collection_name", collectionName).WithField("data", data).WithField("params", params).WithField("info", operationInfo).Info("upsert data success")
	return nil
}

func (q *Qdrant) Search(ctx context.Context, collectionName string, query string, scoreThreshold float32, limit uint64, filter *qdrant.Filter) ([]*qdrant.ScoredPoint, error) {
	embeddings, err := q.embedEngine.Embed(query, true)
	if err != nil {
		logrus.WithField("error", err).Error("embed query failed")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"collection":       collectionName,
		"query":            query,
		"embedding_length": len(embeddings),
	}).Debug("开始搜索")

	queryPoint := &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(embeddings...),
		WithPayload: &qdrant.WithPayloadSelector{
			SelectorOptions: &qdrant.WithPayloadSelector_Enable{
				Enable: true,
			},
		},
		Filter: filter,
	}

	if scoreThreshold > 0 {
		queryPoint.ScoreThreshold = &scoreThreshold
	}

	if limit > 0 {
		queryPoint.Limit = &limit
	}

	searchResult, err := q.client.Query(ctx, queryPoint)

	if err != nil {
		logrus.WithError(err).Error("搜索失败")
		return nil, err
	}

	logrus.WithField("results_count", len(searchResult)).Debug("搜索完成")

	return searchResult, nil
}

func (q *Qdrant) GetPoint(ctx context.Context, collectionName string, id string) (*qdrant.RetrievedPoint, error) {
	client := q.client.GetPointsClient()
	response, err := client.Get(ctx, &qdrant.GetPoints{
		CollectionName: collectionName,
		Ids:            []*qdrant.PointId{qdrant.NewIDUUID(id)},
		WithPayload: &qdrant.WithPayloadSelector{
			SelectorOptions: &qdrant.WithPayloadSelector_Enable{
				Enable: true,
			},
		},
		WithVectors: &qdrant.WithVectorsSelector{
			SelectorOptions: &qdrant.WithVectorsSelector_Enable{
				Enable: true,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(response.GetResult()) == 0 {
		return nil, fmt.Errorf("point not found")
	}
	return response.GetResult()[0], nil
}
