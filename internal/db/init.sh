#!/bin/bash

mongosh --host localhost --port ${DB_PORT} -u ${MONGO_INITDB_ROOT_USERNAME} -p ${MONGO_INITDB_ROOT_PASSWORD} << EOF
use ${DB_NAME}
db.createCollection("accounts")

# Import data
mongoimport --db sample_analytics --collection accounts --jsonArray --file accounts.json

# Create index
db.accounts.createIndex({ 'acoount_id': 1 })
db.accounts.createIndex({'products': 1})

EOF



