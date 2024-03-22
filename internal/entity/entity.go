package entity

type Account struct {
	AccountID int      `bson:"account_id"`
	Limit     int      `bson:"limit"`
	Products  []string `bson:"products"`
}
