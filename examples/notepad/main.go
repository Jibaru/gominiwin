package main

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/jibaru/gominiwin"
	"github.com/jibaru/gominiwin/colors"
	"github.com/jibaru/gominiwin/keys"
)

func main() {
	const height float32 = 600
	const pointSize float32 = 3.0
	const horizontalPointsPerLetter float32 = 10.0
	const verticalPointsPerLetter float32 = 9.0
	const horizontalLetterSize float32 = horizontalPointsPerLetter * pointSize
	const verticalLetterSize float32 = verticalPointsPerLetter * pointSize
	const cursorCode = keys.Space + 1
	var heightPos float32 = 0

	w, err := gominiwin.New("Notepad", 800, int(height))
	if err != nil {
		panic(err)
	}

	chars := map[int]*Letter{}
	parseData("char.txt", 65, chars, pointSize)
	parseData("num.txt", 48, chars, pointSize)
	parseData("other.txt", keys.Space, chars, pointSize)

	cursorLetter := chars[cursorCode]
	typedLetters := []*Letter{}

	cursorLetter.x = 0
	cursorLetter.y = 0
	cursorFreq := &Freq{counter: 0, activeUntil: 70}

	hasSpace := func(w gominiwin.Win) bool {
		nextWordSize := cursorLetter.x + horizontalLetterSize
		return nextWordSize <= float32(w.Width())
	}

	shouldAddKey := func(k int, w gominiwin.Win) bool {
		return k != keys.None && (isNum(k) || isLetter(k) || k == keys.Space) && hasSpace(w)
	}

	go func() {
		for {
			w.Clear()

			// Background
			w.SetColor(colors.White)
			w.FilledRectangle(0, 0, float32(w.Width()), float32(w.Height()))

			k := w.KeyPressed()

			if k == keys.Return {
				heightPos += verticalLetterSize + 10
				cursorLetter.x = 0
				cursorLetter.y = heightPos
			}

			if shouldAddKey(k, w) {
				sprite, ok := chars[k]
				if ok {
					newLetter := *sprite
					newLetter.x = cursorLetter.x
					newLetter.y = cursorLetter.y

					gap := float32(2)
					cursorLetter.x = newLetter.x + horizontalLetterSize + gap
					cursorLetter.y = heightPos
					typedLetters = append(typedLetters, &newLetter)
				}

			}

			if cursorFreq.IsActive() {
				w.SetColor(colors.Red)
				cursorLetter.Draw(w)
			}

			w.SetColor(colors.Black)
			for _, letter := range typedLetters {
				letter.Draw(w)
			}

			w.Refresh()
			cursorFreq.Inc()
			time.Sleep(5 * time.Millisecond)
		}
	}()

	w.Start()
}

type Letter struct {
	spriteData [][]bool
	x, y       float32
	blockSize  float32
}

func NewLetter(data string, x, y, blockSize float32) (*Letter, error) {
	letter := &Letter{x: x, y: y, blockSize: blockSize}

	for _, line := range strings.Split(data, "\n") {
		l := make([]bool, 0)
		for _, v := range strings.Split(strings.TrimSpace(line), "") {
			l = append(l, v == "x")
		}
		letter.spriteData = append(letter.spriteData, l)
	}

	return letter, nil
}

func (d *Letter) Draw(w gominiwin.Win) {
	for i, row := range d.spriteData {
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

func isNum(k int) bool {
	return k >= 48 && k <= 57
}

func isLetter(k int) bool {
	return k >= 65 && k <= 90
}

func readFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	file.Close()
	return string(data)
}

func parseData(filename string, i int, chars map[int]*Letter, pointSize float32) {
	for _, d := range strings.Split(string(readFile(filename)), "\n\n") {
		sprite, err := NewLetter(d, 0, 0, pointSize)
		if err != nil {
			panic(err)
		}
		chars[i] = sprite
		i++
	}
}

type Freq struct {
	counter     int
	activeUntil int
}

func (f *Freq) Inc() {
	f.counter++
	f.resetIfReached()
}

func (f *Freq) IsActive() bool {
	return f.counter >= 0 && f.counter <= f.activeUntil
}

func (f *Freq) resetIfReached() {
	if f.counter > f.activeUntil {
		f.counter = -f.activeUntil
	}
}
