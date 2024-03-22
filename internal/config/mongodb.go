package config

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongo(ctx context.Context, cfg Config) *mongo.Database {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.AppMongoInitConnectionTimeSecond)*time.Second)
	defer cancel()

	option := options.Client().
		ApplyURI(cfg.AppMongoURI).
		SetMinPoolSize(uint64(cfg.AppMongoPoolMin)).
		SetMaxPoolSize(uint64(cfg.AppMongoPoolMax)).
		SetMaxConnIdleTime(time.Duration(cfg.AppMongoMaxIdleTimeSecond) * time.Second)

	client, err := mongo.Connect(ctx, option)
	if err != nil {
		log.Panic().Err(err).Msg("failed to connect to mongoDB")
	}

	database := client.Database(cfg.AppMongoDatabaseName)
	if err := database.Client().Ping(context.Background(), readpref.Primary()); err != nil {
		log.Panic().Err(err).Msg("failed to ping mongoDB")
	}

	return database
}
