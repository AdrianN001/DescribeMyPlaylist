package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"rate_my_playlist/rating"
	"rate_my_playlist/wrapper"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type MusicalElementPageTemplateParameters struct {
	AverageTempo string
	AveragePitch string

	AverageDanceability string
	AverageSingability  string

	Overviews []string

	BackgroundMusicPreviewUrl string
	BackgroundSongArtist      string
	BackgroundSongTitle       string
	BackgroundSongURL         string
}

func MusicalElementPageHandler(w http.ResponseWriter, r *http.Request) {
	// session,_ := store.Get(r, "spotify-code")

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

	result, err := rating.RateMusicalElementOfSaved(r.Context(), user)
	if err != nil {
		return
	}

	var template_args MusicalElementPageTemplateParameters
	template_args.FromRating(result)

	emotional_file := path.Join("static", "view", "musical_element.html")
	template, err := template.ParseFiles(emotional_file)
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

func (template_params *MusicalElementPageTemplateParameters) FromRating(musical_rating rating.MusicalElementsRating) {
	average_tempo := fmt.Sprintf("%.1f BPM", musical_rating.AverageTempo)
	average_pitch := musical_rating.PitchRating

	average_danceability := fmt.Sprintf("%.1f %%", musical_rating.AverageDanceability*100)
	average_singability := fmt.Sprintf("%.1f %%", musical_rating.AverageSpeechines*100)

	overviews := []string{musical_rating.SpeechinesRating, musical_rating.TempoRating, musical_rating.DanceabilityRating}
	if musical_rating.BonusRating != "" {
		overviews = append(overviews, musical_rating.BonusRating)
	}

	rand.Shuffle(len(overviews), func(i, j int) {
		overviews[i], overviews[j] = overviews[j], overviews[i]
	})

	template_params.AverageTempo = average_tempo
	template_params.AveragePitch = average_pitch

	template_params.AverageDanceability = average_danceability
	template_params.AverageSingability = average_singability

	template_params.Overviews = overviews
	template_params.BackgroundMusicPreviewUrl = musical_rating.RandomSong.PreviewURL

	template_params.BackgroundSongArtist = musical_rating.RandomSong.Artists[0].Name
	template_params.BackgroundSongTitle = musical_rating.RandomSong.Name
	template_params.BackgroundSongURL = musical_rating.RandomSong.ExternalURLs["spotify"]
}
