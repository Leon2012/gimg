package gimg

import (
	_ "fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	logger, err := NewLogger("test", 0)
	if err != nil {
		t.Fatal()
	} else {
		//message := fmt.Sprintf("test %s, %d", "test", 200)
		//t.Log(message)

		logger.Info("test %s, %d", "test", 200)
	}
}
