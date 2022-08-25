package audioengine

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func InitializeAudio() {
	rl.InitAudioDevice()

	initializeSounds()
	initializeMusic()

	// Play music should be dynamic in the future but for now we will just play this song
	playMusic(Music)
}

func UnloadAudioComponents() {
	unloadSounds()
	unloadMusic()
}
