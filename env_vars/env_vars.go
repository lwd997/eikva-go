package envvars

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type EnvVar string

const (
	NoSSLVerify               EnvVar = "NO_SSL_VERIFY"
	JwtSectet                 EnvVar = "JWT_SECRET"
	OpenAiApiKey              EnvVar = "OPEN_AI_API_KEY"
	OpenAiBaseUrl             EnvVar = "OPEN_AI_BASE_URL"
	OpenAiCompletionsPathName EnvVar = "OPEN_AI_COMPLETIONS_PATHNAME"
)

var pool map[EnvVar]string = map[EnvVar]string{}

func Init() {
	varNames := []EnvVar{
		NoSSLVerify,
		JwtSectet,
		OpenAiApiKey,
		OpenAiBaseUrl,
		OpenAiCompletionsPathName,
	}

	errs := []string{}

	for _, entry := range varNames {
		asStr := string(entry)
		pool[entry] = os.Getenv(asStr)
		if pool[entry] == "" {
			errs = append(errs, asStr)
		}
	}

	if len(errs) > 0 {
		log.Fatalf("Не указаны обязательные переменные среды: %s \n", strings.Join(errs, ", "))
	}
}

func Get(name EnvVar) string {
	return pool[name]
}

func Dotenv() {
	err := godotenv.Load()

	if err != nil {
		panic(err.Error())
	}
}
