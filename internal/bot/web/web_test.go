package web_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/web"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/memorystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const challengeID = "1200000000000000000"

func TestChallengesPageListsStoredChallenges(t *testing.T) {
	ctx := context.Background()
	store := memorystore.New()

	// IDs are Discord snowflakes; the earlier one has the smaller number.
	require.NoError(t, store.PutChallenge(ctx, "1052000000000000000", challenge.NewChallenge(
		stage.New("Sweet Lamb", location.WAL, stage.Long),
		weather.WET,
		car.New("Lancia Delta S4", class.GroupB4WD),
		nil, nil,
	)))
	require.NoError(t, store.PutChallenge(ctx, "1200000000000000000", challenge.NewChallenge(
		stage.New("Hamra", location.SWE, stage.Short),
		weather.SNOW,
		car.New("Peugeot 205 GTI", class.H2FWD),
		nil, nil,
	)))

	body := getPage(t, store, "/")

	assert.Contains(t, body, "Wales » Sweet Lamb")
	assert.Contains(t, body, "8 Sector")
	assert.Contains(t, body, "Wet")
	assert.Contains(t, body, "Lancia Delta S4")
	assert.Contains(t, body, "Sweden » Hamra")
	assert.Contains(t, body, "Snow")
	assert.Contains(t, body, "2 stored")

	// Newest challenge is listed first.
	assert.Less(t, strings.Index(body, "Hamra"), strings.Index(body, "Sweet Lamb"))
}

func TestChallengesPageWithNoChallenges(t *testing.T) {
	body := getPage(t, memorystore.New(), "/")

	assert.Contains(t, body, "No challenges stored yet.")
	assert.Contains(t, body, "0 stored")
}

func TestChallengesPageShowsCompletionsPanel(t *testing.T) {
	store := storeWithCompletions(t)

	body := getPage(t, store, "/")

	assert.Contains(t, body, "2 completions")
	assert.Contains(t, body, "Alice")
	assert.Contains(t, body, "Bob")
	assert.Contains(t, body, "1:30.000")
	// The panel is embedded in the initial page, inside the collapsed <details>.
	assert.Contains(t, body, `id="completions-`+challengeID+`"`)
	assert.Contains(t, body, "/htmx.min.js")
}

func TestCompletionsFragmentSortsByTime(t *testing.T) {
	store := storeWithCompletions(t)

	body := getPage(t, store, "/challenges/"+challengeID+"/completions?sort=time&dir=asc")

	// Bob (1:30) is faster than Alice (2:00), so he comes first.
	assert.Less(t, strings.Index(body, "Bob"), strings.Index(body, "Alice"))
}

func TestCompletionsFragmentSortsBySubmissionDescending(t *testing.T) {
	store := storeWithCompletions(t)

	body := getPage(t, store, "/challenges/"+challengeID+"/completions?sort=submitted&dir=desc")

	// Alice submitted after Bob, so newest-first puts her on top.
	assert.Less(t, strings.Index(body, "Alice"), strings.Index(body, "Bob"))
}

func TestCompletionsFragmentUnknownChallenge(t *testing.T) {
	server := httptest.NewServer(web.Handler(memorystore.New()))
	t.Cleanup(server.Close)

	response, err := http.Get(server.URL + "/challenges/nope/completions")
	require.NoError(t, err)
	t.Cleanup(func() { _ = response.Body.Close() })

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestServesHTMX(t *testing.T) {
	server := httptest.NewServer(web.Handler(memorystore.New()))
	t.Cleanup(server.Close)

	response, err := http.Get(server.URL + "/htmx.min.js")
	require.NoError(t, err)
	t.Cleanup(func() { _ = response.Body.Close() })

	require.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, response.Header.Get("Content-Type"), "javascript")

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), `version:"2.0.4"`)
}

func storeWithCompletions(t *testing.T) *memorystore.Store {
	t.Helper()

	bob := challenge.NewCompletionAt("bob-id", "Bob", 90*time.Second, time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC))
	alice := challenge.NewCompletionAt("alice-id", "Alice", 120*time.Second, time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC))

	store := memorystore.New()
	require.NoError(t, store.PutChallenge(context.Background(), challengeID, challenge.NewChallenge(
		stage.New("Hamra", location.SWE, stage.Short),
		weather.SNOW,
		car.New("Peugeot 205 GTI", class.H2FWD),
		[]challenge.Completion{alice, bob},
		nil,
	)))

	return store
}

func getPage(t *testing.T, store *memorystore.Store, path string) string {
	t.Helper()

	server := httptest.NewServer(web.Handler(store))
	t.Cleanup(server.Close)

	response, err := http.Get(server.URL + path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = response.Body.Close() })

	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}
