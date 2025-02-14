package model

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	Config configInterface = &ConfigStruct{}
	theme  string
)

type configInterface interface {
	Get() *ConfigStruct
	Set(*ConfigStruct)
}

type ConfigStruct struct {
	Logrus LogrusConfig `yaml:"logrus"`
}

type LogrusConfig struct {
	Output   []string `yaml:"output"`
	Reporter bool     `yaml:"reporter"`
	Format   string   `yaml:"format"`
	Path     string   `yaml:"path"`
	Level    int      `yaml:"level"`
	LogFile  *os.File
}

func (c *ConfigStruct) Get() *ConfigStruct {
	log.Trace()

	return c
}

func (c *ConfigStruct) Set(cfg *ConfigStruct) {
	log.Trace()

	c.Logrus.Output = cfg.Logrus.Output
	c.Logrus.Reporter = cfg.Logrus.Reporter
	c.Logrus.Format = cfg.Logrus.Format
	c.Logrus.Path = cfg.Logrus.Path
	c.Logrus.Level = cfg.Logrus.Level
	c.Logrus.LogFile = cfg.Logrus.LogFile
}
