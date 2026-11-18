package quadtree

import (
	"errors"
	"log"

	"github.com/njaros/QuadTree/helpers"
)

var (
	OutOfBoundsError = errors.New("This point is out of bound.")
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
	Width          int
	Height         int
	Census         int
	pointDirectory map[*T]*qtNode[T]
	root           *qtNode[T]
}

// qtNode may contains several points if they are superimposed (in same coordinates).
// Structure to store points in hashmap to ease to Point reference search.
type qtNode[T Point] struct {
	parent *qtNode[T]
	depth  int
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

// Constructors

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
		Height:         maxY - minY + 1,
		Width:          maxX - minX + 1,
		Census:         0,
		pointDirectory: make(map[*T]*qtNode[T]),
		root:           nil,
	}, nil
}

func newQtNode[T Point](minX int, maxX int, minY int, maxY int, parent *qtNode[T]) *qtNode[T] {
	depth := 0
	if parent != nil {
		depth = (*parent).depth + 1
	}
	// log.Printf("depth: %d\n", depth)
	node := &qtNode[T]{
		parent: parent,
		depth:  depth,
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

	return node
}

// Create a pointer of Node in a pointer location of this one was nil.
func getOrCreateQtNode[T Point](ptr **qtNode[T], minX int, maxX int, minY int, maxY int, parent *qtNode[T]) *qtNode[T] {
	if *ptr == nil {
		node := newQtNode(minX, maxX, minY, maxY, parent) 
		*ptr = node
		
		// if parent != nil {
		// 	var comeFrom string
		// 	if (*parent).NE == node {
		// 		comeFrom = "NE"
		// 	}
		// 	if (*parent).NW == node {
		// 		comeFrom = "NW"
		// 	}
		// 	if (*parent).SE == node {
		// 		comeFrom = "SE"
		// 	}
		// 	if (*parent).SW == node {
		// 		comeFrom = "SW"
		// 	}
		// 	log.Printf("I come from %s\n", comeFrom)
		// }
	}
	return *ptr
}

// Utilities

func (q *QuadTree[T]) InBounds(point *T) bool {
	return (*point).X() >= q.WestPole && (*point).X() <= q.EastPole && (*point).Y() >= q.NorthPole && (*point).Y() <= q.SouthPole
}

// Insert methods

// Insert a new point in the QuadTree. It must have X() int and Y() int method signatures.
// Duplicate coordinates is allowed
func (q *QuadTree[T]) Add(point *T) error {
	if !q.InBounds(point){
		return OutOfBoundsError
	}
	q.Census += 1
	getOrCreateQtNode(&q.root, q.WestPole, q.EastPole, q.NorthPole, q.SouthPole, nil).add(point, &q.pointDirectory)
	return nil
}

func (n *qtNode[T]) add(point *T, dir *map[*T]*qtNode[T]) {
	n.Census += 1
	if n.Census == 1 {
		(*dir)[point] = n
		n.P[point] = true
	} else {
		if len(n.P) > 0 {
			firstPoint, _, _ := helpers.GetOneOnMap(n.P)
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
			getOrCreateQtNode(&(n.NE), n.X+1, n.MaxX, n.MinY, n.Y, n).add(point, dir)
		}
	} else {
		if (*point).X() <= n.X {
			getOrCreateQtNode(&(n.SW), n.MinX, n.X, n.Y+1, n.MaxY, n).add(point, dir)
		} else {
			getOrCreateQtNode(&(n.SE), n.X+1, n.MaxX, n.Y+1, n.MaxY, n).add(point, dir)
		}
	}
}

// GroupAdd method is optimized for add a lot of point in one time. Experimental
func (q *QuadTree[T]) GroupAdd(pts []*T) error {
	// Sort point into
	pointMap := make(map[[2]int]*[]*T)
	for _, pt := range pts {
		if !q.InBounds(pt) {
			return OutOfBoundsError
		}
		elt, ok := pointMap[[2]int{(*pt).X(), (*pt).Y()}]
		if (ok) {
			*elt = append(*elt, pt)
		} else {
			pointMap[[2]int{(*pt).X(), (*pt).Y()}] = &[]*T{pt}
		}
	}
	q.Census += len(pts)

	getOrCreateQtNode(&q.root, q.WestPole, q.EastPole, q.NorthPole, q.SouthPole, nil).groupAdd(&pointMap, &q.pointDirectory)
	return nil
}

// pts is an address to avoid a long hmap copy
// Too long and too ugly method, with dangerous pointer juggling
func (n *qtNode[T]) groupAdd(pts *map[[2]int]*[]*T, dir *map[*T]*qtNode[T]) {
	empty := n.Census == 0
	if len(n.P) > 0 {
		inplacePointPtr, _, _ := helpers.GetOneOnMap(n.P)
		inplacePoint := *inplacePointPtr
		sameCoordPts, ok := (*pts)[[2]int{inplacePoint.X(), inplacePoint.Y()}]
		if ok == true && len(*pts) == 1 {
			n.Census += len(*sameCoordPts)
			for _, pt := range *sameCoordPts {
				n.P[pt] = true
				(*dir)[pt] = n
			}
		} else {
			if !ok {
				// Dangerous memory juggling ?
				sameCoordPts = &[]*T{}
				(*pts)[[2]int{inplacePoint.X(), inplacePoint.Y()}] = sameCoordPts
			}
			for pt := range(n.P) {
				delete(n.P, pt)
				delete(*dir, pt)
				*sameCoordPts = append(*sameCoordPts, pt)
			}
			n.groupChildrenSplit(pts, dir)
		}
	} else if empty && len(*pts) == 1 {
		_, poorable, _ := helpers.GetOneOnMap(*pts)
		for _, pt := range *poorable {
			n.P[pt] = true
			(*dir)[pt] = n
		}
		n.Census += len(*poorable)
	} else {
		n.groupChildrenSplit(pts, dir)
	}
}

func (n *qtNode[T]) groupChildrenSplit(ptsPtr *map[[2]int]*[]*T, dir *map[*T]*qtNode[T]) {
	groupNE := map[[2]int]*[]*T{}
	groupNW := map[[2]int]*[]*T{}
	groupSE := map[[2]int]*[]*T{}
	groupSW := map[[2]int]*[]*T{}

	for coord, pts := range *ptsPtr {
		if coord[1] <= n.Y {
			if coord[0] <= n.X {
				groupNW[coord] = pts
			} else {
				groupNE[coord] = pts
			}
		} else {
			if coord[0] <= n.X {
				groupSW[coord] = pts
			} else {
				groupSE[coord] = pts
			}
		}
		n.Census += len(*pts)
	}
	if len(groupNW) > 2 {
		log.Printf("NW group len: %d\n", len(groupNW))
	} else {
		array := [][2]int{}
		for k := range groupNW {
			array = append(array, k)
		}
		log.Printf("%v\n", array)
	}
	if len(groupNE) > 0 {
		getOrCreateQtNode(&(n.NE), n.X+1, n.MaxX, n.MinY, n.Y, n).groupAdd(&groupNE, dir)
	}
	if len(groupNW) > 0 {
		getOrCreateQtNode(&(n.NW), n.X+1, n.MaxX, n.MinY, n.Y, n).groupAdd(&groupNW, dir)
	}
	if len(groupSE) > 0 {
		getOrCreateQtNode(&(n.SE), n.X+1, n.MaxX, n.MinY, n.Y, n).groupAdd(&groupSE, dir)
	}
	if len(groupSW) > 0 {
		getOrCreateQtNode(&(n.SW), n.X+1, n.MaxX, n.MinY, n.Y, n).groupAdd(&groupSW, dir)
	}
}

// Displaying part

func (n *qtNode[T]) fillGrid(grid *[][]byte, shiftX int, shiftY int, borderX int, borderY int) {
	if n != nil {
		if p, _, ok := helpers.GetOneOnMap(n.P); ok {
			(*grid)[(*p).Y()+shiftY][(*p).X()+shiftX] = helpers.LenToByte(len(n.P))
		} else {
			maxX := n.MaxX
			maxY := n.MaxY
			if maxX == borderX {
				maxX += 1
			}
			if maxY == borderY {
				maxY += 1
			}
			for y := n.MinY; y < maxY; y++ {
				(*grid)[y+shiftY][n.X+shiftX] = '|'
			}
			for x := n.MinX; x < maxX; x++ {
				(*grid)[n.Y+shiftY][x+shiftX] = '-'
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
	display := make([]byte, len+t.Width+1)
	lnCount := 0

	for i := 0; i < t.Height*t.Width; i++ {
		if i%t.Width == 0 {
			display[i+lnCount] = '\n'
			lnCount++
		}
		display[i+lnCount] = g[i/t.Width][i%t.Width]
	}

	display[len+t.Width] = '\n'

	return string(display)
}
