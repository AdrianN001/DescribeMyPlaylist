package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"rate_my_playlist/describe"
	"rate_my_playlist/wrapper"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type PopularityPageTemplate struct {
	MostPopularSongName  string
	LeastPopularSongName string
	AverageRating        int
	DistributionList     string
	Overall              string

	MostPopularSongPreviewUrl string
	BackgroundSongArtist      string
	BackgroundSongTitle       string
	BackgroundSongURL         string
}

func PopularityPageHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("[@] Request to '/popularity' [@]")

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
	result, err := describe.RatePopularityOfSavedTracks(r.Context(), user)
	if err != nil {
		return
	}

	template_args := ConvertRatingToPageTemplateArgs(result)
	popularity_file_path := path.Join("static", "view", "popularity.html")
	template, err := template.ParseFiles(popularity_file_path)
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

func ConvertRatingToPageTemplateArgs(popularity_rating describe.PopularityRating) PopularityPageTemplate {
	most_popular_song_artist := popularity_rating.MostPopularSong.Artists[0].Name
	least_popular_song_artist := popularity_rating.LeastPopularSong.Artists[0].Name

	most_popular_song_name := popularity_rating.MostPopularSong.Name
	least_popular_song_name := popularity_rating.LeastPopularSong.Name

	most_popular_song_full := fmt.Sprintf("%s-%s", most_popular_song_artist, most_popular_song_name)
	least_popular_song_full := fmt.Sprintf("%s-%s", least_popular_song_artist, least_popular_song_name)

	distribution_list_str_repr, err := json.Marshal(popularity_rating.Distribution)
	if err != nil {
		log.Fatal(err)
	}

	preview_song := popularity_rating.MostPopularSong

	return PopularityPageTemplate{
		MostPopularSongName:       most_popular_song_full,
		LeastPopularSongName:      least_popular_song_full,
		AverageRating:             popularity_rating.AveragePopularity,
		Overall:                   popularity_rating.Overall(),
		MostPopularSongPreviewUrl: preview_song.PreviewURL,
		BackgroundSongArtist:      preview_song.Artists[0].Name,
		BackgroundSongTitle:       preview_song.Name,
		DistributionList:          string(distribution_list_str_repr),
		BackgroundSongURL:         preview_song.ExternalURLs["spotify"],
	}

}
