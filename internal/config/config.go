package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

type PostgreSQL struct {
	Username string `yaml:"psql_user" json:"psql_user" env-default:"user"`
	Password string `yaml:"psql_passqord" json:"psql_passqord" env-default:"password"`
	Host     string `yaml:"psql_host" json:"psql_host" env-default:"localhost"`
	Port     string `yaml:"psql_port" json:"psql_port" env-default:"5432"`
	Database string `yaml:"psql_database" json:"psql_database" env-default:"sampledb"`
}

type Config struct {
	Paths []string        `yaml:"paths" json:"paths" env-default:"./"`
	Ext   map[string]bool `yaml:"ext" json:"ext" env-default:"pdf,djvu"`
	PgSQL PostgreSQL      `yaml:"pgsql" json:"pgsql"`
}

var instance *Config
var once sync.Once

func GetConfig(fname string) *Config {
	once.Do(func() {

		instance = &Config{}
		// if err := cleanenv.ReadEnv(instance); err != nil {
		// 	helpText := "Notes system"
		// 	help, _ := cleanenv.GetDescription(instance, &helpText)
		// 	log.Print(help)
		// 	log.Fatal(err)
		// }

		if err := cleanenv.ReadConfig(fname, instance); err != nil {
			helpText := "Notes system"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
		}
	})

	return instance
}

func SaveSampleConfig(fname string) {
	tc := Config{
		Paths: []string{"./", "d:/"},
		Ext:   map[string]bool{".pdf": true, ".djvu": true},
		PgSQL: PostgreSQL{
			Username: "user",
			Password: "password",
			Host:     "localhost",
			Port:     "5432",
			Database: "sampledb"},
	}

	res, err := yaml.Marshal(tc)
	if err != nil {
		log.Println(err)
	}

	file, err := os.Create(fname)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.Write(res)

}
