package repos

import (
	"context"
	"fmt"
	"strings"
	"time"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-vanguard-ca-etf-normalizer/config"
	"github.com/hthl85/aws-vanguard-ca-etf-normalizer/consts"
	"github.com/hthl85/aws-vanguard-ca-etf-normalizer/entities"
	"github.com/hthl85/aws-vanguard-ca-etf-normalizer/infrastructure/repositories/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StockMongo struct
type StockMongo struct {
	db     *mongo.Database
	client *mongo.Client
	log    logger.ContextLog
	conf   *config.MongoConfig
}

// NewStockMongo creates new stock mongo repo
func NewStockMongo(db *mongo.Database, l logger.ContextLog, conf *config.MongoConfig) (*StockMongo, error) {
	if db != nil {
		return &StockMongo{
			db:   db,
			log:  l,
			conf: conf,
		}, nil
	}

	// set context with timeout from the config
	// create new context for the query
	ctx, cancel := createContext(context.Background(), conf.TimeoutMS)
	defer cancel()

	// set mongo client options
	clientOptions := options.Client()

	// set min pool size
	if conf.MinPoolSize > 0 {
		clientOptions.SetMinPoolSize(conf.MinPoolSize)
	}

	// set max pool size
	if conf.MaxPoolSize > 0 {
		clientOptions.SetMaxPoolSize(conf.MaxPoolSize)
	}

	// set max idle time ms
	if conf.MaxIdleTimeMS > 0 {
		clientOptions.SetMaxConnIdleTime(time.Duration(conf.MaxIdleTimeMS) * time.Millisecond)
	}

	// construct a connection string from mongo config object
	cxnString := fmt.Sprintf("mongodb+srv://%s:%s@%s", conf.Username, conf.Password, conf.Host)

	// create mongo client by making new connection
	client, err := mongo.Connect(ctx, clientOptions.ApplyURI(cxnString))
	if err != nil {
		return nil, err
	}

	return &StockMongo{
		db:     client.Database(conf.Dbname),
		client: client,
		log:    l,
		conf:   conf,
	}, nil
}

// Close disconnect from database
func (r *StockMongo) Close() {
	ctx := context.Background()
	r.log.Info(ctx, "close mongo client")

	if r.client == nil {
		return
	}

	if err := r.client.Disconnect(ctx); err != nil {
		r.log.Error(ctx, "disconnect mongo failed", "error", err)
	}
}

///////////////////////////////////////////////////////////
// Implement repo interface
///////////////////////////////////////////////////////////

// FindOverviews finds all fund overviews
func (r *StockMongo) FindOverviews(ctx context.Context) ([]*entities.VanguardOverview, error) {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.VANGUARD_OVERVIEW_COL]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return nil, fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	// filter
	filter := bson.D{}

	// find options
	findOptions := options.Find()

	cur, err := col.Find(ctx, filter, findOptions)

	// only run defer function when find success
	if cur != nil {
		defer func() {
			if deferErr := cur.Close(ctx); deferErr != nil {
				err = deferErr
			}
		}()
	}

	// find was not succeed
	if err != nil {
		r.log.Error(ctx, "find query failed", "error", err)
		return nil, err
	}

	var funds []*entities.VanguardOverview

	// iterate over the cursor to decode document one at a time
	for cur.Next(ctx) {
		// decode cursor to activity model
		var fund entities.VanguardOverview
		if err = cur.Decode(&fund); err != nil {
			r.log.Error(ctx, "decode failed", "error", err)
			return nil, err
		}

		funds = append(funds, &fund)
	}

	if err := cur.Err(); err != nil {
		r.log.Error(ctx, "iterate over cursor failed", "error", err)
		return nil, err
	}

	return funds, nil
}

// InsertStock insert new stock stock
func (r *StockMongo) InsertStock(ctx context.Context, fund *entities.VanguardOverview) error {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	savedStock, err := r.findStockByTicker(ctx, fund.Ticker)
	if err != nil {
		r.log.Error(ctx, "find stock by ticker failed", "error", err, "ticker", fund.Ticker)
		return err
	}

	insertingStock, err := models.NewStockModel(ctx, r.log, fund)
	if err != nil {
		r.log.Error(ctx, "create model failed", "error", err, "ticker", fund.Ticker)
		return err
	}

	if savedStock != nil {
		// Copy dividend history from saved stock to inserting stock
		for k, v := range savedStock.DividendHistory {
			if _, f := insertingStock.DividendHistory[k]; !f {
				insertingStock.DividendHistory[k] = v
			}
		}
	}

	if err = r.insertStock(ctx, insertingStock); err != nil {
		r.log.Error(ctx, "insert stock failed", "error", err, "ticker", fund.Ticker)
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////
// Implement helper function
///////////////////////////////////////////////////////////

// createContext create a new context with timeout
func createContext(ctx context.Context, t uint64) (context.Context, context.CancelFunc) {
	timeout := time.Duration(t) * time.Millisecond
	return context.WithTimeout(ctx, timeout*time.Millisecond)
}

// insertStock inserts new stock
func (r *StockMongo) insertStock(ctx context.Context, m *models.StockModel) error {
	if m == nil {
		r.log.Error(ctx, "invalid param")
		return fmt.Errorf("invalid param")
	}

	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.VANGUARD_FUND_COL]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	m.IsActive = true
	m.Schema = r.conf.SchemaVersion
	m.ModifiedAt = time.Now().UTC().Unix()

	filter := bson.D{{
		Key:   "ticker",
		Value: m.Ticker,
	}}

	update := bson.D{
		{
			Key:   "$set",
			Value: m,
		},
		{
			Key: "$setOnInsert",
			Value: bson.D{{
				Key:   "createdAt",
				Value: time.Now().UTC().Unix(),
			}},
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.log.Error(ctx, "update one failed", "error", err)
		return err
	}

	return nil
}

// findStockByTicker finds stock of a given ticker
func (r *StockMongo) findStockByTicker(ctx context.Context, ticker string) (*models.StockModel, error) {
	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.VANGUARD_FUND_COL]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return nil, fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	// filter
	filter := bson.D{
		{
			Key:   "ticker",
			Value: strings.ToUpper(ticker),
		},
	}

	// find options
	findOptions := options.FindOne()

	var stock models.StockModel
	if err := col.FindOne(ctx, filter, findOptions).Decode(&stock); err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			r.log.Info(ctx, "stock not found", "ticker", ticker)
			return nil, nil
		}

		r.log.Error(ctx, "decode find one failed", "error", err, "ticker", ticker)
		return nil, err
	}

	return &stock, nil
}
