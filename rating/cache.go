package rating

import (
	"context"
	"log"
	"rate_my_playlist/wrapper"
	"sync"

	"github.com/zmb3/spotify/v2"
)

var (
	PlaylistTracksCache map[string][]spotify.PlaylistTrack = make(map[string][]spotify.PlaylistTrack)
	SavedTracksCache    map[string][]spotify.SavedTrack    = make(map[string][]spotify.SavedTrack)

	ArtistCache 		map[spotify.ID]spotify.FullArtist  = make(map[spotify.ID]spotify.FullArtist)
	AudioFeatureCache 	map[spotify.ID]spotify.AudioFeatures  = make(map[spotify.ID]spotify.AudioFeatures)
)

var(
	PlaylistTracksCacheMutex   	sync.Mutex = sync.Mutex{}
	SavedTracksCacheMutex   	sync.Mutex = sync.Mutex{}
	ArtistCacheMutex   			sync.Mutex = sync.Mutex{}
	AudioFeatureCacheMutex   	sync.Mutex = sync.Mutex{}
)


func GetArtists(ctx context.Context, user wrapper.User, artist_ids []spotify.ID) ([]spotify.FullArtist, error) {


	var missing_artists []spotify.ID = make([]spotify.ID, 0, len(artist_ids))
	
	for _, artist_id := range artist_ids{
		if _, ok := ArtistCache[artist_id]; !ok{
			missing_artists = append(missing_artists, artist_id)
		}
	}

	
	if len(missing_artists) == 0{
		var full_artists []spotify.FullArtist = make([]spotify.FullArtist, 0, len(artist_ids))

		for _, artist_id := range artist_ids{
			full_artists = append(full_artists, ArtistCache[artist_id])
		}
		return full_artists, nil
	}

	fetched_missing_artists, err := FetchArtists(ctx, user, missing_artists)
	if err != nil{
		return []spotify.FullArtist{}, err
	}
	var full_artists []spotify.FullArtist = make([]spotify.FullArtist, 0, len(artist_ids))

	/* Place the newly fetched artists into the cache */
	for _, fetched_artist := range fetched_missing_artists{
		ArtistCacheMutex.Lock()
		ArtistCache[fetched_artist.ID] = *fetched_artist
		ArtistCacheMutex.Unlock()



		full_artists = append(full_artists, *fetched_artist)
	}

	for _, artist_id := range artist_ids{
		full_artists = append(full_artists, ArtistCache[artist_id])
	}
	return full_artists, nil
		
}



func GetAudioFeatures(ctx context.Context, user wrapper.User, song_ids []spotify.ID) ([]spotify.AudioFeatures, error) {


	var missing_audio_features []spotify.ID = make([]spotify.ID, 0, len(song_ids))
	
	for _, artist_id := range song_ids{
		if _, ok := AudioFeatureCache[artist_id]; !ok{
			missing_audio_features = append(missing_audio_features, artist_id)
		}
	}

	
	if len(missing_audio_features) == 0{
		var audio_features []spotify.AudioFeatures = make([]spotify.AudioFeatures, 0, len(song_ids))

		for _, artist_id := range song_ids{
			audio_features = append(audio_features, AudioFeatureCache[artist_id])
		}
		return audio_features, nil
	}

	fetched_missing_audio_features, err := FetchSongsAudioFeatures(ctx, user, missing_audio_features)
	if err != nil{
		return []spotify.AudioFeatures{}, err
	}
	var audio_features []spotify.AudioFeatures = make([]spotify.AudioFeatures, 0, len(song_ids))

	/* Place the newly fetched artists into the cache */
	for _, fetched_artist := range fetched_missing_audio_features{
		AudioFeatureCacheMutex.Lock()
		AudioFeatureCache[fetched_artist.ID] = *fetched_artist
		AudioFeatureCacheMutex.Unlock()



		audio_features = append(audio_features, *fetched_artist)
	}

	for _, artist_id := range song_ids{
		audio_features = append(audio_features, AudioFeatureCache[artist_id])
	}
	return audio_features, nil
		
}




func GetTracksFromSaved(ctx context.Context, user wrapper.User, cache_key string) ([]spotify.SavedTrack, error) {
	var cached_playlist []spotify.SavedTrack

	current_user, err := user.CurrentUser(ctx)
	if err != nil{
		return nil, err
	}

	cached_playlist, ok := SavedTracksCache[cache_key]
	if !ok{
		fetched_playlist, err := FetchSavedTracks(ctx, user)
		if err != nil{
			return nil, err
		}
		SavedTracksCacheMutex.Lock()
		SavedTracksCache[cache_key] = fetched_playlist
		SavedTracksCacheMutex.Unlock()
		cached_playlist = fetched_playlist
	}else{
		log.Printf("[$$$] CACHE HIT with user: %s [$$$]\n", current_user.ID)
	}

	return cached_playlist, nil

}

//TODO Implement it
func GetTracksFromPlaylist(){
	
}
