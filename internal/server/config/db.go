package config

// dbConfig is a part of the config which contains setting for the DB.
type dbConfig struct {
	MongoURI string `env:"GOPHKEEPER_DB_URI"  envDefault:"mongodb://localhost:27017"`
	DBName   string `env:"GOPHKEEPER_DB_NAME" envDefault:"gophkeeper"`
}
