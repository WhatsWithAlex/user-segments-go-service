package env

import (
	"os"
	"strings"

	"github.com/WhatsWithAlex/user-segments-go-service/pkg/koanfext"
	"github.com/knadh/koanf/parsers/dotenv"
	kEnv "github.com/knadh/koanf/providers/env"
	kFile "github.com/knadh/koanf/providers/file"
	"github.com/rs/zerolog/log"
)

type (
	DBEnv struct {
		Name     string `koanf:"name"`
		User     string `koanf:"user"`
		Password string `koanf:"password"`
		Addr     string `koanf:"addr"`
		Port     string `koanf:"port"`
	}

	AppEnv struct {
		Addr    string `koanf:"addr"`
		Port    string `koanf:"port"`
		Timeout int    `koanf:"timeout"`
	}

	Env struct {
		Environment string `koanf:"environment"`
		DB          DBEnv  `koanf:"db"`
		App         AppEnv `koanf:"app"`
	}
)

const (
	delim  = "."
	prefix = "ENV_"
)

var k = koanfext.New(delim)

func LoadEnv() (env Env, err error) {
	tranformation_callback := func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, prefix)), "_", delim, -1)
	}

	err = k.Load(kFile.Provider(".env"), dotenv.ParserEnv("", delim, tranformation_callback))
	if err != nil {
		if os.IsNotExist(err) {
			log.Info().Msg("env: .env file not found")
		} else {
			return
		}
	}

	log.Info().Msg("env: loading environment variables")
	err = k.Load(kEnv.Provider(prefix, delim, tranformation_callback), nil)
	if err != nil {
		return
	}

	err = k.StrictUnmarshal("", &env)
	if err != nil {
		return
	}

	return
}
