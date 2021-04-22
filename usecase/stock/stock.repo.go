package stock

import (
	"context"

	"github.com/hthl85/aws-vanguard-ca-etf-normalizer/entities"
)

///////////////////////////////////////////////////////////
// Stock Repository Interface
///////////////////////////////////////////////////////////

// Reader interface
type Reader interface {
	FindOverviews(context.Context) ([]*entities.VanguardOverview, error)
}

// Writer interface
type Writer interface {
	InsertStock(context.Context, *entities.VanguardOverview) error
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
