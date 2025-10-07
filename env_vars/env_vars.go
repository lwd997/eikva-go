package envvars

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type EnvVar string
type EnvVarNumeric string

const (
	NoSSLVerify               EnvVar        = "NO_SSL_VERIFY"
	JwtSectet                 EnvVar        = "JWT_SECRET"
	OpenAiApiKey              EnvVar        = "OPEN_AI_API_KEY"
	OpenAiBaseUrl             EnvVar        = "OPEN_AI_BASE_URL"
	OpenAiCompletionsPathName EnvVar        = "OPEN_AI_COMPLETIONS_PATHNAME"
	LLMTokenTreshold          EnvVarNumeric = "LLM_TOKEN_TRESHOLD"
)

var pool map[EnvVar]string = map[EnvVar]string{}
var poolNumeric map[EnvVarNumeric]int = map[EnvVarNumeric]int{}

func Init() {
	varNames := []EnvVar{
		NoSSLVerify,
		JwtSectet,
		OpenAiApiKey,
		OpenAiBaseUrl,
		OpenAiCompletionsPathName,
	}

	varNamesNumeric := []EnvVarNumeric{
		LLMTokenTreshold,
	}

	errs := []string{}

	for _, entry := range varNames {
		asStr := string(entry)
		pool[entry] = os.Getenv(asStr)
		if pool[entry] == "" {
			errs = append(errs, asStr)
		}
	}

	for _, entry := range varNamesNumeric {
		asStr := string(entry)
		val := os.Getenv(asStr)
		if val == "" {
			errs = append(errs, asStr)
			continue
		}

		numVal, err := strconv.Atoi(val)
		if err != nil {
			errs = append(errs, asStr)
			continue
		}

		poolNumeric[entry] = numVal
	}

	if len(errs) > 0 {
		log.Fatalf("Не указаны или имеют некорректные значения обязательные переменные среды: %s \n", strings.Join(errs, ", "))
	}
}

func Get(name EnvVar) string {
	return pool[name]
}

func GetNumeric(name EnvVarNumeric) int {
	return poolNumeric[name]
}

func Dotenv() {
	if os.Getenv("GIN_MODE") == "release" {
		return
	}

	err := godotenv.Load()

	if err != nil {
		panic(err.Error())
	}
}
