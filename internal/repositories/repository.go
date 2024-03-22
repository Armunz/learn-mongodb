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
	pipeline := mongo.Pipeline{}

	// match stage
	if product != "" {
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "products", Value: product}}}})
	}

	// count stage
	pipeline = append(pipeline, bson.D{primitive.E{Key: "$count", Value: "total_count"}})

	// sort stage
	if orderBy != 0 {
		pipeline = append(pipeline, bson.D{primitive.E{Key: "account_id", Value: orderBy}})
	}

	// pagination stage
	pipeline = append(pipeline, bson.D{primitive.E{Key: "$skip", Value: offset}})
	pipeline = append(pipeline, bson.D{primitive.E{Key: "$limit", Value: limit}})

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
		// Extract total count from the first document
		metadata := results[0]
		totalCount = metadata["total_count"].(int64)

		// Extract and parse account data documents
		for _, doc := range results[1:] { // Skip the first doc (metadata)
			var account entity.Account
			account.AccountID = doc["account_id"].(int)
			account.Limit = doc["limit"].(int)
			products, ok := doc["products"].(bson.A)
			if ok {
				for _, p := range products {
					account.Products = append(account.Products, p.(string))
				}
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
