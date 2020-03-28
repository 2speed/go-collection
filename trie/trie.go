package trie

import (
    "fmt"
    "strings"

    "github.com/2speed/go-collection"
    "github.com/pkg/errors"
)

// Trie
type Trie interface {
    collection.Ordered

    // Completions finds all elements in the Trie that match the provided prefix, and appends the matching elements
    // (if any) to the provided collection.
    Completions(prefix interface{}, collection collection.Collection)


    // LongestCommonPrefix finds all elements in the Trie that share the longest common prefix with the provided
    // element, and appends the matching elements (if any) to the provided collection.
    LongestCommonPrefix(element interface{}, collection collection.Collection)
}

type trie struct {
    root      Node
    head      LeafNode
    tail      LeafNode
    digitizer Digitizer
    capacity  int
    base      int
    size      int
}

func newTrie(capacity int) *trie {
    return newTrieWithDigitizer(NewStringDigitizer(capacity))
}

func newTrieWithDigitizer(digitizer Digitizer) *trie {
    capacity := digitizer.Base()
    head     := newHeadNode()
    tail     := newTailNode()

    head.SetNext(tail)
    tail.SetNext(head)

    return &trie{
        head:      head,
        tail:      tail,
        digitizer: digitizer,
        capacity:  capacity,
    }
}

// NewTrie creates a new Trie with the provided capacity. The capacity is used to set the base (or range of digits) used
// by the StringDigitizer for the trie. The base of the StringDigitizer is used to set the number of children each Node
// within the trie has.
func NewTrie(capacity int) Trie {
    return newTrie(capacity)
}

// NewTrieWithDigitizer
func NewTrieWithDigitizer(digitizer Digitizer) Trie {
    return newTrieWithDigitizer(digitizer)
}

// Add inserts the provided element into the Trie. The returned error will be non-nil for bounded Collection
// implementations that have reached capacity and cannot hold any further elements.
func (t *trie) Add(element interface{}) error {
    _, err := t.insert(element)

    return err
}

// AddAll inserts all elements from the provided collection into the Trie. The returned error will be non-nil
// for bounded Collection implementations that have reached capacity and cannot hold any further elements.
func (t *trie) AddAll(collection collection.Collection) error {
    if collection != nil {
        for _, v := range collection.Values() {
            if err := t.Add(v); err != nil {
                return err
            }
        }
    }

    return nil
}

// ValueWithIndex returns the element at the position specified by the provided index. The returned error will be
// non-nil if the provided index is outside the current bounds of the trie (index < 0 || index > trie.Size() - 1).
func (t *trie) ValueWithIndex(index int) (interface{}, error) {
    if err := t.checkBounds(index); err != nil {
        return nil, err
    }

    iterator := newIterator(t, t.head)
    for i := 0; i < index + 1; i++ {
        iterator.advance()
    }

    return iterator.get(), nil
}

// Remove removes the first occurrence (if any) of an element equivalent to the provided element. If an element was
// removed, the return value will be true, otherwise false will be returned.
func (t *trie) Remove(element interface{}) bool {
    if t.IsEmpty() {
        return false
    }

    sctx := acquireSearchContext()
    defer releaseSearchContext(sctx)

    if t.find(element, sctx) != Matched {
        return false
    }

    t.remove(sctx.pointer)

    return true
}

// Min returns the element with the lowest position in the Trie. More specifically, the first element in the iteration
// order is returned.
func (t *trie) Min() interface{} {
    if !t.IsEmpty() {
        return t.head.Next().Value()
    }

    return nil
}

// Max returns the element with the highest position in the Trie. More specifically, the last element in the iteration
// order is returned.
func (t *trie) Max() interface{} {
    if !t.IsEmpty() {
        return t.tail.Previous().Value()
    }

    return nil
}

// Predecessor returns the element (if any) from the Trie that is less than the provided element. More specifically, the
// element before the first occurrence of the provided element in iteration order is returned.
func (t *trie) Predecessor(element interface{}) interface{} {
    if !t.IsEmpty() {
        sctx := acquireSearchContext()
        defer releaseSearchContext(sctx)

        if t.moveToPredecessor(element, sctx, t.find(element, sctx)) {
            return sctx.pointer.Value()
        }
    }

    return nil
}

// Successor returns the element (if any) from the Trie that is greater than the provided element. More specifically,
// the element after the first occurrence of the provided element in iteration order is returned.
func (t *trie) Successor(element interface{}) interface{} {
    if !t.IsEmpty() {
        sctx := acquireSearchContext()
        defer releaseSearchContext(sctx)

        searchResult := t.find(element, sctx)
        successor    := t.tail
        if searchResult == Matched {
            successor = sctx.pointer.(LeafNode).Next()
        } else if t.moveToPredecessor(element, sctx, t.find(element, sctx)) {
            successor = sctx.pointer.(LeafNode).Next()
        }

        if !successor.IsTail() {
            return successor.Value()
        }
    }

    return nil
}

// Completions finds all elements in the trie that match the provided prefix, and appends the matching elements (if any)
// to the provided collection.
func (t *trie) Completions(prefix interface{}, collection collection.Collection) {
    if !t.IsEmpty() {
        sctx := acquireSearchContext()
        defer releaseSearchContext(sctx)

        searchResult := t.find(prefix, sctx)
        numDigits    := t.digitizer.NumDigitsOf(prefix)
        if t.digitizer.IsPrefixFree() {
            numDigits--
            if sctx.processedEndOfString(prefix) {
                sctx.ascend()
            }
        }

        if searchResult == Prefix || searchResult == Matched || sctx.branchPosition == numDigits {
            sctx.elementsInSubtree(collection)
        }
    }

    return
}

// LongestCommonPrefix finds all elements in the trie that share the longest common prefix with the provided element,
// and appends the matching elements (if any) to the provided collection.
func (t *trie) LongestCommonPrefix(prefix interface{}, collection collection.Collection) {
    if !t.IsEmpty() {
        sctx := acquireSearchContext()
        defer releaseSearchContext(sctx)

        t.find(prefix, sctx)
        if sctx.processedEndOfString(prefix) {
            sctx.ascend()
        }
        sctx.elementsInSubtree(collection)
    }

    return
}

// Size returns the number of elements in the Trie.
func (t *trie) Size() int {
    return t.size
}

// IsEmpty returns true if the Trie contains no elements, otherwise false is returned.
func (t *trie) IsEmpty() bool {
    return t.Size() == 0
}

// Clear removes all elements from the Trie.
func (t *trie) Clear() {
    iterator := newIterator(t, t.head)
    for iterator.advance() {
        iterator.remove()
    }
}

// Contains returns true if an element equivalent to the provided element exists in the Trie, otherwise false is
// returned.
func (t *trie) Contains(element interface{}) bool {
    if t.IsEmpty() {
        return false
    }

    sctx := acquireSearchContext()
    defer releaseSearchContext(sctx)

    return t.find(element, sctx) == Matched
}

// Values returns a slice containing the elements in the trie in the iteration order.
func (t *trie) Values() []interface{} {
    elements := make([]interface{}, 0)
    iterator := newIterator(t, t.head)
    for iterator.advance() {
        elements = append(elements, iterator.get())
    }

    return elements
}

// String returns a string representation of the Trie in it's current state.
func (t *trie) String() string {
    if t.Size() == 0 {
        return "[]"
    }

    elements := make([]string, 0, t.Size())
    locator  := newIterator(t, t.head)
    for locator.advance() {
        elements = append(elements, fmt.Sprintf("%v", locator.get()))
    }

    return "[" + strings.Join(elements, ", ") + "]"
}

func (t *trie) checkBounds(index int) error {
    if index < 0 || index >= t.Size() {
        return errors.Errorf("index out of bounds [no elements exist for requested index = %v]", index)
    }

    return nil
}

func (t *trie) find(element interface{}, sctx *searchContext) searchResult {
    t.prepareSearch(sctx)

    if t.IsEmpty() {
        return Unmatched
    }

    numDigitsInElement := t.digitizer.NumDigitsOf(element)

    for sctx.pointer != nil && !sctx.atLeaf() {
        switch {
        case sctx.branchPosition == numDigitsInElement:
            return Prefix
        case sctx.descendTo(element) == childNotFound:
            return Unmatched
        }
    }

    if sctx.pointer != nil && sctx.branchPosition != numDigitsInElement {
        return Extension
    } else {
        return Matched
    }
}

func (t *trie) prepareSearch(sctx *searchContext) {
    sctx.pointer        = t.root
    sctx.digitizer      = t.digitizer
    sctx.branchPosition = 0
}

func (t *trie) insert(element interface{}) (Node, error) {
    sctx := acquireSearchContext()
    defer releaseSearchContext(sctx)

    searchResult := t.find(element, sctx)
    if searchResult == Matched || (!t.digitizer.IsPrefixFree() && (searchResult == Prefix || searchResult == Extension)) {
        return nil, errors.New(fmt.Sprintf( "element violates prefix-free requirement: %v", element))
    }

    leafNode := newLeafNode()
    leafNode.SetValue(element)
    t.addNode(leafNode, sctx)
    searchResult = Matched

    if t.moveToPredecessor(element, sctx, searchResult) {
        leafNode.AddAfter(sctx.pointer.(LeafNode))
    } else {
        leafNode.AddAfter(t.head)
    }

    t.size++

    return leafNode, nil
}

func (t *trie) addNode(node Node, sctx *searchContext) {
    if sctx.pointer == nil {
        t.root = newRootNode(t.capacity)
        sctx.pointer = t.root
    }

    element := node.Value()

    for sctx.branchPosition < t.digitizer.NumDigitsOf(element) - 1 {
        index     := t.digitizer.DigitOf(element, sctx.branchPosition)
        childNode := newNode(t.capacity)
        sctx.pointer.AddChildWithIndexOf(index, childNode)
        sctx.pointer = childNode
        sctx.branchPosition++
    }

    index := t.digitizer.DigitOf(element, sctx.branchPosition)
    sctx.pointer.AddChildWithIndexOf(index, node)
    sctx.pointer = node
    sctx.branchPosition++
}

func (t *trie) remove(node Node) {
    if leafNode, ok := node.(LeafNode); ok {
        leafNode.Remove()
    }

    element := node.Value()
    level   := t.digitizer.NumDigitsOf(element)

    for !node.IsRoot() && !node.HasChildren() {
        parent := node.Parent()
        level--
        parent.RemoveChildWithIndexOf(t.digitizer.DigitOf(element, level))
        node = parent
    }

    t.size--
}

func (t *trie) moveToPredecessor(element interface{}, sctx *searchContext, searchResult searchResult) bool {
    if sctx.atLeaf() && (searchResult == Greater || searchResult == Extension) {
        return true
    }

    if searchResult != Greater {
        sctx.retraceToLastLeftFork(element)
    }

    if sctx.atRoot() {
        return false
    } else if !sctx.atLeaf() {
        sctx.moveToMaxDescendant()
    }

    return true
}

type iterator struct {
    trie    *trie
    pointer LeafNode
}

func newIterator(trie *trie, pointer LeafNode) *iterator {
    return &iterator{ trie: trie, pointer: pointer }
}

func (i *iterator) inCollection() bool {
    if i.pointer.IsHead() || i.pointer.IsTail() {
        return false
    }

    return !i.pointer.IsDeleted()
}

func (i *iterator) get() interface{} {
    if i.inCollection() {
        return i.pointer.Value()
    }

    return nil
}

func (i *iterator) skipRemovedElements(leafNode LeafNode) LeafNode {
    if leafNode.IsHead() || leafNode.IsTail() || !leafNode.IsDeleted() {
        return leafNode
    }

    leafNode.SetNext(i.skipRemovedElements(leafNode.Next()))

    return leafNode.Next()
}

func (i *iterator) advance() bool {
    if i.pointer.IsTail() {
        return false
    }

    if !i.pointer.IsHead() && i.pointer.IsDeleted() {
        i.pointer = i.skipRemovedElements(i.pointer)
    } else {
        i.pointer = i.pointer.Next()
    }

    return !i.pointer.IsTail()
}

func (i *iterator) retreat() bool {
    if !i.pointer.IsTail() && !i.pointer.IsHead() && i.pointer.IsDeleted() {
        i.pointer = i.skipRemovedElements(i.pointer)
    }

    i.pointer = i.pointer.Previous()

    return !i.pointer.IsHead()
}

func (i *iterator) hasNext() bool {
    if i.pointer.IsDeleted() {
        i.skipRemovedElements(i.pointer)
    }

    return !i.pointer.IsTail() && !i.pointer.Next().IsTail()
}

func (i *iterator) remove() {
    if i.inCollection() {
        i.trie.remove(i.pointer)
    }
}