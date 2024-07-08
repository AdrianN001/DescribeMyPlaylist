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

type GenrePageTemplateParams struct {
	FavoriteGenre 				string
	SongsFromFavGenre			int
	SecondFavoriteGenre			string
	SongsFromSecondFavGenre		int


	NumberOfGenres				int

	BackgroundMusicPreviewUrl  	string

	BackgroundSongArtist		string
	BackgroundSongTitle			string
	BackgroundSongURL			string
}

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
	

	result, err := rating.RateGenreDiversityFromSaved(r.Context(), user)
	if err != nil{
		log.Fatal(err)
	}
	
	var template_args GenrePageTemplateParams
	template_args.FromRating(result)

	genre_file := path.Join("static","view","genres.html")
	template, err := template.ParseFiles(genre_file)
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


func (template_args *GenrePageTemplateParams) FromRating(genre_rating rating.GenreDiversityRating){
	template_args.BackgroundMusicPreviewUrl = genre_rating.SampleSong.PreviewURL
	template_args.FavoriteGenre = genre_rating.FavouriteGenre
	template_args.SongsFromFavGenre = genre_rating.NumberOfSongsFromFavGenre
	
	template_args.SongsFromSecondFavGenre = genre_rating.NumberOfSongsFromSecFavGenre
	template_args.SecondFavoriteGenre = genre_rating.SecondFavouriteGenre
	
	template_args.NumberOfGenres = genre_rating.NumberOfGenres

	template_args.BackgroundSongArtist = genre_rating.SampleSong.Artists[0].Name
	template_args.BackgroundSongTitle = genre_rating.SampleSong.Name
	template_args.BackgroundSongURL = genre_rating.SampleSong.ExternalURLs["spotify"]
}