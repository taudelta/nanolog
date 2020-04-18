package nanolog

import (
	"bytes"
	"testing"
)

func TestLogPriority(t *testing.T) {

	mockDebugWriter := bytes.NewBuffer([]byte{})
	mockInfoWriter := bytes.NewBuffer([]byte{})

	Init(Options{
		Level: InfoLevel,
		Debug: mockDebugWriter,
		Info:  mockInfoWriter,
	})

	INFO.Println("info")
	if mockInfoWriter.String() == "" {
		t.Errorf("info logger don't wrote logged message, but should")
	}
	DEBUG.Println("debug")
	if mockDebugWriter.String() != "" {
		t.Errorf("debug logger wrote message, but should not")
	}
}
