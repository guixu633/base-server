package db

import (
	"context"
	"fmt"
)

type News struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (d *Database) ListNews(ctx context.Context) ([]*News, error) {
	query := `
		SELECT id, content, created_at, updated_at
		FROM news
		ORDER BY created_at DESC
	`

	rows, err := d.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("查询新闻列表失败: %v", err)
	}
	defer rows.Close()

	var newsList []*News
	for rows.Next() {
		news := &News{}
		err := rows.Scan(&news.ID, &news.Content, &news.CreatedAt, &news.UpdatedAt)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("扫描新闻数据失败: %v", err)
		}
		newsList = append(newsList, news)
	}

	if err = rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("遍历新闻数据时发生错误: %v", err)
	}

	if err = rows.Close(); err != nil {
		return nil, fmt.Errorf("关闭数据库连接时发生错误: %v", err)
	}

	return newsList, nil
}

func (d *Database) GetLatestNews(ctx context.Context) (string, error) {
	query := `
		SELECT id, content, created_at, updated_at
		FROM news
		ORDER BY created_at DESC
		LIMIT 1
	`

	row := d.QueryRowContext(ctx, query)
	news := &News{}
	err := row.Scan(&news.ID, &news.Content, &news.CreatedAt, &news.UpdatedAt)
	if err != nil {
		return "", fmt.Errorf("获取最新新闻失败: %v", err)
	}
	return news.Content, nil
}
