package soundengine

import rl "github.com/gen2brain/raylib-go/raylib"

var SwordSound rl.Sound

func InitializeSounds() {
	SwordSound = rl.LoadSound("resources/audio/effects/swing-whoosh.mp3")
}

func UnloadSounds() {
	rl.UnloadSound(SwordSound)
}

func PlaySound(s rl.Sound) {
	// would be good to make this more intelligent in the future, maybe a priority queue w/limited buffer
	rl.PlaySound(s)
}
