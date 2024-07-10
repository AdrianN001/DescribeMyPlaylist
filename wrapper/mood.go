package wrapper

import (
	"math"

	"github.com/zmb3/spotify/v2"
)

const (
	HAPPY_MOOD = iota
	SAD_MOOD

	RELAXED_MOOD
	ANGRY_MOOD

	Moodles = -1
)

type SongMood struct {
	Song spotify.SimpleTrack
	Mood int

	Magnitude float64
}

func GetSongMoodFromAudioFeature(audio_feature spotify.AudioFeatures, song spotify.SimpleTrack) SongMood {

	valence := audio_feature.Valence
	energy := audio_feature.Energy

	is_energetic_mood := energy >= 0.5
	is_happy_mood := valence >= 0.5

	mood := calculate_mood(is_energetic_mood, is_happy_mood)

	magnitude := calculate_magnitude(float64(valence), float64(energy))

	return SongMood{
		Mood:      mood,
		Magnitude: magnitude,
		Song:      song,
	}
}

func calculate_mood(is_energetic_mood, is_happy_mood bool) int {
	if is_energetic_mood && is_happy_mood {
		return HAPPY_MOOD
	} else if is_energetic_mood && !is_happy_mood {
		return ANGRY_MOOD
	} else if !is_energetic_mood && is_happy_mood {
		return RELAXED_MOOD
	}
	return SAD_MOOD
}

func calculate_magnitude(valence, energy float64) float64 {

	/* Valence and energy ranges from 0.0 to 1.0 */
	/* Now, it ranges from -.5 to .5 */

	valence -= .5
	energy -= .5

	return math.Sqrt(valence*valence + energy*energy)
}
