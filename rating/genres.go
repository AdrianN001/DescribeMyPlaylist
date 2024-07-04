package rating

import (
	"context"
	"fmt"
	"log"
	"rate_my_playlist/utils"
	"rate_my_playlist/wrapper"
	"sort"
	"sync"

	"github.com/zmb3/spotify/v2"
)

type GenreDiversityRating struct {
	FavouriteGenre				string
	SecondFavouriteGenre 		string

	LeastFavouriteGenre			string

	NumberOfGenres				int
}

func RateGenreDiversityFromSaved(ctx context.Context, user wrapper.User) (GenreDiversityRating, error){
	current_user, err := user.CurrentUser(ctx)
	if err != nil{
		return GenreDiversityRating{}, err
	}
	cache_key := fmt.Sprintf("%s-%s", current_user.ID, "saved-tracks")
	cached_playlist, err := GetTracksFromSaved(ctx, user, cache_key)
	if err != nil{
		return GenreDiversityRating{}, err
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	
	var should_continue bool = true
	var artists []spotify.FullArtist = make([]spotify.FullArtist, 0, 500)

	var buffer []spotify.ID = make([]spotify.ID, 0, 25)

	for _, song := range cached_playlist{
		if !should_continue{
			break
		}
		artist_ids := utils.Map[spotify.SimpleArtist, spotify.ID](song.Artists, func (artist spotify.SimpleArtist) spotify.ID {
			return artist.ID
		})
		buffer = append(buffer, artist_ids...)
		if len(buffer) <= 15{
			continue
		}
		
		log.Println(buffer)
		wg.Add(1)
		go func( ids *[]spotify.ID) {
			defer wg.Done()
			if !should_continue{
				return
			}
			fetched_buffered_artists, err := GetArtists(ctx, user, *ids)
			if err != nil{
				should_continue = false
				log.Println(err.Error())
				//TODO stays like this untill higher api rate allowed
				return
			}
			mutex.Lock()
			artists = append(artists, fetched_buffered_artists...)
			(*ids) = (*ids)[:0]
			mutex.Unlock()
		}( &buffer)
	}
	wg.Wait()
	log.Println(len(artists))
	genre_frequencies := ArtistsToGenres(artists)
	TopGenresFromFrequencies(genre_frequencies)
	return GenreDiversityRating{}, nil
}

func ArtistsToGenres(artists []spotify.FullArtist) map[string]int {

	var genre_occur map[string]int = map[string]int{}

	for _, artist := range artists{
		for _, genre := range artist.Genres{
			if _, ok := genre_occur[genre]; !ok {
				genre_occur[genre] = 1
			}else{
				genre_occur[genre]++
			}
		}
	}
	return genre_occur
}

func TopGenresFromFrequencies(genres map[string]int) GenreDiversityRating{

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

	fmt.Printf("favourite_genre: %v\n", favourite_genre)
	fmt.Printf("second_favourite_genre: %v\n", second_favourite_genre)
	fmt.Printf("number_of_genres: %v\n", number_of_genres)

	return GenreDiversityRating{}
}

func RateGenreDiversityFromPlaylist(ctx context.Context, user wrapper.User, playlist_id string){

}