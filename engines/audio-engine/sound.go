package audioengine

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var SwordSound rl.Sound

func initializeSounds() {
	SwordSound = rl.LoadSound("resources/audio/effects/swing-whoosh.mp3")
}

func unloadSounds() {
	rl.UnloadSound(SwordSound)
}

func PlaySound(s rl.Sound) {
	// would be good to make this more intelligent in the future, maybe a priority queue w/limited buffer
	rl.PlaySound(s)
}
