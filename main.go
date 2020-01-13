package main

import (
	"github.com/vsantos/steward-go/pkg/cognition"
	"github.com/vsantos/steward-go/pkg/system"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"log"
)

func main() {
	file, err := system.TemporaryFile("/tmp/", "voice_result.*.mp3")
	if err != nil {
		log.Fatal(err)
	}

	// Instantiates a Cognition
	vcs := cognition.Settings{
		Language:      "pt-BR",
		Text:          "Texto aleatório",
		Gender:        texttospeechpb.SsmlVoiceGender_FEMALE,
		Path:          file.Name(),
		Chmod:         0644,
		AudioEncoding: texttospeechpb.AudioEncoding_MP3,
	}

	var spk cognition.Steward
	spk = &cognition.Speech{Settings: vcs}

	err = spk.Read()
	if err != nil {
		log.Fatal(err)
	}

	err = spk.Say(vcs.Path)
	if err != nil {
		log.Fatal(err)
	}

}
