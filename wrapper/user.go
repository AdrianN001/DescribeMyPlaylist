package wrapper

import (
	"context"
	"net/http"

	spotify "github.com/zmb3/spotify/v2"
)

type User struct {
	spotify_client_instance *spotify.Client
}

func (user *User) Init(http_client *http.Client) *User {
	user.spotify_client_instance = spotify.New(http_client, spotify.WithRetry(false))
	return user
}

func (user *User) Client() *spotify.Client {
	return user.spotify_client_instance
}

func (user User) CurrentUser(ctx context.Context) (*spotify.PrivateUser, error) {
	return user.spotify_client_instance.CurrentUser(ctx)
}

func (user User) Playlists(ctx context.Context) ([]spotify.SimplePlaylist, error) {
	playlist, err := user.spotify_client_instance.CurrentUsersPlaylists(ctx, spotify.Limit(30))
	if err != nil {
		return nil, err
	}
	return playlist.Playlists, nil
}

func (user User) SavedTracks(ctx context.Context, opts ...spotify.RequestOption) (*spotify.SavedTrackPage, error) {
	return user.spotify_client_instance.CurrentUsersTracks(ctx, opts...)
}
