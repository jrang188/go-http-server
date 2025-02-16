package poker_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	poker "github.com/jrang188/go-poker"
)

const jsonContentType = "application/json"

func TestGETPlayers(t *testing.T) {
	store := poker.StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := poker.NewPlayerServer(&store)
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

// server_test.go
func TestStoreWins(t *testing.T) {
	store := poker.StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := poker.NewPlayerServer(&store)

	t.Run("it records wins when POST", func(t *testing.T) {
		player := "Pepper"
		request := poker.NewPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusAccepted)

		if len(store.WinCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.WinCalls), 1)
		}
		if store.WinCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.WinCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	store := poker.StubPlayerStore{}
	server := poker.NewPlayerServer(&store)

	t.Run("it returns 200 on /league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		var got []poker.Player

		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
		}

		poker.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []poker.Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiesto", 14},
		}

		store := poker.StubPlayerStore{nil, nil, wantedLeague}
		server := poker.NewPlayerServer(&store)

		request := poker.NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertContentType(t, response, jsonContentType)

		got := poker.GetLeagueFromResponse(t, response.Body)
		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertLeague(t, got, wantedLeague)
	})
}
