package db

import (
	"database/sql"
	"time"

	"github.com/guixu633/base-server/module/config"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Database struct {
	cfg config.Postgres
	*sql.DB
}

func NewDB(cfg config.Postgres) (*Database, error) {
	// 直接构造连接字符串并用单引号包裹密码
	// connStr := fmt.Sprintf("host=%s port=%d user=%s password='%s' dbname=%s sslmode=disable",
	// 	cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)
	// connStr := "postgresql://root:UCA=seIJbDw@dd2d5ee596f1458d8a8c0870246b895ein03.internal.cn-north-9.postgresql.rds.myhuaweicloud.com/xiaoxiang-content"
	connStr := "postgresql://pgroot:%232023INF888@xiaoxiang-blockchain.rwlb.singapore.rds.aliyuncs.com:5432/tokensense-data?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logrus.WithError(err).WithField("url", connStr).Error("Failed to connect to the database")
		return nil, err
	}

	// 连接池设置
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 验证连接
	if err = db.Ping(); err != nil {
		logrus.WithError(err).Error("Database ping failed")
		return nil, err
	}

	logrus.Info("Successfully connected to the database")
	return &Database{DB: db, cfg: cfg}, nil
}
