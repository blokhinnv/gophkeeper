package config

// dbConfig is a part of the config which contains setting for the DB.
type dbConfig struct {
	MongoURI string `env:"MONGO_URI" envDefault:"mongodb://localhost:27017"`
}
