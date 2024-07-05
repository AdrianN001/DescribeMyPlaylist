package wrapper

import (
	"fmt"

	"github.com/zmb3/spotify/v2"
)

var (
	danceability_options = []string{
		"Your music doesn't have a rhythm or beat that is easy to dance to. It is difficult to sync your movements with the music.",
		"It's hard for me to picture someone dancing to this music.",
		"Your songs strike a delicate balance between being suitable for dancing and not.", 
		"Most of your songs are very well suited for dancing.",
		"You're songs are made to dance to.",
	}

	speechiness_options = []string{
		"You're not a fan of lyrics, but you enjoy music with a strong beat.",
		"You rarely listen to music with lyrics that are easy to follow and sing along to.",
		"You enjoy music that emphasizes the beat rather than relying on spoken lyrics",
		"You love songs which have an understandble lyrics, while they have a good beat, too.",
		"You prioritize lyrics that are easy to comprehend.",
	}

	tempo_options = []string{

	}
)

type MusicalElement struct {
	Song 				spotify.SimpleTrack

	Tempo 				float32
	Speechiness			float32					// how clear the lyrics is, and how much of the music is made out of words  “Speechiness detects the presence of spoken words in a track”.
	Accousticness 		float32
	Danceability		float32
}

func GetMusicalElementFromAudioFeatures(audio_feature spotify.AudioFeatures, song spotify.SimpleTrack) MusicalElement {


	return MusicalElement{
		Song: song,
		Speechiness: audio_feature.Speechiness,
		Accousticness: audio_feature.Acousticness,
		Danceability: audio_feature.Danceability,
		Tempo: audio_feature.Tempo,
	}
}


func (average_value MusicalElement) Overall() []string{
	var bonus_sentence string = ""
	if average_value.Tempo > 100 && average_value.Speechiness <= 0.66 && average_value.Danceability > 0.66{
		bonus_sentence = "The songs you prefer are ones that are upbeat and great for dancing to. You enjoy fast-paced music that gets you moving. "
	}else if average_value.Tempo < 90 && average_value.Speechiness >= 0.66{
		bonus_sentence = "You're more into slow(ish) songs, with understandable lyrics."
	}else if average_value.Tempo > 90 && average_value.Speechiness >= 0.66 {
		bonus_sentence = "You're not really into slow-paced songs, but rather quick beats that don't rely on lyrics"
	}
	fmt.Printf("%+v\n", average_value)
	
	// tempo_overview := tempo_options[int(average_value.Tempo /20)]
	danceability_overview := danceability_options[int(average_value.Danceability /.2)]
	speechiness_overview  := speechiness_options[int(average_value.Speechiness   /.2)] 

	return []string{
		bonus_sentence, 
		danceability_overview,
		speechiness_overview,
	}
}