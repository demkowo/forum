package logger

import (
	"io"
	"os"

	model "github.com/demkowo/forum/models"
	log "github.com/sirupsen/logrus"
)

var (
	Start loggerInterface = &loggerStruct{}
	m     *model.ConfigStruct
)

type loggerInterface interface {
	YamlConfig()
	BasicConfig()
}

type loggerStruct struct {
}

func (c *loggerStruct) BasicConfig() {
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&CustomFormatter{})
	log.SetReportCaller(true)

	f, err := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("os.OpenFile failed\n[%s]\n", err)
	}
	multi := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multi)

}

func (c *loggerStruct) YamlConfig() {
	log.Trace()

	m = model.Config.Get()

	setLevel()
	setFormat()
	setReporter()
	setStdout()
}

func setLevel() {
	switch m.Logrus.Level {
	case 0:
		log.SetLevel(log.PanicLevel)
		log.Info("PanicLevel")
	case 1:
		log.SetLevel(log.FatalLevel)
		log.Info("FatalLevel")
	case 2:
		log.SetLevel(log.ErrorLevel)
		log.Info("ErrorLevel")
	case 3:
		log.SetLevel(log.WarnLevel)
		log.Info("WarnLevel")
	case 4:
		log.SetLevel(log.InfoLevel)
		log.Info("InfoLevel")
	case 5:
		log.SetLevel(log.DebugLevel)
		log.Info("DebugLevel")
	case 6:
		log.SetLevel(log.TraceLevel)
		log.Info("TraceLevel")
	default:
		log.SetLevel(log.WarnLevel)
		log.Info("WarnLevel")
	}
}

func setFormat() {

	switch m.Logrus.Format {
	case "text":
		log.SetFormatter(&log.TextFormatter{})
		log.Info("TextFormatter")
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
		log.Info("JSONFormatter")
	case "custom":
		log.SetFormatter(&CustomFormatter{})
		log.Info("CustomFormatter")
	default:
		log.SetFormatter(&log.TextFormatter{})
		log.Info("TextFormatter")
	}
}

func setReporter() {

	if m.Logrus.Reporter {
		log.SetReportCaller(m.Logrus.Reporter)
		log.Info(m.Logrus.Reporter)
	}

}

func setStdout() {
	var writers []io.Writer
	var err error

	for _, out := range m.Logrus.Output {
		if out == "stdout" {
			writer := io.Writer(os.Stdout)
			writers = append(writers, writer)
			log.Info("add writer:   os.Stdout")
		}
		if out == "file" {
			m.Logrus.LogFile, err = os.OpenFile(m.Logrus.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Panicf("os.OpenFile failed\n[%s]\n", err)
			}
			writer := io.Writer(m.Logrus.LogFile)
			writers = append(writers, writer)
			log.Info("add writer:   file")
		}
	}

	if len(writers) < 1 {
		log.SetOutput(os.Stdout)
	}

	multi := io.MultiWriter(writers...)
	log.SetOutput(multi)

}
