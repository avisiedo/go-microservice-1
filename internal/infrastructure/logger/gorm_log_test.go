package logger

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/logger"
)

type TestGiven struct {
	Message      string
	LevelLog     slog.Level
	LevelMessage slog.Level
}
type (
	TestExpected int
	TestCase     struct {
		Name     string
		Given    TestGiven
		Expected TestExpected
	}
)

const (
	RegexpMatch TestExpected = 1
	EmptyString TestExpected = 2
)

// https://regex101.com/
const (
	timeExp    = `time=[[:digit:]]{4}\-[[:digit:]]{2}\-[[:digit:]]{2}T[[:digit:]]{2}:[[:digit:]]{2}:[[:digit:]]{2}.[[:digit:]]+(|Z|\+[[:digit:]]{2}:[[:digit:]]{2})`
	levelExp   = `level=(DEBUG-4|DEBUG|WARNING|NOTICE|INFO|WARN|ERROR|FATAL)`
	regularExp = "^" + timeExp + " " + levelExp + " " + `msg=(.*)\n$`
)

func helperCommonTestCase(levelMsg slog.Level) []TestCase {
	return []TestCase{
		{
			Name: "Display test message for DEBUG level and DEBUG messsage",
			Given: TestGiven{
				LevelLog:     levelMsg,
				LevelMessage: levelMsg,
				Message:      "test",
			},
			Expected: RegexpMatch,
		},
		{
			Name: "NO Display test message for INFO level and DEBUG message",
			Given: TestGiven{
				LevelLog:     LevelSilent,
				LevelMessage: levelMsg,
				Message:      "test",
			},
			Expected: EmptyString,
		},
	}
}

func TestNewGormLog(t *testing.T) {
	r := NewGormLog(nil, true)
	require.NotNil(t, r)
}

func TestLogMode(t *testing.T) {
	l := NewGormLog(nil, true)
	values := []logger.LogLevel{logger.Error, logger.Info, logger.Silent, logger.Warn}
	for _, value := range values {
		assert.NotPanics(t, func() {
			l.LogMode(value)
		})
	}
}

func testCommon(t *testing.T, testCases []TestCase, callback func(l *gormLogger, level slog.Level, msg string)) {
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		b := bytes.Buffer{}
		sh := slog.NewTextHandler(&b, &slog.HandlerOptions{
			Level: testCase.Given.LevelLog,
		})
		sl := slog.New(sh)
		l := &gormLogger{
			slogger:                   sl,
			IgnoreRecordNotFoundError: true,
		}
		callback(l, testCase.Given.LevelMessage, testCase.Given.Message)
		switch testCase.Expected {
		case RegexpMatch:
			assert.Regexp(t, regularExp, b.String())
		case EmptyString:
			assert.Empty(t, b.String())
		}
	}
}

func testCommonSpecific(t *testing.T, testCases []TestCase, callback func(l *gormLogger, msg string)) {
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		b := bytes.Buffer{}
		sh := slog.NewTextHandler(&b, &slog.HandlerOptions{
			Level: testCase.Given.LevelLog,
		})
		sl := slog.New(sh)
		l := &gormLogger{
			slogger:                   sl,
			IgnoreRecordNotFoundError: true,
		}
		callback(l, testCase.Given.Message)
		switch testCase.Expected {
		case RegexpMatch:
			assert.Regexp(t, regularExp, b.String())
		case EmptyString:
			assert.Empty(t, b.String())
		}
	}
}

func TestLogCommon(t *testing.T) {
	testCases := helperCommonTestCase(LevelDebug)
	testCommon(t, testCases, func(l *gormLogger, level slog.Level, msg string) {
		l.logCommon(context.Background(), level, msg)
	})
}

func TestLogMsg(t *testing.T) {
	testCases := helperCommonTestCase(LevelDebug)
	testCommon(t, testCases, func(l *gormLogger, level slog.Level, msg string) {
		l.logMsg(context.Background(), level, "%s", msg)
	})
}

func TestLog(t *testing.T) {
	testCases := helperCommonTestCase(LevelDebug)
	testCommon(t, testCases, func(l *gormLogger, level slog.Level, msg string) {
		l.log(context.Background(), level, "%s", msg)
	})
}

func TestInfo(t *testing.T) {
	testCases := helperCommonTestCase(LevelInfo)
	testCommonSpecific(t, testCases, func(l *gormLogger, msg string) {
		l.Info(context.Background(), "%s", msg)
	})
}

func TestWarn(t *testing.T) {
	testCases := helperCommonTestCase(LevelWarn)
	testCommonSpecific(t, testCases, func(l *gormLogger, msg string) {
		l.Warn(context.Background(), "%s", msg)
	})
}

func TestError(t *testing.T) {
	testCases := helperCommonTestCase(LevelError)
	testCommonSpecific(t, testCases, func(l *gormLogger, msg string) {
		l.Error(context.Background(), "%s", msg)
	})
}

func TestTraceNoError(t *testing.T) {
	// https://regex101.com/
	const regularExpTraceNoError = "^" + timeExp + " " + levelExp + " " + `msg=(.*)\n$`
	b := bytes.Buffer{}
	sh := slog.NewTextHandler(&b, &slog.HandlerOptions{
		Level: LevelTrace,
	})
	sl := slog.New(sh)
	l := &gormLogger{
		slogger:                   sl,
		IgnoreRecordNotFoundError: true,
	}
	l.Trace(context.Background(), time.Now(), func() (sql string, rowsAffected int64) {
		return "SELECT COUNT(*) FROM USERS;", 0
	}, nil)
	assert.Regexp(t, regularExpTraceNoError, b.String())
}

func TestTraceError(t *testing.T) {
	// https://regex101.com/
	const regularExpTraceError = "^" + timeExp + " " + levelExp + " " + `msg=(.*) error=(.*) query=\"(.*)\" elapsed=(.*) rows=(.*)\n$`
	b := bytes.Buffer{}
	sh := slog.NewTextHandler(&b, &slog.HandlerOptions{
		Level: LevelTrace,
	})
	sl := slog.New(sh)
	l := &gormLogger{
		slogger:                   sl,
		IgnoreRecordNotFoundError: true,
	}
	l.Trace(context.Background(), time.Now(), func() (sql string, rowsAffected int64) {
		return "SELECT COUNT(*) FROM USERS;", 0
	}, errors.New("test-error"))
	assert.Regexp(t, regularExpTraceError, b.String())
}
