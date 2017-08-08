package clit

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/eaciit/appconfig"
	"github.com/eaciit/toolkit"
)

var (
	EnableConfig = true
	EnableLog    = true
	Log          *toolkit.LogEngine

	preFn   func() error
	closeFn func() error
	flags   = map[string]*string{}

	configs map[string]*appconfig.Config
	exeDir  string
)

func init() {
	SetFlag("config", "", "Location of the config file")
}

func SetFlag(name, value, usage string) *string {
	f := flag.String(name, value, usage)
	flags[name] = f
	return f
}

func SetPreFn(fn func() error) {
	preFn = fn
}

func SetCloseFn(fn func() error) {
	closeFn = fn
}

func Flag(name string) string {
	f := flags[name]
	return *f
}

func Start() error {
	var err error
	flag.Parse()

	if EnableConfig {
		if err = ReadConfig("", ""); err != nil {
			return fmt.Errorf("error reading config file. %s", err.Error())
		}
	}

	if EnableLog && Log == nil {
		if Log, err = toolkit.NewLog(true, false, "", "", ""); err != nil {
			return fmt.Errorf("error preparing log. %s", err.Error())
		}
	}

	if preFn == nil {
		preFn()
	}

	return nil
}

func ExeDir() string {
	if exeDir == "" {
		exeDir, _ = os.Executable()
		exeDir = filepath.Dir(exeDir)
	}
	return exeDir
}

func ReadConfig(name, path string) error {
	if name == "" {
		name = "default"
	}

	if name == "default" && path == "" {
		path := Flag("config")
		if path == "" {
			path = filepath.Join(ExeDir(), "app.json")
		}
	} else if path == "" {
		return errors.New("path can't be empty")
	}

	initConfigs()
	config := new(appconfig.Config)
	if err := config.SetConfigFile(path); err != nil {
		return err
	}
	return nil
}

func Config(name, key string, def interface{}) interface{} {
	initConfigs()
	if name == "" {
		name = "default"
	}
	config, found := configs[name]
	if !found {
		return def
	}
	return config.GetDefault(name, def)
}

func SetConfig(name, key string, value interface{}) {
	initConfigs()
	if name == "" {
		name = "default"
	}
	config, found := configs[name]
	if !found {
		return
	}
	config.Set(key, value)
}

func WriteConfig(name string) error {
	initConfigs()
	config, found := configs[name]
	if !found {
		return fmt.Errorf("can not write config. config %s is not yet initialized", name)
	}
	return config.Write()
}

func initConfigs() {
	configs = map[string]*appconfig.Config{}
}

func Close() {
	if closeFn != nil {
		closeFn()
	}
}
