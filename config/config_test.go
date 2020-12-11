package config

import (
	"5z7Game/pkg/utils"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMain(m *testing.M) {
	_, fileStr, _, _ := runtime.Caller(0)
	utils.SecurePanic(ReadFromFile(filepath.Dir(fileStr) + "/../app.yaml"))
	m.Run()
}

func TestServer(t *testing.T)  {
	assert.Equal(t, Server().Port, 8091)
}
