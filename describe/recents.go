package describe

import (
	"context"
	"rate_my_playlist/wrapper"

	"github.com/zmb3/spotify/v2"
)

type RecentTopArtistsAndTracks struct {
	TopArtists []spotify.FullArtist
	TopTracks  []spotify.FullTrack
}

func DescribeRecentTopArtistsAndTracks(ctx context.Context, user wrapper.User, time_range spotify.Range) (RecentTopArtistsAndTracks, error) {
	current_user, err := user.CurrentUser(ctx)
	if err != nil {
		return RecentTopArtistsAndTracks{}, err
	}
	client := user.Client()

	artists, err := GetUserTopArtist(ctx, spotify.ID(current_user.ID), client, time_range)
	if err != nil {
		return RecentTopArtistsAndTracks{}, err
	}
	tracks, err := GetUserTopTracks(ctx, spotify.ID(current_user.ID), client, time_range)
	if err != nil {
		return RecentTopArtistsAndTracks{}, err
	}

	return RecentTopArtistsAndTracks{
		TopArtists: artists,
		TopTracks:  tracks,
	}, nil
}
