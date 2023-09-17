package audio_engine

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Music rl.Music
var musicPaused = false

func initializeMusic() {
	Music = rl.LoadMusicStream("resources/audio/tracks/ting ting.mp3")
	// music = rl.LoadMusicStream("resources/audio/Underworld Coffee Shop.mp3")
}

func playMusic(m rl.Music) {
	rl.PlayMusicStream(m)
}

func unloadMusic() {
	rl.UnloadMusicStream(Music)
}

func ToggleMusic() {
	musicPaused = !musicPaused
	if musicPaused {
		PauseMusicStream()
	} else {
		ResumeMusicStream()
	}
}

func UpdateMusicStream() {
	rl.UpdateMusicStream(Music)
}

func PauseMusicStream() {
	rl.PauseMusicStream(Music)
}

func ResumeMusicStream() {
	rl.ResumeMusicStream(Music)
}
