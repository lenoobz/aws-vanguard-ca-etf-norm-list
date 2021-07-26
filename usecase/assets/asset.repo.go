package assets

import (
	"context"

	"github.com/lenoobz/aws-vanguard-ca-etf-norm-list/entities"
)

///////////////////////////////////////////////////////////
// Stock Repository Interface
///////////////////////////////////////////////////////////

// Reader interface
type Reader interface {
	FindFundOverviews(context.Context) ([]*entities.VanguardOverview, error)
}

// Writer interface
type Writer interface {
	InsertAsset(context.Context, *entities.VanguardOverview) error
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
