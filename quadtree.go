package quadtree

import "errors"

type Point interface {
	X() int
	Y() int
}

type QuadTree[T Point] struct {
	WestPole       int
	EastPole       int
	NorthPole      int
	SouthPole      int
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
		Census:         0,
		pointDirectory: make(map[*T]*qtNode[T]),
		root:           nil,
	}, nil
}

// qtNode may contains several points if they are superimposed (in same coordinates).
// Structure to store points in hashmap to ease to Point reference search.
type qtNode[T Point] struct {
	parent *qtNode[T]
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
	return &qtNode[T]{
		parent: parent,
		NW:     nil,
		NE:     nil,
		SE:     nil,
		SW:     nil,
		MinX:   minX,
		MaxX:   maxX,
		MinY:   minY,
		MaxY:   maxY,
		X:      (minX + maxX) / 2,
		Y:      (minY + maxX) / 2,
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
	getOrCreateQtNode(&q.root, q.EastPole, q.WestPole, q.NorthPole, q.SouthPole, nil).add(point, &q.pointDirectory)
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
	if (*point).X() > n.X {
		if (*point).Y() > n.Y {
			getOrCreateQtNode(&n.SE, n.X, n.MaxX, n.Y, n.MaxY, n).add(point, dir)
		} else {
			getOrCreateQtNode(&n.NE, n.X, n.MaxX, n.MinY, n.Y, n).add(point, dir)
		}
	} else {
		if (*point).Y() > n.Y {
			getOrCreateQtNode(&n.SW, n.MinX, n.X, n.Y, n.MaxY, n).add(point, dir)
		} else {
			getOrCreateQtNode(&n.NW, n.MinX, n.X, n.MinY, n.Y, n).add(point, dir)
		}
	}
}
