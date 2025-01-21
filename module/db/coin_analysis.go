package db

import (
	"context"
	"fmt"
)

type CoinAnalysis struct {
	ID        string `json:"id"`
	CoinID    string `json:"coin_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (d *Database) GetAllCoinAnalysis(ctx context.Context) ([]*CoinAnalysis, error) {
	query := `
		SELECT id, coin_id, content, created_at, updated_at
		FROM analysis
	`

	rows, err := d.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("获取所有币种分析失败: %v", err)
	}
	defer rows.Close()

	analysisList := make([]*CoinAnalysis, 0)
	for rows.Next() {
		analysis := &CoinAnalysis{}
		err := rows.Scan(&analysis.ID, &analysis.CoinID, &analysis.Content, &analysis.CreatedAt, &analysis.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("获取所有币种分析失败: %v", err)
		}
		analysisList = append(analysisList, analysis)
	}

	return analysisList, nil
}

func (d *Database) GetLatestCoinAnalysis(ctx context.Context, coinID string) (string, error) {
	query := `
		SELECT id, coin_id, content, created_at, updated_at
		FROM analysis
		WHERE coin_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	row := d.QueryRowContext(ctx, query, coinID)
	analysis := &CoinAnalysis{}
	err := row.Scan(&analysis.ID, &analysis.CoinID, &analysis.Content, &analysis.CreatedAt, &analysis.UpdatedAt)
	if err != nil {
		return "", fmt.Errorf("获取最新币种分析失败: %v", err)
	}
	return analysis.Content, nil
}
