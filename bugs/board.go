package bugs

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type square struct {
	bug   bool
	shown bool
	near  int
}

type board struct {
	height, width int

	win, lose func()

	bugcount int
	bugs     [][]square
}

func (b *board) countHidden() int {
	count := 0

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			item := b.bugs[y][x]

			if !item.shown {
				count++
			}
		}
	}

	return count
}

func (b *board) load(count int) {
	b.bugcount = count
	if b.bugs == nil {
		b.bugs = make([][]square, b.height)

		for y := 0; y < b.height; y++ {
			b.bugs[y] = make([]square, b.width)
		}
	} else {
		for y := 0; y < b.height; y++ {
			for x := 0; x < b.width; x++ {
				b.bugs[y][x].shown = false
				b.bugs[y][x].bug = false
				b.bugs[y][x].near = 0
			}
		}
	}

	for i := 0; i < count; i++ {
		x := rand.Intn(b.width)
		y := rand.Intn(b.height)

		if b.bugs[y][x].bug {
			i--
		} else {
			b.setMine(x, y)
		}
	}
}

func (b *board) incSquare(x, y int) {
	if x < 0 || y < 0 {
		return
	}
	if x >= b.width || y >= b.height {
		return
	}

	if b.bugs[y][x].bug {
		return
	}
	b.bugs[y][x].near++
}

func (b *board) setMine(x, y int) {
	if b.bugs[y][x].bug {
		return
	}
	b.bugs[y][x].bug = true

	b.incSquare(x-1, y-1)
	b.incSquare(x, y-1)
	b.incSquare(x+1, y-1)

	b.incSquare(x-1, y)
	b.incSquare(x+1, y)

	b.incSquare(x-1, y+1)
	b.incSquare(x, y+1)
	b.incSquare(x+1, y+1)
}

func (b *board) reveal(x, y int) {
	if x < 0 || y < 0 {
		return
	}
	if x >= b.width || y >= b.height {
		return
	}

	sq := b.bugs[y][x]
	if sq.shown {
		return
	}
	b.bugs[y][x].shown = true

	if sq.bug {
		if b.lose != nil {
			b.lose()
		}
		return
	}

	if sq.near == 0 {
		b.reveal(x-1, y-1)
		b.reveal(x, y-1)
		b.reveal(x+1, y-1)
		b.reveal(x-1, y)
		b.reveal(x+1, y)
		b.reveal(x-1, y+1)
		b.reveal(x, y+1)
		b.reveal(x+1, y+1)
	}

	if b.countHidden() == b.bugcount && b.win != nil {
		b.win()
	}
}

func squareString(sq square) string {
	if !sq.shown {
		return "?"
	} else if sq.bug {
		return "*"
	} else if sq.near == 0 {
		return " "
	}

	return fmt.Sprintf("%d", sq.near)
}

func (b *board) String() string {
	buf := strings.Builder{}
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			sq := b.bugs[y][x]

			buf.WriteString(squareString(sq))
		}
		buf.WriteByte('\n')
	}

	return buf.String()
}

func newBoard(height, width int) *board {
	rand.Seed(time.Now().Unix())
	return &board{height: height, width: width}
}
