package log

import (
	"fmt"
	"github.com/NidzamuddinMuzakki/go-libs/env"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()

	serviceName    = env.String("MainSetup.ServiceName", "")
	serviceType    = env.String("MainSetup.ServiceType", "")
	serviceCode    = env.String("MainSetup.ServiceCode", "")
	versionRelease = env.String("Version.ReleaseVersion", "")
	versionType    = env.String("Version.VersionType", "")

	logLevel = logrus.ErrorLevel
)

func runPlan() {
	var err error

	year, month, day := time.Now().Local().Date()
	formatDate := fmt.Sprintf("%s_%d_%s_%d.log", env.String("Logging.logFile.FileName", "logs/log"), day, month, year)

	switch env.String("Logging.logFile.Loglevel", "Error") {
	case "Info":
		logLevel = logrus.InfoLevel
	case "Warn":
		logLevel = logrus.WarnLevel
	case "Error":
		logLevel = logrus.ErrorLevel
	case "Fatal":
		logLevel = logrus.FatalLevel
	case "Panic":
		logLevel = logrus.PanicLevel
	}

	log.SetLevel(logLevel)
	log.SetOutput(os.Stderr)

	logPrettyPrintFile, _ := strconv.ParseBool(env.String("Logging.logFile.PrettyPrint", "false"))
	file, err := os.OpenFile(formatDate, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		logrus.Infof("Failed to log to file, using default stderr: %v", err)
	} else {
		log.Out = file
	}

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: env.String("Logging.logFile.TimeFormat", "2006-01-02 15:04:05"),
		PrettyPrint:     logPrettyPrintFile,
	})
}

type (
	Logging struct {
		SkipCaller int
	}
	ILogging interface {
		Trace(traceID, message string, data interface{})
		Debug(traceID, message string, data interface{})
		Info(traceID, message string, data interface{})
		Warning(traceID, message string, data interface{})
		Error(traceID, message string, data interface{})
		Fatal(traceID, message string, data interface{})
		Panic(traceID, message string, data interface{})
		Http(traceID, message, url, method string, header, req, res interface{})
	}
)

func NewLogging(SkipCaller int) ILogging {
	runPlan()
	return &Logging{SkipCaller: SkipCaller}
}

func fieldLog(trace_id string, SkipCaller int, data interface{}) logrus.Fields {
	pc, file, line, _ := runtime.Caller(SkipCaller)
	return logrus.Fields{
		"service_name":    serviceName,
		"service_type":    serviceType,
		"service_code":    serviceCode,
		"version_release": versionRelease,
		"version_type":    versionType,
		"data":            data,
		"trace_id":        trace_id,
		"package":         runtime.FuncForPC(pc).Name(),
		"file":            file,
		"line":            line,
	}
}

func (lib *Logging) Trace(traceID, message string, data interface{}) {
	items := fieldLog(traceID, lib.SkipCaller, data)
	log.WithFields(items).Debug(message)
}

func (lib *Logging) Debug(traceID, message string, data interface{}) {
	items := fieldLog(traceID, lib.SkipCaller, data)
	log.WithFields(items).Debug(message)
}

func (lib *Logging) Info(traceID, message string, data interface{}) {
	items := fieldLog(traceID, lib.SkipCaller, data)
	log.WithFields(items).Info(message)
	logrus.WithFields(items).Info(message)

}

func (lib *Logging) Warning(traceID, message string, data interface{}) {
	items := fieldLog(traceID, lib.SkipCaller, data)
	log.WithFields(items).Warn(message)
	logrus.WithFields(items).Info(message)

}

func (lib *Logging) Error(traceID, message string, data interface{}) {
	items := fieldLog(traceID, lib.SkipCaller, data)
	log.WithFields(items).Error(message)
	logrus.WithFields(items).Info(message)
}

func (lib *Logging) Fatal(traceID, message string, data interface{}) {
	items := fieldLog(traceID, lib.SkipCaller, data)
	log.WithFields(items).Fatal(message)
	logrus.WithFields(items).Info(message)
}

func (lib *Logging) Panic(traceID, message string, data interface{}) {
	items := fieldLog(traceID, lib.SkipCaller, data)
	log.WithFields(items).Panic(message)
	logrus.WithFields(items).Info(message)
}

func (lib *Logging) Http(traceID, message, url, method string, header, req, res interface{}) {
	pc, file, line, _ := runtime.Caller(lib.SkipCaller - 1)
	items := logrus.Fields{
		"host_url":        url,
		"host_method":     method,
		"host_header":     header,
		"host_request":    req,
		"host_response":   res,
		"service_name":    serviceName,
		"service_type":    serviceType,
		"service_code":    serviceCode,
		"version_release": versionRelease,
		"version_type":    versionType,
		"trace_id":        traceID,
		"package":         runtime.FuncForPC(pc).Name(),
		"file":            file,
		"line":            line,
	}
	log.WithFields(items).Info(message)
}
