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
	overviews := average.Overall()

	log.Println(overviews)
	return MusicalElementsRating{}, nil
}


func CalculateAverageMusicalElement(songs []wrapper.MusicalElement) wrapper.MusicalElement{
	sum_of_tempo 				:= 0.0
	sum_of_speechiness			:= 0.0					
	sum_of_accousticness 		:= 0.0
	sum_of_danceability			:= 0.0
	for _, song := range songs{
		sum_of_tempo += float64(song.Tempo)
		sum_of_speechiness += float64(song.Speechiness)
		sum_of_accousticness += float64(song.Accousticness)
		sum_of_danceability += float64(song.Danceability)
	}

	avg_tempo := sum_of_tempo / float64(len(songs))
	avg_speechiness := sum_of_speechiness / float64(len(songs))
	avg_accousticness := sum_of_accousticness / float64(len(songs))
	avg_danceability := sum_of_danceability / float64(len(songs))

	return wrapper.MusicalElement{
		Tempo: float32(avg_tempo),
		Speechiness: float32(avg_speechiness),
		Accousticness: float32(avg_accousticness),
		Danceability: float32(avg_danceability),
	}
}