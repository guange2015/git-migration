package utils

import "testing"

func TestReadLine(t *testing.T) {
	ReadLine("./file.go", func(line string) {
		t.Log(line)
	})
}
