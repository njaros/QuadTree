package quadtree

import "testing"

type TestPoint struct {
	x int
	y int
}

func (p TestPoint) X() int {
	return p.x
}

func (p TestPoint) Y() int {
	return p.y
}

// Some fixtures

func getTestTree(t *testing.T) *QuadTree[TestPoint] {
	tree, err := NewQuadTree[TestPoint](-1000, 1000, -1000, 1000)
	if err != nil {
		t.Errorf("Should have no error, got %s", err.Error())
	}

	tree.Add(&TestPoint{50, 50})
	tree.Add(&TestPoint{0, 0})
	tree.Add(&TestPoint{50, -50})
	tree.Add(&TestPoint{-50, 50})
	tree.Add(&TestPoint{-50, -50})
	tree.Add(&TestPoint{500, 150})
	tree.Add(&TestPoint{500, -50})
	tree.Add(&TestPoint{312, -7})
	tree.Add(&TestPoint{1, 2})
	tree.Add(&TestPoint{3, 4})
	tree.Add(&TestPoint{-1, -500})

	return tree
}

func TestNewQuadTree(t *testing.T) {
	tree, err := NewQuadTree[TestPoint](-101, 102, -103, 104)
	if err != nil {
		t.Errorf("Should have no error, got %s", err.Error())
	}
	if tree.WestPole != -101 {
		t.Errorf("West border must be -101, got %d", tree.WestPole)
	}
	if tree.EastPole != 102 {
		t.Errorf("East border must be 102, got %d", tree.EastPole)
	}
	if tree.NorthPole != -103 {
		t.Errorf("North border must be -103, got %d", tree.NorthPole)
	}
	if tree.SouthPole != 104 {
		t.Errorf("South border must be 104, got %d", tree.SouthPole)
	}

	tree, err = NewQuadTree[TestPoint](14, 5, 0, 100)
	if err == nil {
		t.Errorf("An error should occurs here")
	}
}

func TestQuadTreeAdd(t *testing.T) {
	tree, err := NewQuadTree[TestPoint](-10, 10, -10, 10)
	if err != nil {
		t.Errorf("Should have no error, got %s", err.Error())
	}

	err = tree.Add(&TestPoint{-11, 0})
	if err == nil {
		t.Errorf("An error should occurs here")
	}
	err = tree.Add(&TestPoint{1, 11})
	if err == nil {
		t.Errorf("An error should occurs here")
	}
	err = tree.Add(&TestPoint{-5, 5})
	if err != nil {
		t.Errorf("Should have no error, got %s", err.Error())
	}
}

func TestQuadTreeWithElements(t *testing.T) {
	tree := getTestTree(t)

	if tree.Census != 11 {
		t.Errorf("The census must be 11, found %d", tree.Census)
	}

	
}