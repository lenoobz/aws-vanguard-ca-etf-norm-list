package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	corid "github.com/hthl85/aws-lambda-corid"
	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-list/config"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-list/infrastructure/repositories/mongodb/repos"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-list/usecase/assets"
)

func main() {
	appConf := config.AppConf

	// create new logger
	zap, err := logger.NewZapLogger()
	if err != nil {
		log.Fatal("create app logger failed")
	}
	defer zap.Close()

	// create new repository
	repo, err := repos.NewAssetMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatal("create asset mongo repo failed")
	}
	defer repo.Close()

	// create new service
	svc := assets.NewService(repo, zap)

	// try correlation context
	id, _ := uuid.NewRandom()
	ctx := corid.NewContext(context.Background(), id)

	err = svc.PopulateAssets(ctx)
	if err != nil {
		log.Fatal("populate failed")
	}
}
