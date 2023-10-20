/*
 * Copyright (c) 2023 shenjunzheng@gmail.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"os"

	"github.com/spf13/viper"

	"github.com/sjzar/ips/internal/util"
	"github.com/sjzar/ips/pkg/errors"
)

const (
	DefaultConfigType = "json"
)

var (
	// ConfigName holds the name of the configuration file.
	ConfigName = ""

	// ConfigType specifies the type/format of the configuration file.
	ConfigType = ""

	// ConfigPath denotes the path to the configuration file.
	ConfigPath = ""
)

// Init initializes the configuration settings.
// It sets up the name, type, and path for the configuration file.
func Init(name, _type, path string) error {
	if len(name) == 0 {
		return errors.ErrMissingConfigName
	}

	if len(_type) == 0 {
		_type = DefaultConfigType
	}

	var err error
	if len(path) == 0 {
		path, err = os.UserHomeDir()
		if err != nil {
			path = os.TempDir()
		}
		path += string(os.PathSeparator) + "." + name
	}
	if err := util.PrepareDir(path); err != nil {
		return err
	}

	ConfigName = name
	ConfigType = _type
	ConfigPath = path
	return nil
}

// Load loads the configuration from the previously initialized file.
// It unmarshals the configuration into the provided conf interface.
func Load(conf interface{}) error {
	viper.SetConfigName(ConfigName)
	viper.SetConfigType(ConfigType)
	viper.AddConfigPath(ConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		if err := viper.SafeWriteConfig(); err != nil {
			return err
		}
	}
	if err := viper.Unmarshal(conf); err != nil {
		return err
	}
	SetDefault(conf)
	return nil
}

// LoadFile loads the configuration from a specified file.
// It unmarshals the configuration into the provided conf interface.
func LoadFile(file string, conf interface{}) error {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(conf); err != nil {
		return err
	}
	SetDefault(conf)
	return nil
}

// SetConfig sets a configuration key to a specified value.
// It also writes the updated configuration back to the file.
func SetConfig(key string, value interface{}) error {
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

// ResetConfig resets the configuration to empty.
func ResetConfig() error {
	viper.Reset()
	viper.SetConfigName(ConfigName)
	viper.SetConfigType(ConfigType)
	viper.AddConfigPath(ConfigPath)
	return viper.WriteConfig()
}

// GetConfig retrieves all configuration settings as a map.
func GetConfig() map[string]interface{} {
	return viper.AllSettings()
}
