package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"rate_my_playlist/rating"
	"rate_my_playlist/wrapper"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)



func GenrePageHandler(w http.ResponseWriter, r *http.Request) {
	// session,_ := store.Get(r, "spotify-code")


	session, err := store.Get(r, "spotify-code")
	if err != nil{
		log.Fatal(err)
	}
	serialized_token, ok := session.Values["token"].([]byte) 
	if !ok{
		log.Fatal(ok)
	}

	var token oauth2.Token 

	err = json.Unmarshal(serialized_token, &token )
	if err != nil{
		log.Fatal(err)
	}
	
	
	
	auth := spotifyauth.New( spotifyauth.WithClientID(os.Getenv("CLIENT_ID")),
	spotifyauth.WithRedirectURL(os.Getenv("REDIRECT_URI")), 
	spotifyauth.WithClientSecret(os.Getenv("CLIENT_SECRET")),
	spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserLibraryRead,  spotifyauth.ScopeUserTopRead, spotifyauth.ScopePlaylistReadPrivate, spotifyauth.ScopePlaylistReadCollaborative))
	
	var user wrapper.User
	user.Init(auth.Client(r.Context(), &token))
	

	rating.RateEmotionalVibeOfSaved(r.Context(), user)


	w.Write([]byte("asd23"))
}