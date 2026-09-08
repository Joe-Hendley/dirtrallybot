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
		[]byte("TOKEN=tok\nAPP=appid\nTESTSERVER=guild\n"),
		0o600,
	))
	chdir(t, dir)

	for _, key := range []string{"TOKEN", "APP", "TESTSERVER"} {
		t.Setenv(key, "")
		require.NoError(t, os.Unsetenv(key))
	}

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, "tok", cfg.Token)
	assert.Equal(t, "appid", cfg.App)
	assert.Equal(t, "guild", cfg.TestServerID)
	assert.Equal(t, config.BOLT, cfg.Store)
	assert.Equal(t, config.RandomiserBiased, cfg.Randomiser)
}

func writeEnvFile(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte("TOKEN=tok\n"), 0o600))
	chdir(t, dir)
}

func TestLoadReadsRandomiserFromEnv(t *testing.T) {
	writeEnvFile(t)
	t.Setenv("RANDOMISER", string(config.RandomiserRandom))

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, config.RandomiserRandom, cfg.Randomiser)
}

func TestLoadRejectsUnknownRandomiser(t *testing.T) {
	writeEnvFile(t)
	t.Setenv("RANDOMISER", "nonsense")

	_, err := config.Load()

	require.Error(t, err)
}

func TestLoadErrorsWithoutEnvFile(t *testing.T) {
	chdir(t, t.TempDir())

	_, err := config.Load()

	require.Error(t, err)
}
