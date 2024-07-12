package describe

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"rate_my_playlist/utils"
	"rate_my_playlist/wrapper"
	"sort"
	"sync"

	"github.com/zmb3/spotify/v2"
)

type GenreDiversityRating struct {
	FavouriteGenre               string
	NumberOfSongsFromFavGenre    int
	SecondFavouriteGenre         string
	NumberOfSongsFromSecFavGenre int

	SampleSong spotify.SavedTrack

	NumberOfGenres int
	NumberOfSongs  int
}

func RateGenreDiversityFromSaved(ctx context.Context, user wrapper.User) (GenreDiversityRating, error) {
	current_user, err := user.CurrentUser(ctx)
	if err != nil {
		return GenreDiversityRating{}, err
	}
	cache_key := fmt.Sprintf("%s-%s", current_user.ID, "saved-tracks")
	cached_playlist, err := GetTracksFromSaved(ctx, user, cache_key)
	if err != nil {
		return GenreDiversityRating{}, err
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	var artists []spotify.FullArtist = make([]spotify.FullArtist, 500)

	for i := 0; i < int(len(cached_playlist)/20)*20; i += 20 {
		sub_slice := cached_playlist[i : i+20]
		sub_slice_ids := utils.Map[spotify.SavedTrack, spotify.ID](sub_slice, func(st spotify.SavedTrack) spotify.ID {
			return st.Artists[0].ID
		})
		wg.Add(1)

		go func(song_ids []spotify.ID, subslice_start, subslice_end int) {
			defer wg.Done()

			audio_features, err := GetArtists(ctx, user, song_ids)
			if err != nil {
				log.Fatalln(err)
			}

			mutex.Lock()
			artist_index := 0
			for j := subslice_start; j < subslice_end; j++ {
				artists[j] = audio_features[artist_index]
				artist_index++
			}
			mutex.Unlock()
		}(sub_slice_ids, i, i+20)
	}

	wg.Wait()
	artists = artists[:len(cached_playlist)]
	genre_frequencies := ArtistsToGenres(artists)
	rating := TopGenresFromFrequencies(genre_frequencies)

	sample_song, err := SearchSongFromSavedWithGenre(ctx, user, cached_playlist, rating.FavouriteGenre)
	if err != nil {
		return GenreDiversityRating{}, err
	}
	rating.SampleSong = sample_song

	rating.NumberOfSongs = len(cached_playlist)

	return rating, nil
}

func ArtistsToGenres(artists []spotify.FullArtist) map[string]int {

	var genre_occur map[string]int = map[string]int{}

	for _, artist := range artists {
		for _, genre := range artist.Genres {
			if _, ok := genre_occur[genre]; !ok {
				genre_occur[genre] = 1
			} else {
				genre_occur[genre]++
			}
		}
	}
	return genre_occur
}

func TopGenresFromFrequencies(genres map[string]int) GenreDiversityRating {

	number_of_genres := len(genres)

	var genres_keys []string = make([]string, 0, number_of_genres)

	for key := range genres {
		genres_keys = append(genres_keys, key)
	}

	sort.SliceStable(genres_keys, func(i, j int) bool {
		return genres[genres_keys[i]] < genres[genres_keys[j]]
	})

	favourite_genre := genres_keys[len(genres_keys)-1]
	second_favourite_genre := genres_keys[len(genres_keys)-2]

	return GenreDiversityRating{
		FavouriteGenre:            favourite_genre,
		NumberOfSongsFromFavGenre: genres[favourite_genre],

		SecondFavouriteGenre:         second_favourite_genre,
		NumberOfSongsFromSecFavGenre: genres[second_favourite_genre],

		NumberOfGenres: number_of_genres,
	}
}

func SearchSongFromSavedWithGenre(ctx context.Context, user wrapper.User, songs []spotify.SavedTrack, searched_genre string) (spotify.SavedTrack, error) {
	
	const max_depth int = 20
	depth := rand.Intn(max_depth)
	
	for _, song := range songs {

		main_artist := song.Artists[0].ID

		// At this point, it is 100% cached already
		full_artists, err := GetArtists(ctx, user, []spotify.ID{main_artist})
		if err != nil {
			return spotify.SavedTrack{}, err
		}
		full_artist := full_artists[0]

		for _, genre := range full_artist.Genres {
			if genre == searched_genre && depth == 0 {
				return song, nil
			}else if genre == searched_genre && depth > 0{
				depth--
			}
		}
	}

	return spotify.SavedTrack{}, fmt.Errorf("song not found")
}

func RateGenreDiversityFromPlaylist(ctx context.Context, user wrapper.User, playlist_id string) {

}
