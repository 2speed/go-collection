package trie

import (
    "math"
    "reflect"
)

type radixTree struct {
    *trie
}

// NewRadixTree
func NewRadixTree(capacity int) Trie {
    return NewRadixTreeWithDigitizer(NewStringDigitizer(capacity))
}

// NewRadixTreeWithDigitizer
func NewRadixTreeWithDigitizer(digitizer Digitizer) Trie {
    return &radixTree{ trie: newTrieWithDigitizer(digitizer) }
}

func (rt *radixTree) find(element interface{}, sctx *searchContext) searchResult {
    rt.prepareSearch(sctx)

    if rt.IsEmpty() {
        return Unmatched
    }

    numDigitsInElement := rt.digitizer.NumDigitsOf(element)

    for sctx.branchPosition < numDigitsInElement && !sctx.atLeaf() {
        if sctx.descendTo(element) == childNotFound {
            return Unmatched
        }
    }

    if sctx.branchPosition == numDigitsInElement {
        if sctx.atLeaf() {
            return Matched
        } else {
            return Prefix
        }
    } else {
        sctx.numMatches = sctx.branchPosition
        return rt.checkMatchFromLeaf(element, sctx)
    }
}

func (rt *radixTree) addNode(node Node, searchContext *searchContext) {
    if searchContext.pointer == nil {
        rt.root = newRootNode(rt.capacity)
        searchContext.pointer = rt.root
    }

    element := node.Value()

    if searchContext.atLeaf() {
        leafNode := searchContext.pointer

        searchContext.ascend()

        for searchContext.childIndexOf(element) == searchContext.childIndexOf(leafNode.Value()) {
            searchContext.extendPath(element, newNode(rt.capacity))
        }

        searchContext.extendPath(leafNode.Value(), leafNode)
        searchContext.ascend()
    }

    searchContext.extendPath(element, node)
}

func (rt *radixTree) remove(node Node) {
    // TODO: implement
}

func (rt *radixTree) checkMatchFromLeaf(element interface{}, sctx *searchContext) searchResult {
    leafData := sctx.pointer.Value()
    stop     := int(math.Min(float64(rt.digitizer.NumDigitsOf(element)), float64(rt.digitizer.NumDigitsOf(leafData))))

    for sctx.numMatches < stop {
        targetDigit := rt.digitizer.DigitOf(element, sctx.numMatches)
        leafDigit   := rt.digitizer.DigitOf(leafData, sctx.numMatches)
        comparison  := targetDigit - leafDigit

        if comparison < 0 {
            if rt.digitizer.IsPrefixFree() && targetDigit == 0 {
                return Prefix
            } else {
                return Less
            }
        } else if comparison > 0 {
            return Greater
        } else {
            sctx.numMatches++
        }
    }

    if sctx.numMatches == rt.digitizer.NumDigitsOf(element) {
        if sctx.numMatches == rt.digitizer.NumDigitsOf(leafData) {
            return Matched
        } else {
            return Prefix
        }
    } else {
        return Extension
    }
}

func (rt *radixTree) childIndexOf(parent Node, child Node) int {
    var index int
    c, _ := parent.ChildWithIndexOf(index)
    for !reflect.DeepEqual(c, child) && index < rt.capacity {
        c, _ = parent.ChildWithIndexOf(index)
        index++
    }

    return index
}
