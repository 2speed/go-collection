package trie

import (
    "fmt"
    "reflect"
    "testing"

    "github.com/2speed/go-collection"
    "github.com/2speed/go-collection/list"
)

func TestTrie_Add(t *testing.T) {
    trie := NewTrie(4)

    t.Run("ab", func(t *testing.T) {
        value := "ab"
        err   := trie.Add(value)

        assertError(t, err, nil)
        assertSize(t, trie,1)
        assertContains(t, trie, value, true)
        assertContains(t, trie, "abc", false)
        assertContains(t, trie, "a", false)
        assertContains(t, trie, "acb", false)
    })

    t.Run("abcd", func(t *testing.T) {
        value := "abcd"
        err   := trie.Add(value)

        assertError(t, err, nil)
        assertSize(t, trie,2)
        assertContains(t, trie, value, true)
    })

    t.Run("acb", func(t *testing.T) {
        value := "acb"
        err   := trie.Add(value)

        assertError(t, err, nil)
        assertSize(t, trie,3)
        assertContains(t, trie, value, true)
    })

    t.Run("cbca", func(t *testing.T) {
        value := "cbca"
        err   := trie.Add(value)

        assertError(t, err, nil)
        assertSize(t, trie,4)
        assertContains(t, trie, value, true)
    })
}

func TestTrie_AddAll(t *testing.T) {
    trie   := NewTrie(26)
    values := []interface{}{ "the", "quick", "brown", "fox" }
    err    := trie.AddAll(list.NewArrayListOf(values))

    assertError(t, err, nil)
    assertSize(t, trie,4)
    assertContains(t, trie, "the", true)
    assertContains(t, trie, "quick", true)
    assertContains(t, trie, "brown", true)
    assertContains(t, trie, "fox", true)
    assertContentEquals(t, trie, "[brown, fox, quick, the]")
}

func TestTrie_Remove(t *testing.T) {
    trie   := NewTrie(26)
    values := []interface{}{ "jumped", "over", "the", "lazy", "dog" }
    err    := trie.AddAll(list.NewArrayListOf(values))

    assertError(t, err, nil)
    assertSize(t, trie,5)
    assertContains(t, trie, "jumped", true)
    assertContains(t, trie, "over", true)
    assertContains(t, trie, "the", true)
    assertContains(t, trie, "lazy", true)
    assertContains(t, trie, "dog", true)
    assertContentEquals(t, trie, "[dog, jumped, lazy, over, the]")

    trie.Remove("lazy")
    trie.Remove("the")
    trie.Remove("fox")
    assertSize(t, trie,3)

    assertContains(t, trie, "lazy", false)
    assertContains(t, trie, "the", false)
    assertContentEquals(t, trie, "[dog, jumped, over]")

    trie.Clear()
    assertSize(t, trie,0)
}

func TestTrie_MinMax(t *testing.T) {
    trie   := NewTrie(5)
    values := []interface{}{ "cba", "ab", "bce", "abcd" }
    err    := trie.AddAll(list.NewArrayListOf(values))

    assertError(t, err, nil)
    assertSize(t, trie,4)
    assertContentEquals(t, trie, "[ab, abcd, bce, cba]")
    assertNodeValue(t, trie.Min(), "ab")
    assertNodeValue(t, trie.Max(), "cba")
}

func TestTrie_Predecessor(t *testing.T) {
    trie   := NewTrie(4)
    values := []interface{}{ "bac", "dab", "dabb", "dac", "daca", "dabba", "ab" }
    err    := trie.AddAll(list.NewArrayListOf(values))

    assertError(t, err, nil)
    assertSize(t, trie,7)
    assertContentEquals(t, trie, "[ab, bac, dab, dabb, dabba, dac, daca]")
    assertNodeValue(t, trie.Predecessor("dabba"), "dabb")
    assertNodeValue(t, trie.Predecessor("bac"), "ab")
}

func TestTrie_Successor(t *testing.T) {
    trie   := NewTrie(4)
    values := []interface{}{ "bac", "dab", "dabb", "dac", "daca", "dabba", "ab" }
    err    := trie.AddAll(list.NewArrayListOf(values))

    assertError(t, err, nil)
    assertSize(t, trie,7)
    assertContentEquals(t, trie, "[ab, bac, dab, dabb, dabba, dac, daca]")
    assertNodeValue(t, trie.Successor("dabba"), "dac")
    assertNodeValue(t, trie.Successor("bac"), "dab")
}

func TestTrie_Completions(t *testing.T) {
    trie   := NewTrie(4)
    values := []interface{}{ "acb", "dabc", "daca", "da", "ab" }
    err    := trie.AddAll(list.NewArrayListOf(values))

    assertError(t, err, nil)
    assertSize(t, trie,5)
    assertContentEquals(t, trie, "[ab, acb, da, dabc, daca]")

    l := list.NewArrayList()
    trie.Completions("a", l)
    assertContentEquals(t, l, "[ab, acb]")

    l.Clear()
    trie.Completions("da", l)
    assertContentEquals(t, l, "[da, dabc, daca]")
}

func TestTrie_LongestCommonPrefix(t *testing.T) {
    trie   := NewTrie(4)
    values := []interface{}{ "acb", "dadc", "dada", "da", "ab" }
    err    := trie.AddAll(list.NewArrayListOf(values))

    assertError(t, err, nil)
    assertSize(t, trie,5)
    assertContentEquals(t, trie, "[ab, acb, da, dada, dadc]")

    l := list.NewArrayList()
    trie.LongestCommonPrefix("a", l)
    assertContentEquals(t, l, "[ab, acb]")

    l.Clear()
    trie.LongestCommonPrefix("dadda", l)
    assertContentEquals(t, l, "[dada, dadc]")
}

func assertError(t *testing.T, actual error, expected error) {
    t.Helper()

    if actual != expected {
        t.Errorf("expected error '%s', but found '%s'", expected, actual)
    }
}

func assertContains(t *testing.T, collection collection.Collection, value string, expected bool) {
    t.Helper()

    if collection.Contains(value) != expected {
        if expected {
            t.Errorf("expected to contain value: %s", value)
        } else {
            t.Errorf("expected not to contain value: %s", value)
        }
    }
}

func assertNodeValue(t *testing.T, actual interface{}, expected string) {
    t.Helper()

    if n, ok := actual.(string); ok {
        if n != expected {
            t.Errorf("expected content of '%s', but found '%s'", expected, n)
        }
    } else {
        t.Errorf("expected type of 'string', but found '%v'", reflect.TypeOf(actual))
    }
}

func assertSize(t *testing.T, collection collection.Collection, expected int) {
    t.Helper()

    actual := collection.Size()
    if actual != expected {
        t.Errorf("expected size of '%d', but found '%d'", expected, actual)
    }
}

func assertContentEquals(t *testing.T, collection collection.Collection, expected string) {
    t.Helper()

    actual := fmt.Sprintf("%s", collection)
    if actual != expected {
        t.Errorf("expected content of '%s', but found '%s'", expected, actual)
    }
}