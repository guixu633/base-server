package service

import "context"

func (s *Service) Exists(ctx context.Context, path string) (bool, error) {
	return s.oss.Exists(ctx, path)
}

func (s *Service) IsDir(ctx context.Context, path string) (bool, error) {
	return s.oss.IsDir(ctx, path)
}

func (s *Service) GetFile(ctx context.Context, path string) ([]byte, error) {
	return s.oss.GetFile(ctx, path)
}

func (s *Service) GetDir(ctx context.Context, path string) ([]string, error) {
	return s.oss.GetDir(ctx, path)
}
