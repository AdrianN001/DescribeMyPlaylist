package describe

import (
	"context"
	"fmt"
	"math/rand"
	"rate_my_playlist/wrapper"
	"sort"

	"github.com/zmb3/spotify/v2"
)

type ArtistDiversityRating struct {
	FavouriteArtistName         string
	FavouriteArtistSongsN       int
	SecondFavouriteArtistName   string
	SecondFavouriteArtistSongsN int

	AllSongN int

	RandomSongFromArtist spotify.FullTrack
}

func RateArtistDiversityFromSaved(ctx context.Context, user wrapper.User) (ArtistDiversityRating, error) {
	current_user, err := user.CurrentUser(ctx)
	if err != nil {
		return ArtistDiversityRating{}, err
	}
	cache_key := fmt.Sprintf("%s-%s", current_user.ID, "saved-tracks")
	cached_playlist, err := GetTracksFromSaved(ctx, user, cache_key)
	if err != nil {
		return ArtistDiversityRating{}, err
	}

	var artist_occurences map[string][]spotify.SavedTrack = make(map[string][]spotify.SavedTrack)

	for _, song := range cached_playlist {
		for _, artist := range song.Artists {
			if _, ok := artist_occurences[artist.Name]; !ok {
				artist_occurences[artist.Name] = make([]spotify.SavedTrack, 0, 50)
				artist_occurences[artist.Name] = append(artist_occurences[artist.Name], song)
			} else {
				artist_occurences[artist.Name] = append(artist_occurences[artist.Name], song)
			}
		}
	}

	keys := make([]string, 0, len(artist_occurences))
	for key := range artist_occurences {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return len(artist_occurences[keys[i]]) < len(artist_occurences[keys[j]])
	})

	random_index := rand.Int() % len(artist_occurences[keys[len(keys)-1]])

	random_song_from_top_artist := artist_occurences[keys[len(keys)-1]][random_index]

	return ArtistDiversityRating{
		FavouriteArtistName:       keys[len(keys)-1],
		SecondFavouriteArtistName: keys[len(keys)-2],

		FavouriteArtistSongsN:       len(artist_occurences[keys[len(keys)-1]]),
		SecondFavouriteArtistSongsN: len(artist_occurences[keys[len(keys)-2]]),

		AllSongN:             len(cached_playlist),
		RandomSongFromArtist: random_song_from_top_artist.FullTrack,
	}, nil
}

func (rating ArtistDiversityRating) Overall() string {

	// The artist has the 40% of the playlist
	if (rating.FavouriteArtistSongsN/rating.AllSongN)*100 >= 40 {

		return fmt.Sprintf("You are the biggest fan of %s, since this artist(s) makes up at least 40%% of your playlist.", rating.FavouriteArtistName)
	}

	// The artist has the 25% of the playlist
	if (rating.FavouriteArtistSongsN/rating.AllSongN)*100 >= 25 {

		return fmt.Sprintf("The quarter of your playlist was created by %s. It's not an exaggeration to say that you like this artist.", rating.FavouriteArtistName)
	}

	// The artist has the 10% of the playlist
	if (rating.FavouriteArtistSongsN/rating.AllSongN)*100 >= 10 {

		return fmt.Sprintf("You are definitely a fan of %s, considering that this artist(s) makes up at least 10%% of your playlist", rating.FavouriteArtistName)
	}

	if rating.FavouriteArtistSongsN != rating.SecondFavouriteArtistSongsN {

		return fmt.Sprintf("It was challenging to identify just one artist from the numerous ones you enjoy listening to, but %s has the highest number of saved tracks.", rating.FavouriteArtistName)
	}

	return fmt.Sprintf("From looking at your playlist, it is not apparent who your favorite artist is. However, %s and %s are the top two artists with the most saved tracks.", rating.FavouriteArtistName, rating.SecondFavouriteArtistName)

}
