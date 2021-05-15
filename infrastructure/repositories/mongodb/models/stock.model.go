package models

import (
	"context"
	"strings"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-vanguard-ca-etf-normalizer/consts"
	"github.com/hthl85/aws-vanguard-ca-etf-normalizer/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StockModel struct
type StockModel struct {
	ID              *primitive.ObjectID `bson:"_id,omitempty"`
	IsActive        bool                `bson:"isActive,omitempty"`
	CreatedAt       int64               `bson:"createdAt,omitempty"`
	ModifiedAt      int64               `bson:"modifiedAt,omitempty"`
	Schema          string              `bson:"schema,omitempty"`
	Source          string              `bson:"source,omitempty"`
	Ticker          string              `bson:"ticker,omitempty"`
	Name            string              `bson:"name,omitempty"`
	Type            string              `bson:"type,omitempty"`
	AssetClass      string              `bson:"assetClass,omitempty"`
	Currency        string              `bson:"currency,omitempty"`
	AllocationStock float64             `bson:"allocationStock,omitempty"`
	AllocationBond  float64             `bson:"allocationBond,omitempty"`
	AllocationCash  float64             `bson:"allocationCash,omitempty"`
}

// NewStockModel create stock model
func NewStockModel(ctx context.Context, l logger.ContextLog, e *entities.VanguardOverview) (*StockModel, error) {
	var m = &StockModel{}

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

	return m, nil
}
