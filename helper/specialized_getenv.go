package helper

import (
	"log"
	"os"
	"strconv"
)

func GetEnvInt(key string) int {
	env, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		log.Println(env)
		log.Panicf("failed to get env variable: %s", err)
	}
	return env
}
