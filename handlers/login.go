package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var (
    // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
    key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)
)

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[@] Request to '/login' [@]")

	auth := spotifyauth.New( spotifyauth.WithClientID(os.Getenv("CLIENT_ID")),
						spotifyauth.WithRedirectURL(os.Getenv("REDIRECT_URI")), 
						spotifyauth.WithClientSecret(os.Getenv("CLIENT_SECRET")),
						spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserLibraryRead,  spotifyauth.ScopeUserTopRead, spotifyauth.ScopePlaylistReadPrivate, spotifyauth.ScopePlaylistReadCollaborative))


	url := auth.AuthURL("123")

	fmt.Println(url)

	http.Redirect(w, r, url, 200)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request){
	log.Println("[@] Request to '/callback' [@]")

	auth := spotifyauth.New( spotifyauth.WithClientID(os.Getenv("CLIENT_ID")),
						spotifyauth.WithRedirectURL(os.Getenv("REDIRECT_URI")), 
						spotifyauth.WithClientSecret(os.Getenv("CLIENT_SECRET")),
						spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserLibraryRead,  spotifyauth.ScopeUserTopRead, spotifyauth.ScopePlaylistReadPrivate, spotifyauth.ScopePlaylistReadCollaborative))

	token, err := auth.Token(r.Context(), "123", r)

	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	



	store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
		Secure: true,
	}


	json_token, err := json.Marshal(token)
	if err != nil{
		log.Fatal(err)
	}
	

	session,_ := store.Get(r, "spotify-code")
	session.Values["token"] = json_token
	err = session.Save(r,w)
	if err != nil{
		log.Fatal(err)
	}

	http.Redirect(w, r, "/playlists", 200)



}