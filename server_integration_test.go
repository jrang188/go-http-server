package poker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	poker "github.com/jrang188/go-poker"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := poker.CreateTempFile(t, "[]")
	defer cleanDatabase()
	store, err := poker.NewFileSystemPlayerStore(database)
	poker.AssertNoError(t, err)
	server, err := poker.NewPlayerServer(store, dummyGame)
	if err != nil {
		t.Fatalf("problem creating player server: %v", err)
	}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), poker.NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), poker.NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), poker.NewPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, poker.NewGetScoreRequest(player))
		poker.AssertStatus(t, response, http.StatusOK)

		poker.AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, poker.NewLeagueRequest())
		poker.AssertStatus(t, response, http.StatusOK)

		got := poker.GetLeagueFromResponse(t, response.Body)
		want := []poker.Player{
			{"Pepper", 3},
		}
		poker.AssertLeague(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := poker.NewFileSystemPlayerStore(database)

		poker.AssertNoError(t, err)
	})
}
