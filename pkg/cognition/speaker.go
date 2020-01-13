package cognition

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"context"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	log "github.com/sirupsen/logrus"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"io/ioutil"
	"os"
	"time"
)

type Settings struct {
	Language      string
	Text          string
	Gender        texttospeechpb.SsmlVoiceGender
	Path          string
	KeepRecord    bool
	Chmod         os.FileMode
	AudioEncoding texttospeechpb.AudioEncoding
}

type Steward interface {
	Listen() error
	Say(recordedAudioPath string) error
	Read() error
}

type Speech struct {
	Settings Settings
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func (s *Speech) Say(recordedAudioPath string) error {
	f, err := os.Open(recordedAudioPath)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	defer streamer.Close()
	_ = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	log.Println("Starting recorded audio...")
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done

	defer func() {
		if ! s.Settings.KeepRecord {
			log.Infof("Removing voice result from path: %s", s.Settings.Path)
			err = os.Remove(recordedAudioPath)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	return nil
}

func (s *Speech) Read() error {
	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: s.Settings.Text},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: s.Settings.Language,
			SsmlGender:   s.Settings.Gender,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding:   s.Settings.AudioEncoding,
			SampleRateHertz: 32000,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	// The resp's AudioContent is binary.
	err = ioutil.WriteFile(s.Settings.Path, resp.AudioContent, s.Settings.Chmod)
	if err != nil {
		log.Fatal(err)
	}

	// cleaning memory
	resp.AudioContent = []byte{}

	log.Printf("Saved voice result to path: %s", s.Settings.Path)
	err = client.Close()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *Speech) Listen() error {
	return nil
}
