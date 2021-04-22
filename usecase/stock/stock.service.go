package stock

import (
	"context"

	logger "github.com/hthl85/aws-lambda-logger"
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

// PopulateStocks populate stocks data
func (s *Service) PopulateStocks(ctx context.Context) error {
	s.log.Info(ctx, "populating stocks")

	overviews, err := s.repo.FindOverviews(ctx)
	if err != nil {
		s.log.Error(ctx, "find overviews failed", "error", err)
		return err
	}

	for _, o := range overviews {
		if err := s.repo.InsertStock(ctx, o); err != nil {
			s.log.Error(ctx, "insert stock failed", "error", err)
		}
	}

	return nil
}
