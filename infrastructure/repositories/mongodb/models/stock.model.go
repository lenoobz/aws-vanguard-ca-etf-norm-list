package models

import (
	"context"
	"time"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-vanguard-ca-etf-normalizer/consts"
	"github.com/hthl85/aws-vanguard-ca-etf-normalizer/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockModel struct {
	ID              *primitive.ObjectID      `bson:"_id,omitempty"`
	IsActive        bool                     `bson:"isActive,omitempty"`
	CreatedAt       int64                    `bson:"createdAt,omitempty"`
	ModifiedAt      int64                    `bson:"modifiedAt,omitempty"`
	Schema          string                   `bson:"schema,omitempty"`
	Source          string                   `bson:"source,omitempty"`
	Ticker          string                   `bson:"ticker,omitempty"`
	Name            string                   `bson:"name,omitempty"`
	Type            string                   `bson:"type,omitempty"`
	Sectors         []*SectorModel           `bson:"sector,omitempty"`
	Countries       []*CountryModel          `bson:"countries,omitempty"`
	DividendHistory map[int64]*DividendModel `bson:"dividendHistory,omitempty"`
}

// DividendModel struct
type DividendModel struct {
	PayoutRatio    float64    `bson:"payoutRatio,omitempty"`
	Yield          float64    `bson:"yield,omitempty"`
	Dividend       float64    `bson:"dividend,omitempty"`
	ExDividendDate *time.Time `bson:"exDividendDate,omitempty"`
	RecordDate     *time.Time `bson:"recordDate,omitempty"`
	DividendDate   *time.Time `bson:"payoutDate,omitempty"`
}

// SectorModel struct
type SectorModel struct {
	SectorCode  string  `bson:"sectorCode,omitempty"`
	SectorName  string  `bson:"sectorName,omitempty"`
	FundPercent float64 `bson:"fundPercent,omitempty"`
}

// CountryModel struct
type CountryModel struct {
	CountryCode     string  `bson:"countryCode,omitempty"`
	CountryName     string  `bson:"countryName,omitempty"`
	HoldingStatCode string  `bson:"holdingStatCode,omitempty"`
	FundMktPercent  float64 `bson:"fundMktPercent,omitempty"`
	FundTnaPercent  float64 `bson:"fundTnaPercent,omitempty"`
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

	// map countries entity to model
	var countries []*CountryModel
	for _, v := range e.Countries {
		country := &CountryModel{
			CountryCode:     v.CountryCode,
			CountryName:     v.CountryName,
			HoldingStatCode: v.HoldingStatCode,
			FundMktPercent:  v.FundMktPercent,
			FundTnaPercent:  v.FundTnaPercent,
		}

		countries = append(countries, country)
	}
	m.Countries = countries

	// map sectors entity to model
	var sectors []*SectorModel
	for _, v := range e.Sectors {
		sector := &SectorModel{
			SectorCode:  v.SectorCode,
			SectorName:  v.SectorName,
			FundPercent: v.FundPercent,
		}

		sectors = append(sectors, sector)
	}
	m.Sectors = sectors

	// map dividen history entity to model
	m.DividendHistory = make(map[int64]*DividendModel)
	for _, v := range e.DividendHistory {
		dividend := &DividendModel{
			PayoutRatio:    0,
			Yield:          0,
			Dividend:       v.Amount,
			ExDividendDate: v.AsOfDate,
			RecordDate:     v.AsOfDate,
			DividendDate:   v.AsOfDate,
		}

		dividendTime := v.AsOfDate.Unix()
		m.DividendHistory[dividendTime] = dividend
	}

	return m, nil
}
