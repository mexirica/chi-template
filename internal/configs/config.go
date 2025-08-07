// Handles loading and validation of environment variables
package configs

import (
	"fmt"
	"os" // Importe a biblioteca os

	"github.com/spf13/viper"
)

// Config struct with environment variables
type Config struct {
	NAME        string `mapstructure:"NAME"`
	PORT        string `mapstructure:"PORT"`
	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PORT     string `mapstructure:"DB_PORT"`
	DB_USER     string `mapstructure:"DB_USER"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_NAME     string `mapstructure:"DB_NAME"`
	ENVIRONMENT string `mapstructure:"ENVIRONMENT"`
	REDIS_HOST  string `mapstructure:"REDIS_HOST"`
	REDIS_PORT  string `mapstructure:"REDIS_PORT"`
}

// LoadConfig loads environment variables using viper and .env file
func LoadConfig() (config Config, err error) {
	// Limpe o viper para evitar configurações de tentativas anteriores
	viper.Reset()

	// Diga explicitamente ao Viper qual variável de ambiente corresponde a qual chave.
	// O primeiro argumento é a "chave" que você usará (por exemplo, "DB_HOST").
	// O segundo argumento é o nome da variável de ambiente (por exemplo, "DB_HOST").
	// Esta abordagem é infalível.
	viper.BindEnv("NAME", "NAME")
	viper.BindEnv("PORT", "PORT")
	viper.BindEnv("DB_HOST", "DB_HOST")
	viper.BindEnv("DB_PORT", "DB_PORT")
	viper.BindEnv("DB_USER", "DB_USER")
	viper.BindEnv("DB_PASSWORD", "DB_PASSWORD")
	viper.BindEnv("DB_NAME", "DB_NAME")
	viper.BindEnv("ENVIRONMENT", "ENVIRONMENT")
	viper.BindEnv("REDIS_HOST", "REDIS_HOST")
	viper.BindEnv("REDIS_PORT", "REDIS_PORT")
	
	// O Unmarshal agora usará as ligações que você acabou de criar.
	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	// Verifica se as variáveis de ambiente foram lidas manualmente, para debug
	fmt.Printf("Debug: REDIS_HOST do sistema: %s\n", os.Getenv("REDIS_HOST"))

	fmt.Printf("Loaded config: %+v\n", config)
	fmt.Printf("redis %s:%s\n", config.REDIS_HOST, config.REDIS_PORT)
	return
}