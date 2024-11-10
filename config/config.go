/*
MIT License

# Copyright (c) 2019 David Suarez

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type config struct {
	Solaredge struct {
		InverterAddress string `mapstructure:"inverter_address" yaml:"inverter_address,omitempty"`
		InverterPort    int16  `mapstructure:"inverter_port" yaml:"inverter_port,omitempty"`
		MeterCount      int64  `mapstructure:"meter_count" yaml:"meter_count,omitempty"`
	} `mapstructure:"solaredge" yaml:"solaredge,omitempty"`
	Exporter struct {
		ListenAddress string        `mapstructure:"listen_address" yaml:"listen_address,omitempty"`
		ListenPort    int16         `mapstructure:"listen_port" yaml:"listen_port,omitempty"`
		Interval      time.Duration `mapstructure:"interval" yaml:"interval,omitempty"`
		PanicOnError  bool          `mapstructure:"panic_on_error" yaml:"panic_on_error,omitempty"`
	} `mapstructure:"exporter" yaml:"exporter,omitempty"`
	Log struct {
		Debug bool `mapstructure:"debug" yaml:"debug,omitempty"`
	} `mapstructure:"log" yaml:"log,omitempty"`
}

var (
	Config config
)

// InitConfig initializes the viper configuration
func InitConfig() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SEE")
	viper.AllowEmptyEnv(true)
	viper.EnvKeyReplacer(strings.NewReplacer(
		".", "_",
	))
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/solaredge-exporter")
	viper.AddConfigPath("$HOME/.solaredge-exporter")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("unable to read config: %w", err))
	}
	// workaround because viper does not treat env vars the same as other config
	for _, key := range viper.AllKeys() {
		envKey := viper.GetEnvPrefix() + "_" + strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		err := viper.BindEnv(key, envKey)
		if err != nil {
			panic(fmt.Sprintf("config: unable to bind env: " + err.Error()))
		}
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal config: %w", err))
	}
	if Config.Log.Debug {
		viper.Debug()
	}
}
