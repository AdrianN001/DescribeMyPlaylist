package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"rate_my_playlist/describe"
	"rate_my_playlist/utils"
	"rate_my_playlist/wrapper"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type RecentPagePayload struct {

	TopArtists []spotify.FullArtist		`json:"artists"`
	TopTracks  []spotify.FullTrack		`json:"tracks"`

	BackgroundMusicPreviewUrl string	`json:"preview_url"`
}

func RecentPageRequestHandler(w http.ResponseWriter, r *http.Request){
	log.Println("[@] Request to '/get_recent' [@]")

	session, err := store.Get(r, "spotify-code")
	if err != nil {
		log.Fatal(err)
	}
	serialized_token, ok := session.Values["token"].([]byte)
	if !ok {
		log.Fatal(ok)
	}

	var token oauth2.Token

	err = json.Unmarshal(serialized_token, &token)
	if err != nil {
		log.Fatal(err)
	}

	auth := spotifyauth.New(spotifyauth.WithClientID(os.Getenv("CLIENT_ID")),
		spotifyauth.WithRedirectURL(os.Getenv("REDIRECT_URI")),
		spotifyauth.WithClientSecret(os.Getenv("CLIENT_SECRET")),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserLibraryRead, spotifyauth.ScopeUserTopRead, spotifyauth.ScopePlaylistReadPrivate, spotifyauth.ScopePlaylistReadCollaborative))

	var user wrapper.User
	user.Init(auth.Client(r.Context(), &token))

	time_range := r.URL.Query().Get("time_range")


	recents, err := describe.DescribeRecentTopArtistsAndTracks(r.Context(), user, spotify.Range(time_range))
	if err != nil{
		log.Fatal(err)
	}
	
	payload := RecentPagePayload{
		TopArtists: recents.TopArtists,
		TopTracks: recents.TopTracks,

		BackgroundMusicPreviewUrl: utils.RandElement[spotify.FullTrack](recents.TopTracks).PreviewURL,
	}

	payload_buffer, err := json.Marshal(payload)
	if err != nil{
		log.Fatal(err)
	}
	w.Write(payload_buffer)
}



func RecentPageHandler(w http.ResponseWriter, r *http.Request){
	log.Println("[@] Request to '/recents' [@]")

	session, err := store.Get(r, "spotify-code")
	if err != nil {
		log.Fatal(err)
	}
	serialized_token, ok := session.Values["token"].([]byte)
	if !ok {
		log.Fatal(ok)
	}

	var token oauth2.Token

	err = json.Unmarshal(serialized_token, &token)
	if err != nil {
		log.Fatal(err)
	}

	auth := spotifyauth.New(spotifyauth.WithClientID(os.Getenv("CLIENT_ID")),
		spotifyauth.WithRedirectURL(os.Getenv("REDIRECT_URI")),
		spotifyauth.WithClientSecret(os.Getenv("CLIENT_SECRET")),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserLibraryRead, spotifyauth.ScopeUserTopRead, spotifyauth.ScopePlaylistReadPrivate, spotifyauth.ScopePlaylistReadCollaborative))

	var user wrapper.User
	user.Init(auth.Client(r.Context(), &token))

	genre_file := path.Join("static", "view", "recents.html")
	template, err := template.ParseFiles(genre_file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = template.Execute(w, struct{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
