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

type EmotionalRating struct {
	HappiestSong     wrapper.SongMood
	AngriestSong     wrapper.SongMood
	MostRelaxingSong wrapper.SongMood
	SaddestSong      wrapper.SongMood

	NumberOfHappySong    int
	NumberOfAngrySong    int
	NumberOfRelaxingSong int
	NumberOfSadSong      int
}

func RateEmotionalVibeOfSaved(ctx context.Context, user wrapper.User) (EmotionalRating, error) {
	current_user, err := user.CurrentUser(ctx)
	if err != nil {
		return EmotionalRating{}, err
	}
	cache_key := fmt.Sprintf("%s-%s", current_user.ID, "saved-tracks")
	cached_playlist, err := GetTracksFromSaved(ctx, user, cache_key)
	if err != nil {
		return EmotionalRating{}, err
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	var song_moods []wrapper.SongMood = make([]wrapper.SongMood, 500)

	for i := 0; i < int(len(cached_playlist)/50)*50; i += 50 {
		sub_slice := cached_playlist[i : i+50]
		sub_slice_ids := utils.Map[spotify.SavedTrack, spotify.ID](sub_slice, func(st spotify.SavedTrack) spotify.ID {
			return st.ID
		})
		wg.Add(1)

		go func(song_ids []spotify.ID, subslice_start, subslice_end int) {
			defer wg.Done()

			audio_features, err := GetAudioFeatures(ctx, user, song_ids)
			if err != nil {
				log.Fatalln(err)
			}

			mutex.Lock()
			audio_features_index := 0
			for j := subslice_start; j < subslice_end; j++ {
				song_moods[j] = wrapper.GetSongMoodFromAudioFeature(audio_features[audio_features_index], cached_playlist[j].SimpleTrack)
				audio_features_index++
			}
			mutex.Unlock()
		}(sub_slice_ids, i, i+50)
	}

	wg.Wait()
	song_moods = song_moods[:len(cached_playlist)]

	rating := SearchForTheEdgeSongs(song_moods)

	fmt.Printf("%+v\n", rating)
	return rating, nil
}

func SearchForTheEdgeSongs(songs []wrapper.SongMood) EmotionalRating {

	angriest_song := wrapper.SongMood{Mood: wrapper.Moodles, Magnitude: -1}
	happiest_song := wrapper.SongMood{Mood: wrapper.Moodles, Magnitude: -1}
	most_relaxing_song := wrapper.SongMood{Mood: wrapper.Moodles, Magnitude: -1}
	saddest_song := wrapper.SongMood{Mood: wrapper.Moodles, Magnitude: -1}

	number_of_angry_song := 0
	number_of_happy_song := 0
	number_of_sad_song := 0
	number_of_relaxing_song := 0

	for _, song := range songs {

		switch song.Mood {
		case wrapper.HAPPY_MOOD:
			number_of_happy_song++
			if song.Magnitude > happiest_song.Magnitude {
				happiest_song = song
			}
		case wrapper.SAD_MOOD:
			number_of_sad_song++
			if song.Magnitude > saddest_song.Magnitude {
				saddest_song = song
			}
		case wrapper.RELAXED_MOOD:
			number_of_relaxing_song++
			if song.Magnitude > most_relaxing_song.Magnitude {
				most_relaxing_song = song
			}
		case wrapper.ANGRY_MOOD:
			number_of_angry_song++
			if song.Magnitude > angriest_song.Magnitude {
				angriest_song = song
			}
		}
	}

	return EmotionalRating{
		HappiestSong:     happiest_song,
		SaddestSong:      saddest_song,
		MostRelaxingSong: most_relaxing_song,
		AngriestSong:     angriest_song,

		NumberOfHappySong:    number_of_happy_song,
		NumberOfAngrySong:    number_of_angry_song,
		NumberOfRelaxingSong: number_of_relaxing_song,
		NumberOfSadSong:      number_of_sad_song,
	}
}
