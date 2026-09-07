package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func chdir(t *testing.T, dir string) {
	t.Helper()

	pwd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(pwd) })
}

func TestLoadReadsEnvFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(dir, ".env"),
		[]byte("token=tok\napp=appid\ntestserver=guild\n"),
		0o600,
	))
	chdir(t, dir)

	for _, key := range []string{"token", "app", "testserver"} {
		t.Setenv(key, "")
		require.NoError(t, os.Unsetenv(key))
	}

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, "tok", cfg.Token)
	assert.Equal(t, "appid", cfg.App)
	assert.Equal(t, "guild", cfg.TestServerID)
	assert.Equal(t, config.BOLT, cfg.Store)
}

func TestLoadErrorsWithoutEnvFile(t *testing.T) {
	chdir(t, t.TempDir())

	_, err := config.Load()

	require.Error(t, err)
}
