package db

type DatabaseConfig interface {
	GetMongoDbScheme() string
	GetMongoDbHost() string
	GetMongoDbPort() int
	GetMongoDbDatabase() string
	GetMongoDbUser() string
	GetMongoDbPassword() string
}
