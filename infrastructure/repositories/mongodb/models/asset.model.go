package models

import (
	"context"
	"strings"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-list/consts"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-list/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AssetModel struct
type AssetModel struct {
	ID               *primitive.ObjectID `bson:"_id,omitempty"`
	IsActive         bool                `bson:"isActive,omitempty"`
	CreatedAt        int64               `bson:"createdAt,omitempty"`
	ModifiedAt       int64               `bson:"modifiedAt,omitempty"`
	Schema           string              `bson:"schema,omitempty"`
	Source           string              `bson:"source,omitempty"`
	Ticker           string              `bson:"ticker,omitempty"`
	Name             string              `bson:"name,omitempty"`
	Type             string              `bson:"type,omitempty"`
	AssetClass       string              `bson:"assetClass,omitempty"`
	Currency         string              `bson:"currency,omitempty"`
	AllocationStock  float64             `bson:"allocationStock,omitempty"`
	AllocationBond   float64             `bson:"allocationBond,omitempty"`
	AllocationCash   float64             `bson:"allocationCash,omitempty"`
	DividendSchedule string              `bson:"dividendSchedule,omitempty"`
	Yield12Month     float64             `bson:"yield12Month,omitempty"`
	DistYield        float64             `bson:"distYield,omitempty"`
	DistAmount       float64             `bson:"distAmount,omitempty"`
}

// NewAssetModel create stock model
func NewAssetModel(ctx context.Context, l logger.ContextLog, e *entities.VanguardOverview) (*AssetModel, error) {
	var m = &AssetModel{}

	m.Source = consts.DATA_SOURCE
	m.Type = consts.SECURITY_TYPE

	m.Ticker = e.Ticker

	if e.Name != "" {
		m.Name = e.Name
	}

	if e.AssetClass != "" {
		m.AssetClass = strings.ToUpper(e.AssetClass)
	}

	if e.Currency != "" {
		m.Currency = strings.ToUpper(e.Currency)
	}

	m.AllocationStock = e.AllocationStock
	m.AllocationBond = e.AllocationBond
	m.AllocationCash = e.AllocationCash

	m.DividendSchedule = e.DividendSchedule
	m.Yield12Month = e.Yield12Month
	m.DistYield = e.DistYield
	m.DistAmount = e.DistAmount

	return m, nil
}
