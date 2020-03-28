package list

import (
    "strings"
    "testing"

    "github.com/2speed/go-collection"
)

type element struct {
    value    string
    position int
}

func TestArrayList_Add(t *testing.T) {
    elements := []interface{}{
        element{ value: "piranha plant", position: 0 },
        element{ value: "samus",         position: 1 },
        element{ value: "jigglypuff",    position: 2 },
        element{ value: "r.o.b.",        position: 3 },
        element{ value: "mega man",      position: 4 },
        element{ value: "yoshi",         position: 5 },
    }

    t.Run("Add", func(t *testing.T) {
        list    := NewArrayListOf(elements)
        element := element{ value: "gumball", position: list.Size() }
        err     := list.Add(element)

        assertError(t, err, nil)
        assertSize(t, list, 7)
        assertContains(t, list, element, true)
        assertIndex(t, list, element, 6)
    })

    t.Run("AddFirst", func(t *testing.T) {
       list    := NewArrayListOf(elements)
       element := element{ value: "luffy", position: 0 }
       err     := list.AddFirst(element)

       assertError(t, err, nil)
       assertSize(t, list, 7)
       assertContains(t, list, element, true)
       assertIndex(t, list, element, 0)
    })

    t.Run("AddLast", func(t *testing.T) {
       list    := NewArrayListOf(elements)
       element := element{ value: "snorlax", position: list.Size() }
       err     := list.AddLast(element)

       assertError(t, err, nil)
       assertSize(t, list, 7)
       assertContains(t, list, element, true)
       assertIndex(t, list, element, 0)
    })

    t.Run("AddWithIndex", func(t *testing.T) {
       list     := NewArrayListOf(elements)
       index    := 3
       element  := element{ value: "chopper", position: index }
       err      := list.AddWithIndex(index, element)

       assertError(t, err, nil)
       assertSize(t, list, 7)
       assertContains(t, list, element, true)
       assertIndex(t, list, element, index)
    })

    t.Run("AddAll", func(t *testing.T) {
        list        := NewArrayListOf(elements)
        newElements := []interface{}{
            element{ value: "gumball", position: 6 },
            element{ value: "luffy",   position: 7 },
            element{ value: "chopper", position: 8 },
        }
        err := list.AddAll(NewArrayListOf(newElements))

        assertError(t, err, nil)
        assertSize(t, list, 9)

        assertContains(t, list, newElements[0], true)
        e, _ := newElements[0].(element)
        assertIndex(t, list, newElements[0], e.position)

        assertContains(t, list, newElements[1], true)
        e, _ = newElements[1].(element)
        assertIndex(t, list, newElements[1], e.position)

        assertContains(t, list, newElements[2], true)
        e, _ = newElements[2].(element)
        assertIndex(t, list, newElements[2], e.position)
    })
}

func TestArrayList_Remove(t *testing.T) {
    elements := []interface{}{
        element{ value: "piranha plant", position: 0 },
        element{ value: "samus",         position: 1 },
        element{ value: "jigglypuff",    position: 2 },
        element{ value: "r.o.b.",        position: 3 },
        element{ value: "mega man",      position: 4 },
        element{ value: "yoshi",         position: 5 },
    }

    t.Run("Remove", func(t *testing.T) {
        list    := NewArrayListOf(elements)
        element := elements[3]

        if !list.Remove(element) {
            t.Error("expected result to be true")
        }

        if list.Remove(element) {
            t.Error("expected result to be false")
        }

        assertSize(t, list, 5)
        assertContains(t, list, element, false)
    })

    t.Run("RemoveFirst", func(t *testing.T) {
        list    := NewArrayListOf(elements)
        element := list.RemoveFirst()

        if element == nil {
            t.Error("expected element after remove but was nil")
        }

        assertSize(t, list, 5)
        assertContains(t, list, element, false)
        assertIndex(t, list, elements[1], 0)
    })

    t.Run("RemoveLast", func(t *testing.T) {
        list    := NewArrayListOf(elements)
        element := list.RemoveLast()

        if element == nil {
            t.Error("expected element after remove but was nil")
        }

        assertSize(t, list, 5)
        assertContains(t, list, element, false)
    })

    t.Run("RemoveWithIndex", func(t *testing.T) {
        list         := NewArrayListOf(elements)
        element, err := list.RemoveWithIndex(4)

        assertError(t, err, nil)
        assertSize(t, list, 5)
        assertContains(t, list, element, false)
        assertIndex(t, list, elements[5], 4)
    })

    t.Run("Clear", func(t *testing.T) {
        list := NewArrayListOf(elements)

        assertSize(t, list, 6)

        if list.IsEmpty() {
            t.Error("expected result to be false")
        }

        list.Clear()

        if !list.IsEmpty() {
            t.Error("expected result to be true")
        }
    })
}

func TestArrayList_Functional(t *testing.T) {
    elements := []interface{}{
        "piranha plant",
        element{ value: "samus", position: 1 },
        "jigglypuff",
        nil,
        element{ value: "mega man", position: 4 },
        element{ value: "yoshi", position: 5 },
    }

    t.Run("Filter,Map,ForEach", func(t *testing.T) {
        byElementType := func(v interface{}) bool {
            switch v.(type) {
            case element:
                return true
            default:
                return false
            }
        }

        valueToUpperCase := func(v interface{}) interface{} {
            e, _ := v.(element)
            e.value = strings.ToUpper(e.value)
            return e
        }

        list := NewArrayList()
        collect := func(v interface{}) {
            _ = list.Add(v)
        }

        NewArrayListOf(elements).
             Filter(byElementType).
             Map(valueToUpperCase).
             ForEach(collect)

        assertSize(t, list, 3)
        assertContains(t, list, element{ value: "SAMUS", position: 1 }, true)
        assertContains(t, list, element{ value: "MEGA MAN", position: 4 }, true)
        assertContains(t, list, element{ value: "YOSHI", position: 5 }, true)
    })

}

func assertContains(t *testing.T, collection collection.Collection, value interface{}, expected bool) {
    t.Helper()

    if collection.Contains(value) != expected {
        if expected {
            t.Errorf("expected to contain value: %+v", value)
        } else {
            t.Errorf("expected not to contain value: %+v", value)
        }
    }
}

func assertError(t *testing.T, actual error, expected error) {
    t.Helper()

    if actual != expected {
        t.Errorf("expected error '%s', but found '%s'", expected, actual)
    }
}

func assertIndex(t *testing.T, list List, value interface{}, expected int) {
    t.Helper()

    actual, err := list.IndexOf(value)
    if err != nil {
        t.Errorf("expected index of '%d', but found '%d'", expected, actual)
    }
}

func assertSize(t *testing.T, list List, expected int) {
    t.Helper()

    actual := list.Size()
    if actual != expected {
        t.Errorf("expected size of '%d', but found '%d'", expected, actual)
    }
}

