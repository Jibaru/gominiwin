package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jibaru/gominiwin"
	"github.com/jibaru/gominiwin/colors"
	"github.com/jibaru/gominiwin/keys"
)

type Sprite struct {
	data      [][]bool
	x, y      float32
	blockSize float32
}

func NewSprite(filename string, x, y, blockSize float32) (*Sprite, error) {
	dino := &Sprite{x: x, y: y, blockSize: blockSize}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	file.Close()

	for _, line := range strings.Split(string(data), "\n") {
		l := make([]bool, 0)
		for _, v := range strings.Split(strings.TrimSpace(line), " ") {
			l = append(l, v == "x")
		}
		dino.data = append(dino.data, l)
	}

	return dino, nil
}

func (d *Sprite) CollidesWith(other *Sprite) bool {
	for i, row := range d.data {
		for j, cell := range row {
			if cell {
				dCellX := d.x + float32(j)*d.blockSize
				dCellY := d.y + float32(i)*d.blockSize

				for k, otherRow := range other.data {
					for l, otherCell := range otherRow {
						if otherCell {
							oCellX := other.x + float32(l)*other.blockSize
							oCellY := other.y + float32(k)*other.blockSize

							if dCellX < oCellX+other.blockSize &&
								dCellX+d.blockSize > oCellX &&
								dCellY < oCellY+other.blockSize &&
								dCellY+d.blockSize > oCellY {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func (d *Sprite) Draw(w gominiwin.Win) {
	for i, row := range d.data {
		for j, cell := range row {
			if cell {
				left := d.x + float32(j)*d.blockSize
				top := d.y + float32(i)*d.blockSize
				right := left + d.blockSize
				bottom := top + d.blockSize

				w.FilledRectangle(left, top, right, bottom)
			}
		}
	}
}

func main() {
	const startY float32 = 320
	const initialSpeed float32 = 17

	dino, err := NewSprite("dino.txt", 50, startY, 2)
	if err != nil {
		panic(err)
	}

	obstacleTemplate, err := NewSprite("obstacle.txt", 200, 344, 2)
	if err != nil {
		panic(err)
	}

	gameOver, err := NewSprite("over.txt", 100, 90, 15)
	if err != nil {
		panic(err)
	}

	w, err := gominiwin.New("DinoGame", 800, 500)
	if err != nil {
		panic(err)
	}

	obstacleSpeed := float32(8)
	jumpSpeed := initialSpeed
	gravity := float32(1)
	onGround := true

	obstacles := []*Sprite{
		{data: obstacleTemplate.data, x: 800, y: 344, blockSize: 2},
		{data: obstacleTemplate.data, x: 1100, y: 344, blockSize: 2},
		{data: obstacleTemplate.data, x: 1400, y: 344, blockSize: 2},
	}
	score := 0.0

	go func() {
		for {
			score += 0.1

			w.Clear()

			// Background
			w.SetColor(colors.White)
			w.FilledRectangle(0, 0, float32(w.Width()), float32(w.Height()))

			// Score
			w.SetColor(colors.Black)
			w.SetText(10, 10, fmt.Sprintf("SCORE: %v", int(score)))

			// Ground
			w.SetColorRGB(255, 153, 51)
			w.FilledRectangle(0, float32(w.Height()-100), float32(w.Width()), float32(w.Height()))

			// Dino
			w.SetColorRGB(128, 128, 128)
			dino.Draw(w)

			// Obstacles
			w.SetColorRGB(119, 179, 0)
			for i := 0; i < len(obstacles); i++ {
				obstacles[i].Draw(w)

				if dino.CollidesWith(obstacles[i]) {
					// Game over
					w.SetColor(colors.Red)
					gameOver.Draw(w)
					w.Refresh()
					return
				}

				obstacles[i].x -= obstacleSpeed

				if obstacles[i].x+obstacles[i].blockSize*float32(len(obstacles[i].data[0])) < 0 {
					obstacles[i].x = 800 + (rand.Float32() * 200)
				}
			}

			// Jumps
			if w.KeyPressed() == keys.Space && onGround {
				onGround = false
			}

			if !onGround {
				dino.y -= jumpSpeed
				jumpSpeed -= gravity

				if dino.y >= startY {
					dino.y = startY
					onGround = true
					jumpSpeed = initialSpeed
				}
			}

			w.Refresh()
			time.Sleep(5 * time.Millisecond)
		}
	}()

	w.Start()
}
