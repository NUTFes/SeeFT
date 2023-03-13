package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Root     string
		Password string
		Dbname   string
	}
	Server struct {
		Address string
		Port    string
	}
}

var C config

func NewConfig() {

	Conf := &C

	viper.SetConfigName("config")

	viper.SetConfigType("yml")

	viper.AddConfigPath(filepath.Join("$GOPATH", "src", "github.com", "NUTFes", "seeft", "config"))

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("config file read error")
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		fmt.Println("config file Unmarshal error")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(Conf)

	//return conf
}