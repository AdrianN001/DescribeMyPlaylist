package describe

import (
	"context"
	"log"
	"rate_my_playlist/wrapper"
	"sync"

	"github.com/zmb3/spotify/v2"
)

// TODO Fix this (the saved track is "working")
func FetchPlaylistTracks(ctx context.Context, user wrapper.User, playlist_id string) ([]spotify.PlaylistTrack, error) {
	client := user.Client()
	track_page, err := client.GetPlaylistTracks(ctx, spotify.ID(playlist_id), spotify.Limit(50), spotify.Offset(0))
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var tracks []spotify.PlaylistTrack = make([]spotify.PlaylistTrack, 500)

	for i := 50; i == int(track_page.Total); i += 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			track_page, err := client.GetPlaylistTracks(ctx, spotify.ID(playlist_id), spotify.Limit(50), spotify.Offset(i))
			if err != nil {
				log.Fatalln(err)
			}
			mutex.Lock()
			tracks = append(tracks, track_page.Tracks...)
			mutex.Unlock()
		}()
	}
	wg.Wait()
	log.Println("[!] Fetch finished [!] ", len(tracks))
	return tracks, nil
}

func FetchSavedTracks(ctx context.Context, user wrapper.User) ([]spotify.SavedTrack, error) {
	track_page, err := user.SavedTracks(ctx, spotify.Limit(50), spotify.Offset(0))
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var tracks []spotify.SavedTrack = make([]spotify.SavedTrack, 0, 500)

	tracks = append(tracks, track_page.Tracks...)

	j := 0

	for i := 50; i < int(track_page.Total); i += 50 {
		wg.Add(1)
		go func(loop_index int) {

			defer wg.Done()
			track_page, err := user.SavedTracks(ctx, spotify.Limit(50), spotify.Offset(loop_index))
			if err != nil {
				log.Println(err)
				return
			}
			mutex.Lock()
			j += len(track_page.Tracks)
			tracks = append(tracks, track_page.Tracks...)
			mutex.Unlock()
		}(i)
	}
	wg.Wait()
	log.Println("[!] Fetch finished [!] ", j)
	return tracks, nil
}

func FetchArtists(ctx context.Context, user wrapper.User, artist_ids []spotify.ID) ([]*spotify.FullArtist, error) {
	client := user.Client()

	return client.GetArtists(ctx, artist_ids...)
}

func FetchSongsAudioFeatures(ctx context.Context, user wrapper.User, song_ids []spotify.ID) ([]*spotify.AudioFeatures, error) {
	client := user.Client()

	return client.GetAudioFeatures(ctx, song_ids...)
}
