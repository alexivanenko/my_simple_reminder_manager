package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-ini/ini"
)

const LOCAL = false
const VERSION = "0.0.1"

var config *ini.File
var stdLogger *log.Logger
var rootDir string

func getValue(section string, key string) *ini.Key {
	return config.Section(section).Key(key)
}

func String(section string, key string) string {
	return getValue(section, key).String()
}

func Is(section string, key string) bool {
	if val, err := getValue(section, key).Bool(); err != nil {
		return false
	} else {
		return val
	}
}

func Log(msg string) {
	stdLogger.Printf("[LOG] %v | %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}

func GetVersion() string {
	return VERSION
}

func GetRootDir() string {

	if rootDir != "" {
		return rootDir
	}

	if LOCAL {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			panic("No caller information")
		}

		rootDir = path.Dir(filename)
		rootDir = strings.Replace(rootDir, "/config", "", 1)
	} else {
		var err error
		if rootDir, err = filepath.Abs(fmt.Sprintf("%s/", filepath.Dir(os.Args[0]))); err != nil {
			panic(err)
		}
	}

	return rootDir
}

func init() {
	stdLogger = log.New(os.Stdout, "", 0)
	Log(fmt.Sprintf("Initializing application ver %s", GetVersion()))

	var err error
	configPath := fmt.Sprintf("%s/config.ini", GetRootDir())

	Log(fmt.Sprintf("Loading base config from %s", configPath))

	if config, err = ini.Load(configPath); err != nil {
		panic(err)
	}
}
