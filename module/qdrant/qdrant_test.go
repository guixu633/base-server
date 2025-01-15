package qdrant

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/guixu633/base-server/module/config"
	"github.com/guixu633/base-server/module/embedding"
	"github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestListCollection(t *testing.T) {
	qdrant := getQdrant(t)
	collections, err := qdrant.ListCollection(context.Background())
	assert.NoError(t, err)
	for _, collection := range collections {
		fmt.Println(collection)
	}
	fmt.Println(len(collections))
}

func TestExistsCollection(t *testing.T) {
	qdrant := getQdrant(t)
	exists, err := qdrant.ExistsCollection(context.Background(), "test")
	assert.NoError(t, err)
	fmt.Println(exists)
}

func TestCreateCollection(t *testing.T) {
	qdrant := getQdrant(t)
	err := qdrant.CreateCollection(context.Background(), "test")
	assert.NoError(t, err)
}

func TestUpsert(t *testing.T) {
	vdb := getQdrant(t)
	title := "哈哈哈"
	data := "人工智能技术正在快速发展，深度学习模型能够处理越来越复杂的任务。从图像识别到自然语言处理，AI 系统展现出惊人的能力。特别是在医疗诊断、自动驾驶和智能助手等领域，AI 已经开始改变我们的生活方式。"
	publishTime := time.Now().Unix()
	// RFC 3339 format
	date := time.Unix(publishTime, 0).Format(time.RFC3339)
	// timestamp := timestamppb.New(time.Unix(publishTime, 0))
	err := vdb.Upsert(context.Background(), "test", "", data, map[string]any{
		"title":        title,
		"publish_time": date,
	})
	assert.NoError(t, err)
}

func TestGetPoint(t *testing.T) {
	qdrant := getQdrant(t)
	point, err := qdrant.GetPoint(context.Background(), "test", "517c6d82-b058-41e9-98a5-6e2be63d3172")
	assert.NoError(t, err)
	fmt.Println(point)
	fmt.Println(point.Payload)
}

func TestSearch(t *testing.T) {
	qdrant := getQdrant(t)
	results, err := qdrant.Search(context.Background(), "test", "人工智能", 0, 0, nil)
	assert.NoError(t, err)

	if len(results) == 0 {
		t.Skip("no results")
	}

	for _, result := range results[:] {
		title := result.Payload["title"].GetStringValue()
		content := result.Payload["content"].GetStringValue()
		fmt.Println(result.Score)
		fmt.Printf("标题: %s\n", title)
		fmt.Printf("内容: %s\n", content)
		fmt.Println(result.Payload["publish_time"].GetStringValue())
		fmt.Println("--------------------------------")
	}
}

func TestSearchWithTimeFilter(t *testing.T) {
	vdb := getQdrant(t)

	startTime := time.Date(2025, 1, 13, 11, 14, 55, 0, time.Local)
	// startTime := time.Now().Add(-240 * time.Hour)
	endTime := time.Now()

	filter := &qdrant.Filter{
		Must: []*qdrant.Condition{
			qdrant.NewDatetimeRange("publish_time", &qdrant.DatetimeRange{
				Gte: timestamppb.New(startTime),
				Lte: timestamppb.New(endTime),
			}),
		},
	}

	results, err := vdb.Search(context.Background(), "test", "人工智能", 0.65, 10, filter)
	assert.NoError(t, err)

	for _, result := range results {
		title := result.Payload["title"].GetStringValue()
		publishTime := result.Payload["publish_time"].GetStringValue()
		fmt.Printf("标题: %s, 发布时间: %s\n", title, publishTime)
	}
}

func getQdrant(t *testing.T) *Qdrant {
	cfg, err := config.LoadConfig("../../config.toml")
	assert.NoError(t, err)
	embedEngine := embedding.NewEmbedEngine(&cfg.Embedding, &http.Client{})
	qdrant, err := NewClient(&cfg.Qdrant, &cfg.CryptoArticle, embedEngine)
	assert.NoError(t, err)
	return qdrant
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now().Format(time.RFC3339))
	timeStr := "2024-01-15T10:00:00+08:00"
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		t.Errorf("解析时间失败: %v", err)
		return
	}
	fmt.Printf("成功解析时间: %v\n", parsedTime)

	invalidTimeStr := "invalid-time"
	_, err = time.Parse(time.RFC3339, invalidTimeStr)
	if err == nil {
		t.Error("应该解析失败但成功了")
	}
}
