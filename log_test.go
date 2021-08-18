package nanolog

import (
	"bytes"
	"testing"
)

func TestDefaults(t *testing.T) {
	Error().Println("test success")
}

func TestLogMessage(t *testing.T) {

	mockDebugWriter := bytes.NewBuffer([]byte{})

	Init(Options{
		Level: DebugLevel,
		Debug: LoggerOptions{
			Writer: mockDebugWriter,
			Flags:  -1,
		},
	})

	text := "test debug"
	message := FormatPrefix(DefaultPrefix, DebugColor, DebugLevel) + text + "\n"
	Debug().Println(text)
	currentMessage := mockDebugWriter.String()
	if currentMessage != message {
		t.Errorf("message %q != %q", currentMessage, message)
	}

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

func TestDefaultLogger(t *testing.T) {

	mockDebugWriter := bytes.NewBuffer([]byte{})

	Init(Options{
		Level: DebugLevel,
		Debug: LoggerOptions{
			Writer: mockDebugWriter,
			Flags:  -1,
		},
	})

	text := "test debug"
	message := FormatPrefix(DefaultPrefix, DebugColor, DebugLevel) + text + "\n"
	DefaultLogger().Debug().Println(text)
	currentMessage := mockDebugWriter.String()
	if currentMessage != message {
		t.Errorf("message %q != %q", currentMessage, message)
	}

}
