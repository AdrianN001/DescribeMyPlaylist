package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"rate_my_playlist/wrapper"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type PlaylistRendererArgs struct {
	Playlists []spotify.SimplePlaylist
}


func PlaylistsPageHandler(w http.ResponseWriter, r *http.Request) {
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
	
	user_playlists, err := user.Playlists(r.Context())
	log.Println(len(user_playlists))
	if err != nil{
		log.Fatal(err)
	}


	template_args := PlaylistRendererArgs{user_playlists}
	playlist_file_path := path.Join("static","view","playlists.html")
	template, err := template.ParseFiles(playlist_file_path)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	err = template.Execute(w, template_args)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
}