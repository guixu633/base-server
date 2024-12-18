package oss

import (
	"context"
	"io"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/guixu633/base-server/module/config"
	"github.com/sirupsen/logrus"
)

type Oss struct {
	cfg *config.Oss
	oss.Bucket
}

func NewOss(cfg *config.Oss) (*Oss, error) {
	client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		logrus.WithField("err", err).Error("fail to create oss client")
		return nil, err
	}

	// connect bucket
	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		logrus.WithField("err", err).WithField("endpoint", cfg.Endpoint).WithField("bucket", cfg.Bucket).Error("fail to connect to bucket")
		return nil, err
	}
	logrus.WithField("endpoint", cfg.Endpoint).WithField("bucket", cfg.Bucket).Info("load oss success")
	return &Oss{cfg: cfg, Bucket: *bucket}, nil
}

func (o *Oss) Exists(ctx context.Context, path string) (bool, error) {
	path = strings.TrimPrefix(path, "/")

	// 先检查是否作为文件存在
	exists, err := o.IsObjectExist(path)
	if err != nil {
		logrus.WithField("err", err).WithField("path", path).Error("检查路径是否存在失败")
		return false, err
	}
	if exists {
		return true, nil
	}

	// 如果不是文件，检查是否作为目录存在
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	lsRes, err := o.ListObjects(
		oss.Prefix(path),
		oss.MaxKeys(1),
	)
	if err != nil {
		logrus.WithField("err", err).WithField("path", path).Error("检查路径是否存在失败")
		return false, err
	}

	return len(lsRes.Objects) > 0, nil
}

func (o *Oss) IsDir(ctx context.Context, path string) (bool, error) {
	path = strings.TrimPrefix(path, "/")
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// 使用 MaxKeys=1 来只获取一个对象，提高效率
	lsRes, err := o.ListObjects(
		oss.Prefix(path),
		oss.MaxKeys(1),
	)
	if err != nil {
		logrus.WithField("err", err).WithField("path", path).Error("检查路径类型失败")
		return false, err
	}

	// 如果能列出对象，说明是目录
	return len(lsRes.Objects) > 0, nil
}

func (o *Oss) GetDir(ctx context.Context, path string) ([]string, error) {
	var fileList []string
	path = strings.TrimPrefix(path, "/")
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	lsRes, err := o.ListObjects(
		oss.Prefix(path),
		oss.Delimiter("/"), // 添加分隔符参数，只列出一级目录
	)
	if err != nil {
		logrus.WithField("err", err).WithField("path", path).Error("列出文件失败")
		return nil, err
	}

	// 添加目录
	for _, prefix := range lsRes.CommonPrefixes {
		fileList = append(fileList, strings.TrimSuffix(prefix, "/"))
	}

	// 添加文件
	for _, object := range lsRes.Objects {
		if object.Key != path { // 排除路径本身
			fileList = append(fileList, object.Key)
		}
	}

	return fileList, nil
}

func (o *Oss) GetFile(ctx context.Context, path string) ([]byte, error) {
	path = strings.TrimPrefix(path, "/")
	// 检查文件是否存在
	exist, err := o.IsObjectExist(path)
	if err != nil {
		logrus.WithField("err", err).WithField("path", path).Error("检查文件是否存在失败")
		return nil, err
	}

	if !exist {
		logrus.WithField("path", path).Warn("文件不存在")
		return nil, nil
	}

	// 获取文件内容
	body, err := o.GetObject(path)
	if err != nil {
		logrus.WithField("err", err).WithField("path", path).Error("获取文件失败")
		return nil, err
	}
	defer body.Close()

	// 读取文件内容
	data, err := io.ReadAll(body)
	if err != nil {
		logrus.WithField("err", err).WithField("path", path).Error("读取文件内容失败")
		return nil, err
	}

	return data, nil
}

func (o *Oss) ListAllFilesInPath(ctx context.Context, path string) ([]string, error) {
	var fileList []string

	marker := ""
	for {
		lsRes, err := o.ListObjects(oss.Prefix(path), oss.Marker(marker))
		if err != nil {
			logrus.WithField("err", err).WithField("path", path).Error("列出文件失败")
			return nil, err
		}

		for _, object := range lsRes.Objects {
			fileList = append(fileList, object.Key)
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}

	return fileList, nil
}
