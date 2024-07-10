package rating

import (
	"context"
	"fmt"
	"rate_my_playlist/wrapper"

	"github.com/zmb3/spotify/v2"
)

type PopularityRating struct {
	MostPopularSong  spotify.SavedTrack
	LeastPopularSong spotify.SavedTrack

	Distribution      []int
	AveragePopularity int
}

func RatePopularityOfSavedTracks(ctx context.Context, user wrapper.User) (PopularityRating, error) {
	current_user, err := user.CurrentUser(ctx)
	if err != nil {
		return PopularityRating{}, err
	}
	cache_key := fmt.Sprintf("%s-%s", current_user.ID, "saved-tracks")
	cached_playlist, err := GetTracksFromSaved(ctx, user, cache_key)
	if err != nil {
		return PopularityRating{}, err
	}
	var least_popular_song_value int = 110
	var least_popular_song spotify.SavedTrack

	var most_popular_song_value int = -1
	var most_popular_song spotify.SavedTrack

	var sum_popularity_value int = 0

	for _, song := range cached_playlist {
		popularity_value := int(song.Popularity)

		if popularity_value < least_popular_song_value {
			least_popular_song = song
			least_popular_song_value = popularity_value
		}
		if popularity_value > most_popular_song_value {
			most_popular_song = song
			most_popular_song_value = popularity_value

		}
		sum_popularity_value += popularity_value
	}

	avg_popularity := sum_popularity_value / len(cached_playlist)

	return PopularityRating{
		MostPopularSong:   most_popular_song,
		LeastPopularSong:  least_popular_song,
		Distribution:      GetPopularityDistributionOfSaved(cached_playlist),
		AveragePopularity: avg_popularity,
	}, nil
}

func RatePopularityOfPlaylist(ctx context.Context, user wrapper.User, playlist_id string) {

}

func GetPopularityDistributionOfSaved(songs []spotify.SavedTrack) []int {
	var distribution_list []int = make([]int, 5)

	for _, song := range songs {
		popularity := song.Popularity

		distribution_place := int(popularity / 20)

		distribution_list[distribution_place]++
	}
	return distribution_list
}

func (rate PopularityRating) Overall() string {

	if 80 <= rate.AveragePopularity && rate.AveragePopularity < 100 {
		return "According to your playlist, it appears that you enjoy listening to popular music. Your song selections indicate a preference for well-known and widely loved tracks."
	} else if 60 <= rate.AveragePopularity && rate.AveragePopularity < 80 {
		return "Most of the songs in your playlist are popular hits, but you also have a taste for undiscovered gems. Your music collection includes a mix of mainstream and lesser-known tracks."
	} else if 40 <= rate.AveragePopularity && rate.AveragePopularity < 60 {
		return "Your playlist mostly consists of non-mainstream hits, with a few exceptions. It is evident that you prefer music that aligns with your personal taste."
	} else if 20 <= rate.AveragePopularity && rate.AveragePopularity < 40 {
		return "Based on your playlist, it appears that you do not prefer popular songs and instead gravitate towards more unique and lesser-known tracks. Your music taste seems to diverge greatly from mainstream trends, as evidenced by the songs you have chosen to include in your playlist."
	}
	return "Based on your playlist, it is evident that you have a unique taste in music and do not necessarily follow the latest popular hits. You prefer to stay true to your own preferences and style when it comes to choosing songs. Your playlist reflects your individuality and independence in selecting music."

}
