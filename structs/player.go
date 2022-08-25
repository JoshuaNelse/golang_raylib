package structs

import (
	"math"
	audioengine "raylib/playground/engines/audio-engine"
	"raylib/playground/game"
	util "raylib/playground/game/utils"
	"raylib/playground/structs/draw2d"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

type Test struct {
	Test   string
	Sprite Sprite
}

type Player struct {
	Sprite         Sprite
	SpriteFlipped  bool
	Obj            *resolv.Object
	Weapon         *Weapon
	Hand           Point
	Moving         bool
	Attacking      bool
	AttackCooldown int
}

func (p *Player) Move(dx, dy float64) {
	p.Obj.X += dx
	p.Obj.Y += dy
	p.Obj.Update()
	p.Sprite.Dest.X = util.RectFromObj(p.Obj).X
	p.Sprite.Dest.Y = util.RectFromObj(p.Obj).Y
	p.Weapon.Move(dx, dy)
}

func (p *Player) Draw() {
	if game.FrameCount%8 == 1 {
		p.Sprite.Frame++
	}
	if p.Sprite.Frame > 3 {
		p.Sprite.Frame = 0
	}
	var weaponOffset float32 = 0
	if p.Moving {
		p.Sprite.Src.X = 192                                                                       // pixel where run animation starts
		p.Sprite.Src.X += float32(p.Sprite.Frame) * float32(math.Abs(float64(p.Sprite.Src.Width))) // rolling the animation
		weaponOffset = -4
	} else {
		if rl.GetScreenWidth()/2 <= int(rl.GetMouseX()) {
			util.FlipRight(&p.Sprite.Src)
			p.SpriteFlipped = false
		} else {
			util.FlipLeft(&p.Sprite.Src)
			p.SpriteFlipped = true
		}
		p.Sprite.Src.X = 128                                                                       // pixel where rest idle starts
		p.Sprite.Src.X += float32(p.Sprite.Frame) * float32(math.Abs(float64(p.Sprite.Src.Width))) // rolling the animation
	}
	p.Weapon.SpriteFlipped = p.SpriteFlipped
	p.Moving = false
	rl.DrawTexturePro(draw2d.Texture, p.Sprite.Src, p.Sprite.Dest, rl.NewVector2(p.Sprite.Dest.Width, p.Sprite.Dest.Height), 0, rl.White)
	updateFrame := game.FrameCount%8 == 0
	p.Weapon.Draw(p.Sprite.Frame, updateFrame, weaponOffset)
}

func (p *Player) Attack() []Projectile {
	audioengine.PlaySound(audioengine.SwordSound)
	p.Weapon.AttackFrame = 0 // find a better way to trigger animation than this.
	p.AttackCooldown = p.Weapon.Cooldown

	playerCenter := Point{
		X: float32(p.Obj.X + p.Obj.W/2),
		Y: float32(p.Obj.Y + p.Obj.H/2),
	}
	rl.DrawCircleLines(int32(playerCenter.X), int32(playerCenter.Y), 32, rl.Green)
	angle := util.GetPlayerToMouseAngleDegress()

	// TODO use weapon attributes in the future to determine this logic
	projectileCount := p.Weapon.ProjectileCount
	projectileReach := p.Weapon.Projectilelength
	projectileSpread := p.Weapon.ProjectileSpreadDegrees
	projectileTTL := p.Weapon.ProjectileTTLFrames
	projectileVelocity := p.Weapon.ProjectileVelocity
	projectileSpreadItter := int(float64(angle) - math.Floor(float64(projectileCount)/2)*float64(projectileSpread))

	var newProjectiles []Projectile
	for i := 0; i < projectileCount; i++ {
		x2 := int(float64(projectileReach) * math.Sin(util.DegreesToRadians(float64(projectileSpreadItter))))
		y2 := int(float64(projectileReach) * math.Cos(util.DegreesToRadians(float64(projectileSpreadItter))))
		var projectileTrajectory float64
		if projectileVelocity > 0 {
			projectileTrajectory = float64(projectileSpreadItter)
		}
		newProjectile := Projectile{
			Start:      rl.NewVector2(playerCenter.X, playerCenter.Y),
			End:        rl.NewVector2(playerCenter.X+float32(x2), playerCenter.Y+float32(y2)),
			Ttl:        projectileTTL,
			Velocity:   projectileVelocity,
			Trajectory: projectileTrajectory,
			Sprite: Sprite{
				Src:  p.Weapon.ProjectileSpriteSrc.Src,
				Dest: p.Weapon.ProjectileSpriteSrc.Dest,
			},
		}
		newProjectiles = append(newProjectiles, newProjectile)
		projectileSpreadItter += projectileSpread
	}
	p.Attacking = false
	return newProjectiles
}

func (p *Player) EquipWeapon(w *Weapon) {
	// create new object from updated dest X/Y
	w.Sprite.Dest.X = p.Hand.X + float32(p.Obj.X)
	w.Sprite.Dest.Y = p.Hand.Y + float32(p.Obj.Y)
	w.Obj = util.ObjFromRect(w.Sprite.Dest)

	// update player weapon
	p.Weapon = w
}
