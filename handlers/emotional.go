package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"rate_my_playlist/rating"
	"rate_my_playlist/wrapper"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)




type EmotionalPageTemplateParameters struct{

	HappiestSongName				string
	SaddestSongName					string
	MostRelaxingSongName			string
	AngriestSongName				string


	NumberOfHappySong				int
	NumberOfSadSong					int
	NumberOfRelaxingSong			int
	NumberOfAngrySong				int

	Overview 						string
	BackgroundMusicPreviewUrl		string

	BackgroundSongArtist		string
	BackgroundSongTitle			string
	BackgroundSongURL			string
}

func EmotionalPageHandler(w http.ResponseWriter, r *http.Request) {
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
	

	result, err := rating.RateEmotionalVibeOfSaved(r.Context(), user)
	if err != nil{
		return 
	}

	template_args := template_params_from_audio_features(result)
	emotional_file := path.Join("static","view","emotional.html")
	template, err := template.ParseFiles(emotional_file)
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

func template_params_from_audio_features(audio_rating rating.EmotionalRating) EmotionalPageTemplateParameters{

	var max_magnitude_value float64 = math.Sqrt(0.5*0.5 + 0.5*0.5)

	angriest_song_full := 			fmt.Sprintf("%s - %s (%.1f %% angry)", audio_rating.AngriestSong.Song.Artists[0].Name ,audio_rating.AngriestSong.Song.Name,( audio_rating.AngriestSong.Magnitude / float64(max_magnitude_value)) * 100)
	most_relaxing_song_full := 		fmt.Sprintf("%s - %s (%.1f %% relaxing)", audio_rating.MostRelaxingSong.Song.Artists[0].Name ,audio_rating.MostRelaxingSong.Song.Name, ( audio_rating.MostRelaxingSong.Magnitude / float64(max_magnitude_value)) * 100)
	happiest_song_full := 			fmt.Sprintf("%s - %s (%.1f %% happy)", audio_rating.HappiestSong.Song.Artists[0].Name ,audio_rating.HappiestSong.Song.Name, ( audio_rating.HappiestSong.Magnitude / float64(max_magnitude_value)) * 100)
	saddest_song_full := 			fmt.Sprintf("%s - %s (%.1f %% sad)", audio_rating.SaddestSong.Song.Artists[0].Name ,audio_rating.SaddestSong.Song.Name,  ( audio_rating.SaddestSong.Magnitude / float64(max_magnitude_value)) * 100)

	return EmotionalPageTemplateParameters{
		AngriestSongName: angriest_song_full,
		MostRelaxingSongName: most_relaxing_song_full,
		HappiestSongName: happiest_song_full,
		SaddestSongName: saddest_song_full,

		NumberOfHappySong:  audio_rating.NumberOfHappySong,
		NumberOfSadSong:    audio_rating.NumberOfSadSong,
		NumberOfAngrySong: audio_rating.NumberOfAngrySong,
		NumberOfRelaxingSong: audio_rating.NumberOfAngrySong,

		BackgroundMusicPreviewUrl: audio_rating.HappiestSong.Song.PreviewURL,
		

		BackgroundSongArtist: audio_rating.HappiestSong.Song.Artists[0].Name,
		BackgroundSongTitle: audio_rating.HappiestSong.Song.Name,
		BackgroundSongURL: audio_rating.HappiestSong.Song.ExternalURLs["spotify"],
		Overview: "Baszo izelesed van",
	}
}