package config

import (
	"flag"
)

var (
	MongoURI = flag.String("mongo", "mongodb://mongo", "MongoDB URI.")
	MongoDBName = flag.String("db", "media_stats", "MongoDB internal database name")
)

func init() {
	flag.Parse()
}
