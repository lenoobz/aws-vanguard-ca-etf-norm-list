package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
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

	lambda.Start(svc.PopulateAssets)
}
