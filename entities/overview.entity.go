package entities

import (
	"time"
)

// VanguardOverview struct
type VanguardOverview struct {
	PortID           string              `json:"portId,omitempty"`
	AssetClass       string              `json:"assetClass,omitempty"`
	Strategy         string              `json:"strategy,omitempty"`
	DividendSchedule string              `json:"dividendSchedule,omitempty"`
	Name             string              `json:"name,omitempty"`
	Currency         string              `json:"currency,omitempty"`
	Isin             string              `json:"isin,omitempty"`
	Sedol            string              `json:"sedol,omitempty"`
	Ticker           string              `json:"ticker,omitempty"`
	TotalAssets      float64             `json:"totalAssets,omitempty"`
	Yield12Month     float64             `json:"yield12Month,omitempty"`
	Price            float64             `json:"price,omitempty"`
	ManagementFee    float64             `json:"managementFee,omitempty"`
	MerFee           float64             `json:"merFee,omitempty"`
	DistYield        float64             `json:"distYield,omitempty"`
	AllocationStock  float64             `json:"allocationStock,omitempty"`
	AllocationBond   float64             `json:"allocationBond,omitempty"`
	AllocationCash   float64             `json:"allocationCash,omitempty"`
	Sectors          []*SectorBreakdown  `json:"sectors,omitempty"`
	Countries        []*CountryBreakdown `json:"countries,omitempty"`
	DividendHistory  []*DividendHistory  `json:"dividendHistory,omitempty"`
}

// SectorBreakdown struct
type SectorBreakdown struct {
	SectorCode  string  `json:"sectorCode,omitempty"`
	SectorName  string  `json:"sectorName,omitempty"`
	FundPercent float64 `json:"fundPercent,omitempty"`
}

// CountryBreakdown struct
type CountryBreakdown struct {
	CountryCode     string  `json:"countryCode,omitempty"`
	CountryName     string  `json:"countryName,omitempty"`
	FundMktPercent  float64 `json:"fundMktPercent,omitempty"`
	FundTnaPercent  float64 `json:"fundTnaPercent,omitempty"`
	HoldingStatCode string  `json:"holdingStatCode,omitempty"`
}

// DividendHistory struct
type DividendHistory struct {
	Amount       float64    `json:"amount,omitempty"`
	CurrencyCode string     `json:"currencyCode,omitempty"`
	AsOfDate     *time.Time `json:"asOfDate,omitempty"`
}
