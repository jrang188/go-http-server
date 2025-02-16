package main

import (
	"log"
	"net/http"

	poker "github.com/jrang188/go-poker"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()
	server := poker.NewPlayerServer(store)
	if err := http.ListenAndServe(":3333", server); err != nil {
		log.Fatal(http.ListenAndServe(":3333", server))
	}
}
