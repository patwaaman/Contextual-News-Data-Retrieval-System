package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	MongoURI    string
	OpenAIKey   string
	DebugMongo  bool
	DebugOpenAI bool
}

func Load() Config {
	return Config{
		MongoURI:  mustEnv("MONGO_URI"),
		OpenAIKey: mustEnv("OPENAI_KEY"),

		DebugMongo:  optionalBoolEnv("DEBUG_MONGO", true),
		DebugOpenAI: optionalBoolEnv("DEBUG_OPENAI", true),
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return v
}

func optionalBoolEnv(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Printf("invalid boolean value for %s: %q, using default %v", key, v, def)
		return def
	}

	return b
}
