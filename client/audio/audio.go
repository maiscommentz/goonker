package audio

import (
	"Goonker/client/assets"
	"bytes"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// Constants for audio playback
const (
	SampleRate  = 44100
	MusicVolume = 0.3
	SoundVolume = 0.3
)

// AudioManager manages audio playback
type AudioManager struct {
	context *audio.Context
	players map[string]*audio.Player
	assets  fs.FS
}

// NewAudioManager creates a new audio manager
func NewAudioManager() *AudioManager {
	return &AudioManager{
		context: audio.NewContext(SampleRate),
		players: make(map[string]*audio.Player),
		assets:  assets.AssetsFS,
	}
}

// LoadMusic loads a mp3 music file and creates a player for it
func (am *AudioManager) LoadMusic(name, path string) error {
	// Read the file
	binData, err := fs.ReadFile(am.assets, path)
	if err != nil {
		return err
	}

	// Decode the MP3 file
	dataStream, err := mp3.DecodeWithSampleRate(SampleRate, bytes.NewReader(binData))
	if err != nil {
		return err
	}

	// Create a loop
	loop := audio.NewInfiniteLoop(dataStream, dataStream.Length())

	// Create a player
	player, err := am.context.NewPlayer(loop)
	if err != nil {
		return err
	}
	// Set the volume
	player.SetVolume(MusicVolume)

	// Add the player to the map
	am.players[name] = player
	return nil
}

// LoadSound loads a wav sound file and creates a player for it
func (am *AudioManager) LoadSound(name, path string) error {
	// Read the file
	binData, err := fs.ReadFile(am.assets, path)
	if err != nil {
		return err
	}

	// Decode the wav file
	dataStream, err := wav.DecodeWithSampleRate(SampleRate, bytes.NewReader(binData))
	if err != nil {
		return err
	}

	// Create a player
	player, err := am.context.NewPlayer(dataStream)
	if err != nil {
		return err
	}
	// Set the volume
	player.SetVolume(SoundVolume)

	am.players[name] = player
	return nil
}

// Play plays the audio
func (am *AudioManager) Play(name string) {
	if p, ok := am.players[name]; ok {
		if !p.IsPlaying() {
			err := p.Rewind()
			if err != nil {
				log.Printf("Error rewinding audio '%s': %v", name, err)
			}
			p.Play()
		}
	} else {
		log.Printf("Audio '%s' not found", name)
	}
}

// Stop stops the audio
func (am *AudioManager) Stop(name string) {
	if p, ok := am.players[name]; ok {
		if p.IsPlaying() {
			p.Pause()
			err := p.Rewind()
			if err != nil {
				log.Printf("Error rewinding audio '%s': %v", name, err)
			}
		}
	}
}

// Pause pauses the audio
func (am *AudioManager) Pause(name string) {
	if p, ok := am.players[name]; ok {
		if p.IsPlaying() {
			p.Pause()
		}
	}
}

// IsPlaying checks if an audio is playing
func (am *AudioManager) IsPlaying(name string) bool {
	if p, ok := am.players[name]; ok {
		return p.IsPlaying()
	}
	return false
}
