package qdrant

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/qdrant/go-client/qdrant"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CryptoArticle struct {
	Id          string  `json:"id"`
	Score       float64 `json:"score"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	PublishTime string  `json:"publish_time"`
	Src         string  `json:"src"`
	Url         string  `json:"url"`
}

func (c *CryptoArticle) ToRetrievalRecord() RetrievalRecord {
	return RetrievalRecord{
		Content: c.Content,
		Score:   c.Score,
		Title:   c.Title,
		Metadata: map[string]any{
			"id":           c.Id,
			"publish_time": c.PublishTime,
			"src":          c.Src,
			"url":          c.Url,
		},
	}
}

func (q *Qdrant) GetCryptoArticleInfo() (*qdrant.CollectionInfo, error) {
	info, err := q.client.GetCollectionInfo(context.Background(), q.cfgCryptoArticle.Collection)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (q *Qdrant) InitCryptoArticleVDB() error {
	collectionName := q.cfgCryptoArticle.Collection
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
		logrus.WithField("error", err).Error("create collection failed")
		return err
	}
	logrus.WithField("collection", collectionName).Info("collection created")
	_, err = q.client.CreateFieldIndex(context.Background(), &qdrant.CreateFieldIndexCollection{
		CollectionName: collectionName,
		FieldName:      "publish_time",
		FieldType:      qdrant.FieldType_FieldTypeDatetime.Enum(),
	})
	if err != nil {
		logrus.WithField("error", err).Error("create time index failed")
		return err
	}
	logrus.WithField("collection", collectionName).Info("init time index success")
	return nil
}

func (q *Qdrant) DeleteCryptoArticleVDB() error {
	return q.client.DeleteCollection(context.Background(), q.cfgCryptoArticle.Collection)
}

func (q *Qdrant) UpsertCryptoArticle(ctx context.Context, article *CryptoArticle) error {
	return q.Upsert(ctx, q.cfgCryptoArticle.Collection, article.Id, article.Content, map[string]any{
		"title":        article.Title,
		"publish_time": article.PublishTime,
		"src":          article.Src,
		"url":          article.Url,
	})
}

func (q *Qdrant) SearchCryptoArticle(ctx context.Context, query string, topk int, scoreThreshold float64) ([]*CryptoArticle, error) {
	startTime := time.Now().Add(-time.Duration(q.cfgCryptoArticle.HoursLimit) * time.Hour)
	endTime := time.Now()

	filter := &qdrant.Filter{
		Must: []*qdrant.Condition{
			qdrant.NewDatetimeRange("publish_time", &qdrant.DatetimeRange{
				Gte: timestamppb.New(startTime),
				Lte: timestamppb.New(endTime),
			}),
		},
	}

	results, err := q.Search(ctx, q.cfgCryptoArticle.Collection, query, float32(scoreThreshold), uint64(topk), filter)
	if err != nil {
		return nil, err
	}
	data := make([]*CryptoArticle, 0, len(results))
	for _, result := range results {
		data = append(data, &CryptoArticle{
			Id:          result.Id.GetUuid(),
			Score:       float64(result.Score),
			Title:       result.Payload["title"].GetStringValue(),
			Content:     result.Payload["content"].GetStringValue(),
			PublishTime: result.Payload["publish_time"].GetStringValue(),
			Src:         result.Payload["src"].GetStringValue(),
			Url:         result.Payload["url"].GetStringValue(),
		})
	}

	return data, nil
}

func (q *Qdrant) SearchCryptoArticleStr(ctx context.Context, query string, topk int, scoreThreshold float64) (string, error) {
	data, err := q.SearchCryptoArticle(ctx, query, topk, scoreThreshold)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for i, article := range data {
		fmt.Fprintf(&sb, "参考文章 %d\n", i+1)
		fmt.Fprintf(&sb, "标题: %s\n", article.Title)
		fmt.Fprintf(&sb, "时间: %s\n", article.PublishTime)
		fmt.Fprintf(&sb, "正文: %s\n", article.Content)
		sb.WriteString("---\n")
	}
	return sb.String(), nil
}

func (q *Qdrant) RetrievalCryptoArticle(ctx context.Context, req RetrievalRequest) (RetrievalResponse, error) {
	articles, err := q.SearchCryptoArticle(ctx, req.Query, req.Settings.TopK, req.Settings.ScoreThreshold)
	if err != nil {
		return RetrievalResponse{}, err
	}
	records := make([]RetrievalRecord, 0, len(articles))
	for _, article := range articles {
		records = append(records, article.ToRetrievalRecord())
	}
	return RetrievalResponse{Records: records}, nil
}
