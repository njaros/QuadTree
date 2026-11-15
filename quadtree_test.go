package quadtree

import (
	"testing"
)

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
	t.Log("Instantiate QuadTree (-10, 10, -10, 10)")
	tree, err := NewQuadTree[TestPoint](-10, 10, -10, 10)
	if err != nil {
		t.Errorf("Should have no error, got %s", err.Error())
	}

	t.Log("Add 2 out of bounds points")

	err = tree.Add(&TestPoint{-11, 0})
	if err == nil {
		t.Errorf("An error should occurs here")
	}
	err = tree.Add(&TestPoint{1, 11})
	if err == nil {
		t.Errorf("An error should occurs here")
	}

	t.Log("Add point (-5, 5)")

	err = tree.Add(&TestPoint{-5, 5})

	if err != nil {
		t.Errorf("Should have no error, got %s", err.Error())
	}
	if tree.Census != 1 || tree.root.Census != 1 {
		t.Errorf("tree census and root census must be 1, have tree census %d and root census %d", tree.Census, tree.root.Census)
	}
	if tree.root.NE != nil || tree.root.NW != nil || tree.root.SE != nil || tree.root.SW != nil {
		t.Errorf("tree root with one element must not have any children.")
	}

	t.Log("Add point (5, 5)")

	err = tree.Add(&TestPoint{5, 5})
	if err != nil {
		t.Errorf("Should have no error, got %s", err.Error())
	}
	if tree.Census != 2 || tree.root.Census != 2 {
		t.Errorf("tree census and root census must be 2, have tree census %d and root census %d", tree.Census, tree.root.Census)
	}
	if !(tree.root.NE == nil && tree.root.NW == nil && tree.root.SE != nil && tree.root.SW != nil) {
		t.Errorf(
			"Only SE and SW root's children must be instantiate at this point, have NE: %p NW: %p SE: %p SW: %p",
			tree.root.NE,
			tree.root.NW,
			tree.root.SE,
			tree.root.SW,
		)
	}

	t.Log("Add point (-5, -5)")

	err = tree.Add(&TestPoint{-5, -5})
	if err != nil {
		t.Errorf("Should have no error, got %s", err.Error())
	}
	if tree.Census != 3 || tree.root.Census != 3 {
		t.Errorf("tree census and root census must be 3, have tree census %d and root census %d", tree.Census, tree.root.Census)
	}
	if !(tree.root.NE == nil && tree.root.NW != nil && tree.root.SE != nil && tree.root.SW != nil) {
		t.Errorf(
			"Only NW, SE and SW root's children must be instantiate at this point, have NE: %p NW: %p SE: %p SW: %p",
			tree.root.NE,
			tree.root.NW,
			tree.root.SE,
			tree.root.SW,
		)
	}

	t.Log("Add point (5, -5)")

	err = tree.Add(&TestPoint{5, -5})
	if err != nil {
		t.Errorf("Should have no error, got %s", err.Error())
	}
	if tree.Census != 4 || tree.root.Census != 4 {
		t.Errorf("tree census and root census must be 4, have tree census %d and root census %d", tree.Census, tree.root.Census)
	}
	if tree.root.NE == nil || tree.root.NW == nil || tree.root.SE == nil || tree.root.SW == nil {
		t.Errorf(
			"All root's children must have been instantiated at this point, have NE: %p NW: %p SE: %p SW: %p",
			tree.root.NE,
			tree.root.NW,
			tree.root.SE,
			tree.root.SW,
		)
	}

	t.Log("Add point (5, -5) again")

	err = tree.Add(&TestPoint{5, -5})
	if err != nil {
		t.Errorf("Should have no error, got %s", err.Error())
	}
	if tree.Census != 5 || tree.root.Census != 5 {
		t.Errorf("tree census and root census must be 4, have tree census %d and root census %d", tree.Census, tree.root.Census)
	}
	if tree.root.NE == nil || tree.root.NW == nil || tree.root.SE == nil || tree.root.SW == nil {
		t.Errorf(
			"All root's children must have been instantiated at this point, have NE: %p NW: %p SE: %p SW: %p",
			tree.root.NE,
			tree.root.NW,
			tree.root.SE,
			tree.root.SW,
		)
	}
	if tree.root.NE.Census != 2 {
		t.Errorf("tree.root.NE.Census expectet 2, have NE %d NW %d SE %d SW %d", tree.root.NE.Census, tree.root.NW.Census, tree.root.SE.Census, tree.root.SW.Census)
	}
	if len(tree.root.NE.P) != 2 {
		t.Errorf("tree.root.NE.P expected to contain 2 points, have %d", len(tree.root.NE.P))
	}

	t.Log("Add point (3, -3)")

	err = tree.Add(&TestPoint{3, -3})
	if err != nil {
		t.Errorf("Should have no error, got %s", err.Error())
	}
	if tree.Census != 6 || tree.root.Census != 6 {
		t.Errorf("tree census and root census must be 4, have tree census %d and root census %d", tree.Census, tree.root.Census)
	}
	if tree.root.NE == nil || tree.root.NW == nil || tree.root.SE == nil || tree.root.SW == nil {
		t.Errorf(
			"All root's children must have been instantiated at this point, have NE: %p NW: %p SE: %p SW: %p",
			tree.root.NE,
			tree.root.NW,
			tree.root.SE,
			tree.root.SW,
		)
	}
	nodeTarget := tree.root.NE
	if nodeTarget.Census != 3 {
		t.Errorf("tree.root.NE.Census expectet 3, have %d", nodeTarget.Census)
	}
	if len(nodeTarget.P) != 0 {
		t.Errorf("tree.root.NE.P expected to contain 0 points, have %d", len(nodeTarget.P))
	}
	if !(nodeTarget.SW != nil && nodeTarget.SE == nil && nodeTarget.NE == nil && nodeTarget.NW != nil) {
		t.Errorf(
			"Root's NE's children instantiated at this point must be NW and SW, have NE: %p NW: %p SE: %p SW: %p",
			nodeTarget.NE,
			nodeTarget.NW,
			nodeTarget.SE,
			nodeTarget.SW,
		)
	}

	t.Log("Add lot of points, then get tree stringified")

	tree.Add(&TestPoint{2, 2})
	tree.Add(&TestPoint{-1, 2})
	tree.Add(&TestPoint{4, -3})
	tree.Add(&TestPoint{8, 8})
	tree.Add(&TestPoint{-9, 9})
	tree.Add(&TestPoint{10, 10})
	tree.Add(&TestPoint{-10, -10})

	expected := "\n1.|..|....|....|.....\n..|..|....|....|.....\n-----|....|....|.....\n..|..|....|....|.....\n..|..|....|....|.....\n-----1----|----2-----\n.....|....|..|.|.....\n.....|....|..11|.....\n.....|....|----|.....\n.....|....|..|.|.....\n---------------------\n.....|....|..|.|.....\n.....|...1|.1|.|.....\n.....|....|----|.....\n.....|....|..|.|.....\n-----1----|----1-----\n.....|....|....|..|..\n.....|....|....|..|..\n.....|....|....|--1--\n.1...|....|....|..|..\n.....|....|....|..|.1\n"
	if expected != tree.String() {
		t.Errorf("expected tree: %s, have %s, len expected: %d len found: %d", expected, tree.String(), len(expected), len(tree.String()))
	}
}
