package quadtree

import (
	"errors"
	"log"

	"github.com/njaros/QuadTree/helpers"
)

type Point interface {
	X() int
	Y() int
}

type QuadTree[T Point] struct {
	WestPole       int
	EastPole       int
	NorthPole      int
	SouthPole      int
	Width int
	Height int
	Census         int
	pointDirectory map[*T]*qtNode[T]
	root           *qtNode[T]
}

// minX, maxX, minY, maxY defines the border of the map, it is necessary to
// efficiently subdivide the map.
func NewQuadTree[T Point](minX int, maxX int, minY int, maxY int) (*QuadTree[T], error) {
	if minX >= maxX || minY >= maxY {
		return nil, errors.New("QuadTree borders wrongly defined.")
	}

	return &QuadTree[T]{
		WestPole:       minX,
		EastPole:       maxX,
		NorthPole:      minY,
		SouthPole:      maxY,
		Height: maxY - minY + 1,
		Width: maxX - minX + 1,
		Census:         0,
		pointDirectory: make(map[*T]*qtNode[T]),
		root:           nil,
	}, nil
}

// qtNode may contains several points if they are superimposed (in same coordinates).
// Structure to store points in hashmap to ease to Point reference search.
type qtNode[T Point] struct {
	parent *qtNode[T]
	depth int
	NW     *qtNode[T]
	NE     *qtNode[T]
	SE     *qtNode[T]
	SW     *qtNode[T]
	MinX   int
	MaxX   int
	MinY   int
	MaxY   int
	X      int
	Y      int
	Census int
	P      map[*T]bool
}

func newqtNode[T Point](minX int, maxX int, minY int, maxY int, parent *qtNode[T]) *qtNode[T] {
	depth := 0
	if parent != nil {
		depth = (*parent).depth + 1
	}
	return &qtNode[T]{
		parent: parent,
		depth: depth,
		NW:     nil,
		NE:     nil,
		SE:     nil,
		SW:     nil,
		MinX:   minX,
		MaxX:   maxX,
		MinY:   minY,
		MaxY:   maxY,
		X:      helpers.MyRoundedMean(minX, maxX),
		Y:      helpers.MyRoundedMean(minY, maxY),
		Census: 0,
		P:      make(map[*T]bool),
	}
}

// Create a pointer of Node in a pointer location of this one was nil.
func getOrCreateQtNode[T Point](ptr **qtNode[T], minX int, maxX int, minY int, maxY int, parent *qtNode[T]) *qtNode[T] {
	if *ptr == nil {
		*ptr = newqtNode(minX, maxX, minY, maxY, parent)
	}
	return *ptr
}

func (q *QuadTree[T]) Add(point *T) error {
	if (*point).X() < q.WestPole || (*point).X() > q.EastPole || (*point).Y() < q.NorthPole || (*point).Y() > q.SouthPole {
		return errors.New("This point is out of bound.")
	}
	q.Census += 1
	getOrCreateQtNode(&q.root, q.WestPole, q.EastPole, q.NorthPole, q.SouthPole, nil).add(point, &q.pointDirectory)
	return nil
}

func getFirstPoint[T Point](m map[*T]bool) *T {
	for k := range m {
		return k
	}
	return nil
}

func (n *qtNode[T]) add(point *T, dir *map[*T]*qtNode[T]) {
	n.Census += 1
	if n.Census == 1 {
		(*dir)[point] = n
		n.P[point] = true
	} else {
		if len(n.P) > 0 {
			firstPoint := getFirstPoint(n.P)
			if (*firstPoint).X() == (*point).X() && (*firstPoint).Y() == (*point).Y() {
				(*dir)[point] = n
				n.P[point] = true
			} else {
				for p := range n.P {
					delete(*dir, p)
					delete(n.P, p)
					n.addInGoodNodeChild(p, dir)
				}
				n.addInGoodNodeChild(point, dir)
			}
		} else {
			n.addInGoodNodeChild(point, dir)
		}
	}
}

func (n *qtNode[T]) addInGoodNodeChild(point *T, dir *map[*T]*qtNode[T]) {
	if (*point).Y() <= n.Y {
		if (*point).X() <= n.X {
			getOrCreateQtNode(&(n.NW), n.MinX, n.X, n.MinY, n.Y, n).add(point, dir)
		} else {
			getOrCreateQtNode(&(n.NE), n.X + 1, n.MaxX, n.MinY, n.Y, n).add(point, dir)
		}
	} else {
		if (*point).X() <= n.X {
			getOrCreateQtNode(&(n.SW), n.MinX, n.X, n.Y + 1, n.MaxY, n).add(point, dir)
		} else {
			getOrCreateQtNode(&(n.SE), n.X + 1, n.MaxX, n.Y + 1, n.MaxY, n).add(point, dir)
		}
	}
}

func (n *qtNode[T]) fillGrid(grid *[][]byte, shiftX int, shiftY int, borderX int, borderY int) {
	if n != nil {
		if len(n.P) > 0 {
			p := getFirstPoint(n.P)
			(*grid)[(*p).Y() + shiftY][(*p).X() + shiftX] = helpers.LenToByte(len(n.P))
		} else {
			maxX := n.MaxX
			maxY := n.MaxY
			if (maxX == borderX) {
				maxX += 1
			}
			if (maxY == borderY) {
				maxY += 1
			}
			for y := n.MinY; y < maxY; y++ {
				(*grid)[y + shiftY][n.X + shiftX] = '|'
			}
			for x := n.MinX; x < maxX; x++ {
				(*grid)[n.Y + shiftY][x + shiftX] = '-'
			}
			n.NE.fillGrid(grid, shiftX, shiftY, borderX, borderY)
			n.NW.fillGrid(grid, shiftX, shiftY, borderX, borderY)
			n.SE.fillGrid(grid, shiftX, shiftY, borderX, borderY)
			n.SW.fillGrid(grid, shiftX, shiftY, borderX, borderY)
		}
	}
}

func (t *QuadTree[T]) ToGrid() [][]byte {
	grid := make([][]byte, t.Height)
	for y := range len(grid) {
		grid[y] = make([]byte, t.Width)
		for x := range len(grid[y]) {
			grid[y][x] = '.'
		}
	}

	shiftX := t.EastPole
	shiftY := t.SouthPole
	t.root.fillGrid(&grid, shiftX, shiftY, t.EastPole, t.SouthPole)

	return grid
}

func (t *QuadTree[T]) Print() {
	log.Println(t)
}

func (t *QuadTree[T]) String() string {
	g := t.ToGrid()
	len := t.Height * t.Width
	display := make([]byte, len + t.Width + 1)
	lnCount := 0

	for i := 0; i < t.Height * t.Width; i++ {
		if i % t.Width == 0 {
			display[i + lnCount] = '\n'
			lnCount++
		}
		display[i + lnCount] = g[i / t.Width][i % t.Width]
	}

	display[len + t.Width] = '\n'
	
	return string(display)
}