package assets

import (
	"context"

	logger "github.com/lenoobz/aws-lambda-logger"
)

// Service sector
type Service struct {
	repo Repo
	log  logger.ContextLog
}

// NewService create new service
func NewService(r Repo, l logger.ContextLog) *Service {
	return &Service{
		repo: r,
		log:  l,
	}
}

// PopulateAssets populate assets
func (s *Service) PopulateAssets(ctx context.Context) error {
	s.log.Info(ctx, "populating assets")

	overviews, err := s.repo.FindFundOverviews(ctx)
	if err != nil {
		s.log.Error(ctx, "find overviews failed", "error", err)
		return err
	}

	for _, o := range overviews {
		if err := s.repo.InsertAsset(ctx, o); err != nil {
			s.log.Error(ctx, "insert stock failed", "error", err)
		}
	}

	return nil
}
