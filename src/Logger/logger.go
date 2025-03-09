package Logger

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

type ENUM_LOG_LEVEL int32

const (
	LOG_LEVEL_INFO ENUM_LOG_LEVEL = iota
	LOG_LEVEL_WARNING
	LOG_LEVEL_ERROR
	LOG_LEVEL_DEBUG
)

func (this ENUM_LOG_LEVEL) String() string {
	return []string{
		"INFO",
		"WARNING",
		"ERROR",
		"DEBUG",
	}[this]
}

func generateLogString(level ENUM_LOG_LEVEL, args ...any) string {
	var stream *bytes.Buffer = &bytes.Buffer{}

	stream.WriteString(time.Now().Format("[2006-01-02 15:04:05] "))
	stream.WriteString(level.String() + ": ")

	fmt.Fprint(stream, args...)

	data, err := io.ReadAll(stream)

	if err != nil {
		return ""
	}

	return string(data)
}

func Log(level ENUM_LOG_LEVEL, args ...any) {
	var output string = generateLogString(level, args...)

	fmt.Println(output)
}

func F_INFO(args ...any) string {
	return generateLogString(LOG_LEVEL_INFO, args...)
}

func F_WARN(args ...any) string {
	return generateLogString(LOG_LEVEL_WARNING, args...)
}

func F_ERROR(args ...any) string {
	return generateLogString(LOG_LEVEL_ERROR, args...)
}

func F_DEBUG(args ...any) string {
	return generateLogString(LOG_LEVEL_DEBUG, args...)
}

func INFO(args ...any) {
	Log(LOG_LEVEL_INFO, args...)
}

func WARN(args ...any) {
	Log(LOG_LEVEL_WARNING, args...)
}

func ERROR(args ...any) {
	Log(LOG_LEVEL_ERROR, args...)
}

func DEBUG(args ...any) {
	Log(LOG_LEVEL_DEBUG, args...)
}
