package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"rate_my_playlist/rating"
	"rate_my_playlist/wrapper"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type ArtistPageTemplateParameters struct {
	FavoriteArtistName       string
	SecondFavoriteArtistName string

	FavoriteArtistSongCount       int
	SecondFavoriteArtistSongCount int

	RemainingSongCount int

	Overview string

	BackgroundMusicPreviewUrl string
	BackgroundSongArtist      string
	BackgroundSongTitle       string
	BackgroundSongURL         string
}

func ArtistPageHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("[@] Request to '/artist' [@]")

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
	result, err := rating.RateArtistDiversityFromSaved(r.Context(), user)
	if err != nil {
		return
	}

	overview := result.Overall()

	template_args := ArtistPageTemplateParameters{
		Overview:                 overview,
		FavoriteArtistName:       result.FavouriteArtistName,
		SecondFavoriteArtistName: result.SecondFavouriteArtistName,

		FavoriteArtistSongCount:       result.FavouriteArtistSongsN,
		SecondFavoriteArtistSongCount: result.SecondFavouriteArtistSongsN,

		RemainingSongCount: result.AllSongN - result.SecondFavouriteArtistSongsN - result.FavouriteArtistSongsN,

		BackgroundMusicPreviewUrl: result.RandomSongFromArtist.PreviewURL,

		BackgroundSongArtist: result.RandomSongFromArtist.Artists[0].Name,
		BackgroundSongTitle:  result.RandomSongFromArtist.Name,
		BackgroundSongURL:    result.RandomSongFromArtist.ExternalURLs["spotify"],
	}
	artist_file_path := path.Join("static", "view", "artist.html")
	template, err := template.ParseFiles(artist_file_path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = template.Execute(w, template_args)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
