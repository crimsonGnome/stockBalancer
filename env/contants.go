package env

import (
	"fmt"
	"os"
)

func getEnvVariable(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		err := fmt.Errorf("environmentalVariable: %s does not exist", key)
		panic(err)
	}
	return value
}

var ENV_END_DATE = getEnvVariable("ENV_END_DATE")
