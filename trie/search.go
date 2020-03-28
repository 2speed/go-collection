package trie

import (
    "reflect"
    "sync"

    "github.com/2speed/go-collection"
)

type searchResult int

const (
    Extension = searchResult(1)
    Greater   = searchResult(2)
    Less      = searchResult(3)
    Matched   = searchResult(4)
    Prefix    = searchResult(5)
    Unmatched = searchResult(6)
)

const childNotFound = -1

var searchContextPool = sync.Pool{
    New: func() interface{} { return &searchContext{} },
}

func acquireSearchContext() *searchContext {
    return searchContextPool.Get().(*searchContext)
}

func releaseSearchContext(sctx *searchContext) {
    searchContextPool.Put(sctx)
}

type searchContext struct {
    pointer        Node
    digitizer      Digitizer
    branchPosition int
    numMatches     int
}

func (s *searchContext) atLeaf() bool {
    if _, ok := s.pointer.(LeafNode); ok {
        return true
    }

    return false
}

func (s *searchContext) atRoot() bool {
    return s.pointer.Parent() == nil
}

func (s *searchContext) moveToMaxDescendant() {
    for !s.atLeaf() {
        index := s.digitizer.Base() - 1
        for s.descendToIndex(index) == childNotFound {
            index--
        }
    }
}

func (s *searchContext) childBranchPosition(index int) int {
    return s.branchPosition + 1
}

func (s *searchContext) parentBranchPosition() int {
    return s.branchPosition - 1
}

func (s *searchContext) childIndexOf(element interface{}) int {
    return s.digitizer.DigitOf(element, s.branchPosition)
}

func (s *searchContext) descendTo(element interface{}) int {
    index := s.digitizer.DigitOf(element, s.branchPosition)

    return s.descendToIndex(index)
}

func (s *searchContext) descendToIndex(index int) int {
    child, err := s.pointer.ChildWithIndexOf(index)
    if err != nil || child == nil {
        return childNotFound
    }

    s.branchPosition = s.childBranchPosition(index)
    s.pointer = child

    return index
}

func (s *searchContext) ascend() int {
    s.branchPosition = s.parentBranchPosition()
    s.pointer = s.pointer.Parent()

    return s.branchPosition
}

func (s *searchContext) extendPath(element interface{}, node Node) int {
    index := s.digitizer.DigitOf(element, s.branchPosition)
    node.AddChildWithIndexOf(index, node)

    return s.descendToIndex(index)
}

func (s *searchContext) processedEndOfString(element interface{}) bool {
    childNode, _ := s.pointer.Parent().ChildWithIndexOf(0)

    return s.digitizer.IsPrefixFree() &&
        !s.pointer.IsRoot() &&
        reflect.DeepEqual(childNode, s.pointer)
}

func (s *searchContext) retraceToLastLeftFork(element interface{}) {
    for {
        if !s.atLeaf() {
            index := s.digitizer.DigitOf(element, s.branchPosition)
            for i := index - 1; i >= 0; i-- {
                if s.descendToIndex(i) != childNotFound {
                    return
                }
            }
        }

        if s.atRoot() {
            return
        } else {
            s.ascend()
        }
    }
}

func (s *searchContext) elementsInSubtree(collection collection.Collection) {
    if s.atLeaf() {
        collection.Add(s.pointer.Value())
    } else {
        for i := 0; i < s.digitizer.Base(); i++ {
            if s.descendToIndex(i) != childNotFound {
                s.elementsInSubtree(collection)
                s.ascend()
            }
        }
    }
}