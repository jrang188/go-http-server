package poker_test

import (
	"strings"
	"testing"

	poker "github.com/jrang188/go-poker"
)

func TestCLI(t *testing.T) {
	t.Run("record 'Chris win' from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record 'Tracy win' from user input", func(t *testing.T) {
		in := strings.NewReader("Tracy wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Tracy")
	})
}
