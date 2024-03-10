package configs

import (
	"github.com/spf13/viper"
)

type conf struct {
	MaxReqPermitidas string `mapstructure:"MAX_REQ_PERMITIDAS"`
	TempoBloqueioSec string `mapstructure:"TEMPO_BLOQUEIO_SEC"`
	MetodoBloqueio   string `mapstructure:"METODO_BLOQUEIO"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, err
}
