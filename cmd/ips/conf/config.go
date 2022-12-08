package conf

import (
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

const (
	EnvConfigPath = "IPS_CONFIG_PATH"
	ConfigDirName = "ips"
	ConfigName    = "config"
	ConfigType    = "json"
)

var ConfigPath string
var Conf Config

type Config struct {
	IPv4Format string   `mapstructure:"ipv4_format"`
	IPv4File   string   `mapstructure:"ipv4_file"`
	IPv6Format string   `mapstructure:"ipv6_format"`
	IPv6File   string   `mapstructure:"ipv6_file"`
	Fields     []string `mapstructure:"fields"`
}

func init() {

	// find config path
	ConfigPath = os.Getenv(EnvConfigPath)
	if len(ConfigPath) == 0 {
		ConfigPath = filepath.Join(xdg.ConfigHome, ConfigDirName)
	}
	PrepareDir(ConfigPath)

	// set default config
	viper.SetDefault("ipv4_file", "city.free.ipdb")
	viper.SetDefault("ipv6_file", "city.free.ipdb")
	//viper.SetDefault("fields", []string{"country", "province", "city", "isp"})

	// read config
	viper.SetConfigName(ConfigName)
	viper.SetConfigType(ConfigType)
	viper.AddConfigPath(ConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		if err := viper.SafeWriteConfig(); err != nil {
			log.Fatalf("viper.SafeWriteConfig() failed: %s\n", err)
		}
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatalf("viper.Unmarshal() failed: %s\n", err)
	}
}

func PrepareDir(path string) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0755); err != nil {
				log.Fatalf("mkdir dir(%s) failed: %s\n", path, err)
			}
		} else {
			log.Fatalf("os.Stat() failed: %s\n", err)
		}
	} else if !stat.IsDir() {
		log.Fatalf("path is not dir: %s\n", path)
	}
}
