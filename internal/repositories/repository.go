package repositories

import (
	"context"
	"time"

	"github.com/Armunz/learn-mongodb/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, account entity.Account) error
	List(ctx context.Context, product string, orderBy int, limit int, offset int) ([]entity.Account, int64, error)
	GetByAccountID(ctx context.Context, accountID int) (entity.Account, error)
	Update(ctx context.Context, account entity.Account) error
	Delete(ctx context.Context, accountID int) error
}

type repoImpl struct {
	collection *mongo.Collection
	timeoutMs  int
}

func New(database *mongo.Database, timeoutMs int) Repository {
	return &repoImpl{
		collection: database.Collection(ACCOUNTS_COLLECTION_NAME),
		timeoutMs:  timeoutMs,
	}
}

// Create implements Repository.
func (r *repoImpl) Create(ctx context.Context, account entity.Account) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	_, err := r.collection.InsertOne(ctxTimeout, account)

	return err
}

// Delete implements Repository.
func (r *repoImpl) Delete(ctx context.Context, accountID int) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	filter := bson.M{"account_id": accountID}
	_, err := r.collection.DeleteOne(ctxTimeout, filter)

	return err
}

// GetByAccountID implements Repository.
func (r *repoImpl) GetByAccountID(ctx context.Context, accountID int) (entity.Account, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	var account entity.Account
	filter := bson.M{"account_id": accountID}
	err := r.collection.FindOne(ctxTimeout, filter).Decode(&account)

	return account, err
}

// List implements Repository.
func (r *repoImpl) List(ctx context.Context, product string, orderBy int, limit int, offset int) ([]entity.Account, int64, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	// build pipeline
	var metadataStage bson.A
	var dataStage bson.A

	// match stage
	if product != "" {
		dataStage = append(dataStage, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "products", Value: product}}}})
		metadataStage = append(metadataStage, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "products", Value: product}}}})
	}

	// sort stage
	if orderBy != 0 {
		dataStage = append(dataStage, bson.D{primitive.E{Key: "$sort", Value: bson.D{primitive.E{Key: "account_id", Value: orderBy}}}})
	}

	// pagination stage
	dataStage = append(dataStage, bson.D{primitive.E{Key: "$skip", Value: offset}})
	dataStage = append(dataStage, bson.D{primitive.E{Key: "$limit", Value: limit}})

	// count stage
	metadataStage = append(metadataStage, bson.D{primitive.E{Key: "$count", Value: "total_count"}})

	// facet stage
	facet := bson.M{
		"metadata": metadataStage,
		"data":     dataStage,
	}

	pipeline := mongo.Pipeline{
		{primitive.E{Key: "$facet", Value: facet}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err := cursor.All(ctxTimeout, &results); err != nil {
		return nil, 0, err
	}

	var totalCount int64
	var accounts []entity.Account
	if len(results) > 0 {
		metadata, ok := results[0]["metadata"].(bson.A)[0].(bson.M)
		if !ok {
			return nil, 0, errMetaDataTypeAssertion
		}

		data, ok := results[0]["data"].(bson.A)
		if !ok {
			return nil, 0, errDataTypeAssertion
		}

		totalCount = int64(metadata["total_count"].(int32))
		for _, d := range data {
			ac, ok := d.(bson.M)
			if !ok {
				return nil, 0, errAccountsTypeAssertion
			}

			acProducts := ac["products"].(bson.A)
			products := make([]string, 0, len(acProducts))
			for _, p := range acProducts {
				products = append(products, p.(string))
			}

			account := entity.Account{
				AccountID: int(ac["account_id"].(int32)),
				Limit:     int(ac["limit"].(int32)),
				Products:  products,
			}

			accounts = append(accounts, account)
		}
	}

	return accounts, totalCount, nil
}

// Update implements Repository.
func (r *repoImpl) Update(ctx context.Context, account entity.Account) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	filter := bson.M{"account_id": account.AccountID}
	_, err := r.collection.UpdateOne(ctxTimeout, filter, bson.M{"$set": account})

	return err
}
