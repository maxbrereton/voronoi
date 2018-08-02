// Note - Currently this is just a standard Binary Search Tree Implementation
// The plan is to eventually move to a redblack tree
package main

import (
	"container/heap"
	"fmt"
)

const (
	red   bool = true
	black bool = false
)

// Internal Nodes - These will have a breakpoint but arcSite will be nil
// Leaf Nodes - These have an arcSite but no breakpoint
type node struct {
	left, right, parent, next, previous *node
	colour                              bool
	breakpoint                          *breakpoint
	arcSite                             *site
	key                                 int
	circleEvent                         *Item
	edge                                *halfEdge
}

type redblacktree struct {
	root *node
}

func (rbtree *redblacktree) insert(newKey int, newSite *site, eventQueue *PriorityQueue) {
	if rbtree.root == nil {
		rbtree.root = &node{key: newKey, colour: black, arcSite: newSite}
	} else {
		rbtree.root = rbtree.root.insert(rbtree.root, newKey, newSite, eventQueue)
	}
}

func (n *node) insert(currentNode *node, newKey int, newSite *site, eventQueue *PriorityQueue) *node {
	// Check if this is a leaf node
	if currentNode.breakpoint == nil {

		if currentNode.circleEvent != nil {
			// Remove circle event from event queue as it is a false alarm
			heap.Remove(eventQueue, currentNode.circleEvent.index)
		}

		// Define the breakpoints that will be used in the two new internal nodes
		leftBreakpoint := breakpoint{leftSite: currentNode.arcSite, rightSite: newSite}
		rightBreakpoint := breakpoint{leftSite: newSite, rightSite: currentNode.arcSite}

		// The 3 leaf nodes that represent the arcs
		leftLeafNode := node{arcSite: currentNode.arcSite, previous: currentNode.previous, key: currentNode.key}
		middleLeafNode := node{arcSite: newSite, previous: &leftLeafNode, key: newKey}
		rightLeafNode := node{arcSite: currentNode.arcSite, next: currentNode.next, previous: &middleLeafNode,
			key: currentNode.key}
		middleLeafNode.next = &rightLeafNode
		leftLeafNode.next = &middleLeafNode

		leftInternalNode := node{left: &leftLeafNode, right: &middleLeafNode, breakpoint: &leftBreakpoint}
		rightInternalNode := node{left: &leftInternalNode, right: &rightLeafNode, breakpoint: &rightBreakpoint}

		// Set parent nodes
		leftInternalNode.parent = &rightInternalNode
		rightLeafNode.parent = &rightInternalNode
		middleLeafNode.parent = &leftInternalNode
		leftLeafNode.parent = &leftInternalNode

		// Check for circle event (i.e. check for unique triples of sites on beachline (a,b,c))
		leftLeafNode.circleEvent = checkCircleEvent(&leftLeafNode, newSite.y, eventQueue)
		rightLeafNode.circleEvent = checkCircleEvent(&rightLeafNode, newSite.y, eventQueue)

		return &rightInternalNode
	}

	// The directrix will be at the same y coordinate as the new site being added
	breakpointXCoordinate := getBreakpointXCoordinate(currentNode.breakpoint, newSite.y)

	fmt.Println("Breakpoint x-coord --> ", breakpointXCoordinate)
	fmt.Println("New site being added x-coord --> ", newSite.x)

	if newSite.x < breakpointXCoordinate {
		currentNode.left = currentNode.insert(currentNode.left, newKey, newSite, eventQueue)
		currentNode.left.parent = currentNode
	} else if newSite.x > breakpointXCoordinate {
		currentNode.right = currentNode.insert(currentNode.right, newKey, newSite, eventQueue)
		currentNode.right.parent = currentNode
	}

	return currentNode
}

func (rbtree *redblacktree) removeArc(leafNode *node, eventQueue *PriorityQueue) {
	if leafNode.parent == nil {
		// This happens if only 1 node is in the tree and you call remove - this should never happen
		rbtree.root = nil
	}

	leftLeafNode := leafNode.previous
	rightLeafNode := leafNode.next

	// Check if left or right have circle events - these won't be valid once leaf node has been removed
	if leftLeafNode.circleEvent != nil {
		heap.Remove(eventQueue, leftLeafNode.circleEvent.index)
	}
	if rightLeafNode.circleEvent != nil {
		heap.Remove(eventQueue, rightLeafNode.circleEvent.index)
	}
}

func (rbtree *redblacktree) inorderTraversal() {
	rbtree.root.inorderTraversal(rbtree.root)
}

func (n *node) inorderTraversal(currentNode *node) {
	if currentNode != nil {
		currentNode.inorderTraversal(currentNode.left)
		//fmt.Printf("%d ", currentNode.key)
		fmt.Println(currentNode)
		currentNode.inorderTraversal(currentNode.right)
	}
}
