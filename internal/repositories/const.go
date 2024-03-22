package repositories

import "errors"

const (
	ACCOUNTS_COLLECTION_NAME string = "accounts"
)

var (
	errMetaDataTypeAssertion = errors.New("failed to do type assertion on metadata")

	errDataTypeAssertion = errors.New("failed to do type assertion on data")

	errAccountsTypeAssertion = errors.New("failed to do type assertion on account")
)
