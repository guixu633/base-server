package oss

import (
	"context"
	"fmt"
	"testing"

	"github.com/guixu633/base-server/module/config"
	"github.com/stretchr/testify/assert"
)

func TestIsDir(t *testing.T) {
	oss := getOss(t)
	isDir, err := oss.IsDir(context.Background(), "workflow/aaa")
	assert.NoError(t, err)
	fmt.Println(isDir)
}

func TestGetDir(t *testing.T) {
	oss := getOss(t)
	files, err := oss.ListAllFilesInPath(context.Background(), "workflow/tools")
	assert.NoError(t, err)
	for _, file := range files {
		fmt.Println(file)
	}
}

func TestGetObject(t *testing.T) {
	oss := getOss(t)
	object, err := oss.GetFile(context.Background(), "workflow/tools/rerank.svg")
	assert.NoError(t, err)
	fmt.Println(object)
}

func getOss(t *testing.T) *Oss {
	cfg, err := config.LoadConfig("../../config.toml")
	assert.NoError(t, err)
	oss, err := NewOss(&cfg.Oss)
	assert.NoError(t, err)
	return oss
}
