package rating

import (
	"context"
	"fmt"
	"log"
	"rate_my_playlist/utils"
	"rate_my_playlist/wrapper"
	"sync"

	"github.com/zmb3/spotify/v2"
)
type MusicalElementsRating struct {
	AverageTempo			float32
	TempoRating				string

	AverageDanceability  	float32
	DanceabilityRating   	string	

	AverageSpeechines		float32
	SpeechinesRating		string

	AveragePitch			int
	PitchRating				string


	BonusRating				string

	RandomSong              spotify.SimpleTrack
}

func RateMusicalElementOfSaved(ctx context.Context, user wrapper.User) (MusicalElementsRating, error){
	current_user, err := user.CurrentUser(ctx)
	if err != nil{
		return MusicalElementsRating{}, err
	}
	cache_key := fmt.Sprintf("%s-%s", current_user.ID, "saved-tracks")
	cached_playlist, err := GetTracksFromSaved(ctx, user, cache_key)
	if err != nil{
		return MusicalElementsRating{}, err
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	
	var song_moods []wrapper.MusicalElement = make([]wrapper.MusicalElement, 500)


	for i := 0; i < int(len(cached_playlist)/50)*50; i += 50{
		sub_slice := cached_playlist[i: i+50]
		sub_slice_ids := utils.Map[spotify.SavedTrack, spotify.ID](sub_slice, func(st spotify.SavedTrack) spotify.ID {
			return st.ID
		})
		wg.Add(1)

		go func(song_ids []spotify.ID, subslice_start, subslice_end int){
			defer wg.Done()

			audio_features, err := GetAudioFeatures(ctx, user, song_ids)
			if err != nil{
				log.Fatalln(err)
			}

			mutex.Lock()
			audio_features_index := 0
			for j := subslice_start; j < subslice_end; j++{
				song_moods[j] = wrapper.GetMusicalElementFromAudioFeatures(audio_features[audio_features_index], cached_playlist[j].SimpleTrack)
				audio_features_index++
			}
			mutex.Unlock()
		}(sub_slice_ids, i, i+50)
	}

	wg.Wait()
	song_moods = song_moods[: len(cached_playlist)]

	average := CalculateAverageMusicalElement(song_moods)
	rating := create_rating_from_average(average)
	return rating,  nil
}


func CalculateAverageMusicalElement(songs []wrapper.MusicalElement) wrapper.MusicalElement{
	sum_of_tempo 				:= 0.0
	sum_of_speechiness			:= 0.0					
	sum_of_accousticness 		:= 0.0
	sum_of_danceability			:= 0.0
	sum_of_keys					:= 0			// Probably shouldn't do average, but check the most common key
	for _, song := range songs{
		sum_of_tempo += float64(song.Tempo)
		sum_of_speechiness += float64(song.Speechiness)
		sum_of_accousticness += float64(song.Accousticness)
		sum_of_danceability += float64(song.Danceability)
		sum_of_keys 		+= song.Key
	}

	avg_tempo := sum_of_tempo / float64(len(songs))
	avg_speechiness := sum_of_speechiness / float64(len(songs))
	avg_accousticness := sum_of_accousticness / float64(len(songs))
	avg_danceability := sum_of_danceability / float64(len(songs))
	avg_key			 :=  sum_of_keys/len(songs)

	random_song := utils.RandElement(songs).Song
	for random_song.PreviewURL == ""{
		random_song = utils.RandElement(songs).Song
	}

	return wrapper.MusicalElement{
		Tempo: float32(avg_tempo),
		Speechiness: float32(avg_speechiness),
		Accousticness: float32(avg_accousticness),
		Danceability: float32(avg_danceability),

		Song: random_song,
		Key: avg_key,
	}
}

func create_rating_from_average(average wrapper.MusicalElement ) MusicalElementsRating{
	bonus_sentence, tempo_overview, danceability_overview, speechiness_overview := average.Overall()

	return MusicalElementsRating{
		AverageTempo: average.Tempo,
		TempoRating: tempo_overview,

		AverageDanceability: average.Danceability,
		DanceabilityRating: danceability_overview,

		AverageSpeechines: average.Speechiness,
		SpeechinesRating: speechiness_overview,

		AveragePitch: average.Key,
		PitchRating: wrapper.PitchClass[average.Key],

		BonusRating: bonus_sentence,
		RandomSong: average.Song,
	}

}