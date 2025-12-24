package main

import (
	"sort"
)

// BVHNode is a node in a Bounding Volume Hierarchy tree.
// It implements the Hittable interface and accelerates ray-object intersection
// tests by organizing objects in a spatial tree structure.
//
// BVH reduces intersection tests from O(n) to O(log n) by allowing rays to
// skip entire subtrees that don't intersect with the bounding box.
type BVHNode struct {
	left  Hittable // Left child (can be another BVHNode or a primitive)
	right Hittable // Right child (can be another BVHNode, primitive, or nil)
	box   AABB     // Bounding box enclosing all objects in this subtree
}

// NewBVH constructs a BVH from a slice of Hittable objects.
// It recursively partitions the objects into a binary tree structure.
func NewBVH(objects []Hittable) *BVHNode {
	if len(objects) == 0 {
		panic("BVH: cannot build from empty object list")
	}
	return buildBVH(objects, 0, len(objects))
}

// buildBVH recursively constructs a BVH subtree for objects[start:end].
func buildBVH(objects []Hittable, start, end int) *BVHNode {
	node := &BVHNode{}
	objectCount := end - start

	// Base case: 1 object - create leaf node
	if objectCount == 1 {
		node.left = objects[start]
		node.right = nil
		node.box = getBoundingBox(objects[start])
		return node
	}

	// Base case: 2 objects - create node with two leaves
	if objectCount == 2 {
		node.left = objects[start]
		node.right = objects[start+1]
		box1 := getBoundingBox(objects[start])
		box2 := getBoundingBox(objects[start+1])
		node.box = box1.Union(box2)
		return node
	}

	// Recursive case: partition objects and build subtrees
	axis := chooseSplitAxis(objects, start, end)
	mid := partitionObjects(objects, start, end, axis)

	node.left = buildBVH(objects, start, mid)
	node.right = buildBVH(objects, mid, end)
	node.box = getBoundingBox(node.left).Union(getBoundingBox(node.right))

	return node
}

// Hit tests if a ray intersects any object in this BVH subtree.
// It first checks the bounding box for early rejection, then recursively
// tests children, using progressive tmax culling to find the nearest hit.
func (b *BVHNode) Hit(r Ray, tmin, tmax float64, hr *HitRecord) bool {
	// Early rejection: if ray doesn't hit bounding box, skip entire subtree
	if !b.box.Hit(r, tmin, tmax) {
		return false
	}

	// Test left child
	hitLeft := false
	closest := tmax
	if b.left != nil {
		hitLeft = b.left.Hit(r, tmin, closest, hr)
		if hitLeft {
			closest = hr.T // Progressive culling: narrow search range
		}
	}

	// Test right child with updated tmax
	hitRight := false
	if b.right != nil {
		hitRight = b.right.Hit(r, tmin, closest, hr)
	}

	return hitLeft || hitRight
}

// BoundingBox returns the bounding box for this BVH node.
func (b *BVHNode) BoundingBox() AABB {
	return b.box
}

// getBoundingBox extracts the bounding box from a Hittable.
// This uses type assertion to call the BoundingBox() method.
func getBoundingBox(h Hittable) AABB {
	type bounded interface {
		BoundingBox() AABB
	}

	if b, ok := h.(bounded); ok {
		return b.BoundingBox()
	}

	// Fallback: if object doesn't have BoundingBox method, panic
	// (all objects in our raytracer should have BoundingBox)
	panic("BVH: object does not implement BoundingBox()")
}

// chooseSplitAxis determines which axis (X=0, Y=1, Z=2) to split on.
// It chooses the axis with the largest extent across all objects.
func chooseSplitAxis(objects []Hittable, start, end int) int {
	// Compute bounding box of all objects in range
	bbox := getBoundingBox(objects[start])
	for i := start + 1; i < end; i++ {
		bbox = bbox.Union(getBoundingBox(objects[i]))
	}

	// Choose axis with largest extent
	extent := bbox.Max.Sub(bbox.Min)
	if extent.X > extent.Y && extent.X > extent.Z {
		return 0 // X axis
	} else if extent.Y > extent.Z {
		return 1 // Y axis
	}
	return 2 // Z axis
}

// partitionObjects sorts objects[start:end] along the given axis and returns
// the midpoint index. Objects are sorted by their bounding box centroids.
func partitionObjects(objects []Hittable, start, end, axis int) int {
	// Sort objects by centroid along chosen axis
	sort.Slice(objects[start:end], func(i, j int) bool {
		boxI := getBoundingBox(objects[start+i])
		boxJ := getBoundingBox(objects[start+j])

		// Compute centroids
		centerI := boxI.Min.Add(boxI.Max).DivS(2)
		centerJ := boxJ.Min.Add(boxJ.Max).DivS(2)

		// Compare along chosen axis
		switch axis {
		case 0:
			return centerI.X < centerJ.X
		case 1:
			return centerI.Y < centerJ.Y
		case 2:
			return centerI.Z < centerJ.Z
		default:
			panic("BVH: invalid axis")
		}
	})

	// Return midpoint index
	return start + (end-start)/2
}
