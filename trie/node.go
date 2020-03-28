package trie

import (
    "fmt"

    "github.com/pkg/errors"
)

// Node
type Node interface {
    SetParent(parent Node)
    Parent() Node
    AddChildWithIndexOf(index int, child Node) error
    ChildWithIndexOf(index int) (Node, error)
    RemoveChildWithIndexOf(index int) bool
    HasChildren() bool
    SetValue(element interface{})
    Value() interface{}
    IsRoot() bool
    IsLeaf() bool
}

type node struct {
    parent      Node
    children    []Node
    numChildren int
    element     interface{}
    isRoot      bool
}

func newNode(capacity int) Node {
    if capacity <= 0 {
        return &node{}
    }

    return &node{ children: make([]Node, capacity) }
}

func newRootNode(capacity int) Node {
    return &node{
        children: make([]Node, capacity),
        isRoot:   true,
    }
}

// SetParent
func (n *node) SetParent(parent Node) {
    n.parent = parent
}

// Parent
func (n *node) Parent() Node {
    return n.parent
}

// AddChildWithIndexOf
func (n *node) AddChildWithIndexOf(index int, child Node) error {
    if index < 0 || index >= len(n.children) {
        return errors.Errorf("index out of bounds [Node.capacity = %v, requested index = %v]", cap(n.children), index)
    }

    if n.children[index] != nil {
        return errors.Errorf("child exists at index %v", index)
    }

    if n.children[index] == nil {
        n.numChildren++
    }

    n.children[index] = child
    child.SetParent(n)

    return nil
}

// ChildWithIndexOf
func (n *node) ChildWithIndexOf(index int) (Node, error) {
    if err := n.checkBounds(index); err != nil {
        return nil, err
    }

    return n.children[index], nil
}

// RemoveChildWithIndexOf
func (n *node) RemoveChildWithIndexOf(index int) bool {
    if err := n.checkBounds(index); err != nil {
        return false
    }

    if n.children[index] != nil {
        n.children[index] = nil
        n.numChildren--

        return true
    }

    return false
}

// HasChildren
func (n *node) HasChildren() bool {
    return n.numChildren > 0
}

// SetValue
func (n *node) SetValue(element interface{}) {
    n.element = element
}

// Value
func (n *node) Value() interface{} {
    return n.element
}

// IsRoot
func (n *node) IsRoot() bool {
    return n.isRoot
}

// IsLeaf
func (n *node) IsLeaf() bool {
    return false
}

// String
func (n *node) String() string {
    return fmt.Sprintf("%v", n.element)
}

func (n *node) checkBounds(index int) error {
    if index < 0 || index > len(n.children) {
        return errors.Errorf("index out of bounds [Node.capacity = %v, requested index = %v]", cap(n.children), index)
    }

    return nil
}

// LeafNode
type LeafNode interface {
    Node

    AddAfter(leafNode LeafNode)
    SetNext(next LeafNode)
    Next() LeafNode
    SetPrevious(previous LeafNode)
    Previous() LeafNode
    IsDeleted() bool
    Remove()
    IsHead() bool
    IsTail() bool
}

type leafNode struct {
    Node

    next     LeafNode
    previous LeafNode
    isHead   bool
    isTail   bool
}

func newLeafNode() LeafNode {
    return &leafNode{ Node: newNode(0) }
}

func newHeadNode() LeafNode {
    return &leafNode{
        Node:   newNode(0),
        isHead: true,
    }
}

func newTailNode() LeafNode {
    return &leafNode{
        Node: newNode(0),
        isTail: true,
    }
}

// IsLeaf
func (l *leafNode) IsLeaf() bool {
    return true
}

// AddAfter
func (l *leafNode) AddAfter(leafNode LeafNode) {
    l.SetNext(leafNode.Next())
    leafNode.SetNext(l)
    l.SetPrevious(leafNode)
    l.next.SetPrevious(l)
}

// SetNext
func (l *leafNode) SetNext(next LeafNode) {
    l.next = next
}

// Next
func (l *leafNode) Next() LeafNode {
    return l.next
}

// SetPrevious
func (l *leafNode) SetPrevious(previous LeafNode) {
    l.previous = previous
}

// Previous
func (l *leafNode) Previous() LeafNode {
    return l.previous
}

// IsDeleted
func (l *leafNode) IsDeleted() bool {
    return l.previous == nil
}

// Remove
func (l *leafNode) Remove() {
    l.previous.SetNext(l.next)
    l.next.SetPrevious(l.previous)
    l.markDeleted()
}

// IsHead
func (l *leafNode) IsHead() bool {
    return l.isHead
}

// IsTail
func (l *leafNode) IsTail() bool {
    return l.isTail
}

func (l *leafNode) markDeleted() {
    l.previous = nil
}