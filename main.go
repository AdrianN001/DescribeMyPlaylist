package main

import (
	"log"
	"net/http"
	"rate_my_playlist/handlers"

	"github.com/joho/godotenv"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))

	godotenv.Load()
	http.HandleFunc("/callback", handlers.CallbackHandler)
	http.HandleFunc("/login", handlers.LoginPageHandler)
	http.HandleFunc("/playlists", handlers.PlaylistsPageHandler)
	http.HandleFunc("/popularity", handlers.PopularityPageHandler)
	http.HandleFunc("/artist", handlers.ArtistPageHandler)
	http.HandleFunc("/genre", handlers.GenrePageHandler)
	http.HandleFunc("/emotional", handlers.EmotionalPageHandler)
	http.HandleFunc("/musical_elements", handlers.MusicalElementPageHandler)

	
	http.HandleFunc("/{$}", handlers.HomePageHandler)

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("[@] Server Started as http://localhost:8080 [@]")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
