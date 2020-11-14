package main

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	v := os.Getenv("GORM_DIALECT")
	t.Log(v)
}
