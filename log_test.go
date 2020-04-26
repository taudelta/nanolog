package nanolog

import (
	"bytes"
	"testing"
)

func TestDefaults(t *testing.T) {
	Error().Println("test success")
}

func TestLogPriority(t *testing.T) {

	mockDebugWriter := bytes.NewBuffer([]byte{})
	mockInfoWriter := bytes.NewBuffer([]byte{})

	Init(Options{
		Level: InfoLevel,
		Debug: LoggerOptions{Writer: mockDebugWriter},
		Info:  LoggerOptions{Writer: mockInfoWriter},
	})

	Info().Println("info")
	if mockInfoWriter.String() == "" {
		t.Errorf("info logger don't wrote logged message, but should")
	}
	Debug().Println("debug")
	if mockDebugWriter.String() != "" {
		t.Errorf("debug logger wrote message, but should not")
	}
}
