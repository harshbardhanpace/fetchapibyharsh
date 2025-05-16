package loggerconfig

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"space/constants"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

var Info = func(args ...interface{}) {
	InfoImpl(args...)
}

func InfoImpl(args ...interface{}) {
	message := buildMessage(args...)
	logger.Info(message)
}

var Warn = func(args ...interface{}) {
	WarnImpl(args...)
}

func WarnImpl(args ...interface{}) {
	message := buildMessage(args...)
	logger.Warn(message)
}

var Panic = func(args ...interface{}) {
	PanicImpl(args...)
}

func PanicImpl(args ...interface{}) {
	message := buildMessage(args...)
	logger.Panic(message)
}

func buildMessage(args ...interface{}) string {
	var message string
	for i, arg := range args {
		if i == 0 {
			message = fmt.Sprint(arg)
		} else {
			message += " " + fmt.Sprint(arg)
		}
	}
	return message
}

var logger *logrus.Logger

func LogrusInitialize() {

	env := os.Getenv("GO_ENV")

	logger = logrus.New()
	if env == "local" || env == "debug" {
		logger.SetLevel(logrus.DebugLevel)
		logger.Debug("InitLogrus Debug logging enabled.")
	}

	// Set logrus configuration
	logger.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "(service-SPACE)", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		},
	}
	logger.SetFormatter(formatter)

	elasticApmServiceName := os.Getenv(constants.ElasticApmServiceName)
	elasticAPM, err := apm.NewTracer(elasticApmServiceName, constants.ServiceVersion)
	if err != nil {
		log.Printf("Getting Error-%v\n", err)
	}

	var logLevels []logrus.Level

	for _, level := range logrus.AllLevels {
		logLevels = append(logLevels, level)
	}

	logger.AddHook(&apmlogrus.Hook{
		Tracer:    elasticAPM,
		LogLevels: logLevels,
	})

	// Get the current date and time
	currentTime := time.Now().In(constants.LocationKolkata)
	// Format the date and time as a string
	timeStamp := currentTime.Format(constants.DDMMYYYY)

	if constants.FileLoggingEnabled {

		// logFileName := fmt.Sprintf("%s_log_papertrading", timeStamp)
		logFileName := timeStamp + "_" + constants.LogSpace

		// Open log file
		file, err := os.OpenFile(constants.LogFilePath+logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("Failed to create logfile: ", logFileName, "err: ", err)
			panic(err)
		}

		// Create a hook for file logging
		fileHook := NewFileHook(file)
		logger.AddHook(fileHook)
	}
}

// FileHook is a Logrus hook that writes logs to a file
type FileHook struct {
	file *os.File
}

// NewFileHook creates a new instance of FileHook
func NewFileHook(file *os.File) *FileHook {
	return &FileHook{
		file: file,
	}
}

// Fire is called when a log event is fired
func (hook *FileHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	if constants.FileLoggingEnabled {
		_, err = hook.file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Levels returns the log levels that the hook should be triggered for
func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}

type CustomAPMHook struct {
	Tracer *apm.Tracer
}

func (hook *CustomAPMHook) Fire(entry *logrus.Entry) error {
	return nil
}
func (hook *CustomAPMHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type logEntry struct {
	count     int
	timestamp time.Time
}

var (
	logTracker = make(map[string]*logEntry)
	logMutex   sync.Mutex
)

const (
	P2ThresholdToP1 = 10
	P2ThresholdToP0 = 100
	P1ThresholdToP0 = 10
	timeWindow      = time.Minute
)

var Error = func(args ...interface{}) {
	ErrorImpl(args...)
}

func ErrorImpl(args ...interface{}) {
	message := buildMessage(args...)
	severity, mainBody, fullMessage := parseLogMessage(message)
	if severity == "" {
		logrus.Error(fullMessage) // Ignore escalation if severity is missing
		return
	}

	caller := getCallerFunctionName(fullMessage)
	logKey := caller + "::" + mainBody
	escalatedSeverity, reason := checkAndEscalate(logKey, severity)

	// Modify the log message only if the severity changes
	finalLog := fullMessage
	if severity != escalatedSeverity {
		finalLog = strings.Replace(fullMessage, severity, escalatedSeverity, 1)
		if reason != "" {
			finalLog += " (" + reason + ")"
		}
	}

	logrus.Error(finalLog)
}
func checkAndEscalate(logKey, severity string) (string, string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if entry, exists := logTracker[logKey]; exists {
		if time.Since(entry.timestamp) > timeWindow {
			entry.count = 1
			entry.timestamp = time.Now()
		} else {
			entry.count++
		}
	} else {
		logTracker[logKey] = &logEntry{count: 1, timestamp: time.Now()}
	}

	entry := logTracker[logKey]
	var reason string

	if severity == "P2-Mid" {
		if entry.count >= P2ThresholdToP0 {
			reason = fmt.Sprintf("Escalated from P2 to P0 due to %d occurrences in the last minute", entry.count)
			return "P0-Critical", reason
		} else if entry.count >= P2ThresholdToP1 {
			reason = fmt.Sprintf("Escalated from P2 to P1 due to %d occurrences in the last minute", entry.count)
			return "P1-High", reason
		}
	} else if severity == "P1-High" && entry.count >= P1ThresholdToP0 {
		reason = fmt.Sprintf("Escalated from P1 to P0 due to %d occurrences in the last minute", entry.count)
		return "P0-Critical", reason
	}

	return severity, ""
}

func parseLogMessage(logMsg string) (severity, mainBody, fullMessage string) {
	fullMessage = logMsg // Preserve full message for final log output

	// Extract severity
	severityRegex := regexp.MustCompile(`Alert Severity:(P[0-2]-[A-Za-z]+)`)
	severityMatch := severityRegex.FindStringSubmatch(logMsg)
	if len(severityMatch) < 2 {
		return "", "", fullMessage // Ignore escalation if severity is missing
	}
	severity = severityMatch[1]

	// Extract caller function name (PascalCase or CamelCase word after severity)
	callerRegex := regexp.MustCompile(`Alert Severity:P[0-2]-[A-Za-z]+,\s*([A-Z][a-zA-Z0-9]*)`)
	callerMatch := callerRegex.FindStringSubmatch(logMsg)
	if len(callerMatch) < 2 {
		return severity, logMsg, fullMessage // If no caller is found, use full message
	}
	caller := callerMatch[1]

	// Remove known dynamic parts (platform, requestId, ClientID, etc.) for tracking
	ignoreFields := []string{"platform:", "requestid=", "userId=", "ClientID=", "clientVersion="}
	sanitizedMessage := logMsg
	for _, field := range ignoreFields {
		sanitizedMessage = removeField(sanitizedMessage, field)
	}

	// Extract main log body
	mainBodyRegex := regexp.MustCompile(`Alert Severity:P[0-2]-[A-Za-z]+,\s*` + caller + `\s*(.+)`)
	mainBodyMatch := mainBodyRegex.FindStringSubmatch(sanitizedMessage)
	if len(mainBodyMatch) < 2 {
		return severity, sanitizedMessage, fullMessage
	}

	return severity, mainBodyMatch[1], fullMessage
}

func removeField(logMsg, field string) string {
	index := strings.Index(logMsg, field)
	if index != -1 {
		endIndex := strings.Index(logMsg[index:], " ")
		if endIndex != -1 {
			logMsg = logMsg[:index] + logMsg[index+endIndex+1:]
		} else {
			logMsg = logMsg[:index]
		}
	}
	return strings.TrimSpace(logMsg)
}

func getCallerFunctionName(logMsg string) string {
	callerRegex := regexp.MustCompile(`Alert Severity:P[0-2]-[A-Za-z]+,\s*([A-Z][a-zA-Z0-9]*)`)
	callerMatch := callerRegex.FindStringSubmatch(logMsg)
	if len(callerMatch) >= 2 {
		return callerMatch[1]
	}
	return "UnknownCaller"
}
