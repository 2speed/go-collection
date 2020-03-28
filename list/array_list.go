package list

import (
    "fmt"
    "reflect"
    "strings"

    "github.com/2speed/go-collection"
    "github.com/pkg/errors"
)

// ArrayList is a simple implementation of a List whose elements are maintained by an internal slice. ArrayList does
// not make any guarantees for concurrent access.
//
// Example:
//
//   elements := []interface{}{ "one", 1, "two", 2, "three", 3 } // create a slice of elements of varying types
//
//   NewArrayListOf(elements).
//      Filter(func(element interface{}) bool {      // filter elements that are not of type string
//          _, ok := element.(string)
//          return ok
//      }).
//      Map(func(element interface{}) interface{} {  // map each string element to it's upper case variant
//          return strings.ToUpper(element.(string))
//      }).
//      ForEach(func(element interface{}) {          // print each result
//          fmt.Println(element)
//      })
//
//   Result:
//
//     ONE
//     TWO
//     THREE
//
type arrayList struct {
    elements []interface{}
}

// NewArrayList creates a new ArrayList.
func NewArrayList() List {
    return &arrayList{ elements: make([]interface{}, 0) }
}

// NewArrayListOf creates a new ArrayList containing the provided elements.
func NewArrayListOf(elements []interface{}) List {
    l := NewArrayList()
    if elements != nil {
        for _, e := range elements {
            _ = l.Add(e)
        }
    }

    return l
}

// NewArrayListFrom creates a new ArrayList containing the elements from the provided Collection.
func NewArrayListFrom(collection collection.Collection) List {
    var elements []interface{}
    if collection != nil {
        elements = collection.Values()
    }

    return NewArrayListOf(elements)
}

// Add inserts the provided element into the ArrayList.
func (l *arrayList) Add(element interface{}) error {
    l.elements = append(l.elements, element)

    return nil
}

// AddAll inserts all elements from the provided List into the ArrayList.
func (l *arrayList) AddAll(collection collection.Collection) error {
    if collection != nil {
        l.elements = append(l.elements, collection.Values()...)
    }

    return nil
}

// AddFirst inserts the provided element at the front (index == 0) of the ArrayList. The positions of the existing
// elements are increased by one.
func (l *arrayList) AddFirst(element interface{}) error {
    l.elements = append([]interface{}{element }, l.elements...)

    return nil
}

// AddLast inserts the provided element at the end of the ArrayList (index == ArrayList.Size()).
func (l *arrayList) AddLast(element interface{}) error {
    return l.Add(element)
}

// AddWithIndex inserts the provided element into the ArrayList specified by index. The position of the elements that
// were at positions index to ArrayList.Size() - 1 increase by one. The returned error will be non-nil if the provided
// index is outside the current bounds of the ArrayList (index < 0 || index > ArrayList.Size() - 1).
func (l *arrayList) AddWithIndex(index int, element interface{}) error {
    if err := l.checkBounds(index); err != nil {
        return err
    }

    l.elements = append(l.elements, nil)
    copy(l.elements[index + 1:], l.elements[index:])
    l.elements[index] = element

    return nil
}

// ValueWithIndex returns the element at the position specified by the provided index. The returned error will be
// non-nil if the provided index is outside the current bounds of the ArrayList
// (index < 0 || index > ArrayList.Size() - 1).
func (l *arrayList) ValueWithIndex(index int) (interface{}, error) {
    if err := l.checkBounds(index); err != nil {
        return nil, err
    }

    return l.elements[index], nil
}

// IndexOf returns the position of the first occurrence (if any) of an element equivalent to the provided element. The
// returned error will be non-nil if provided element is not found in the ArrayList, and the returned index will be
// equal to collection.ErrorElementNotFound.
func (l *arrayList) IndexOf(element interface{}) (int, error) {
    i, err := l.findFirst(element)
    if err != nil {
        return i, err
    }

    return i, nil
}

// Remove removes the first occurrence (if any) of an element equivalent to the provided element. If an element was
// removed, the return value will be true, otherwise false will be returned.
func (l *arrayList) Remove(element interface{}) bool {
    if l.Contains(element) {
        i, _ := l.IndexOf(element)
        if _, err := l.RemoveWithIndex(i); err != nil {
            return false
        }

        return true
    }

    return false
}

// RemoveFirst removes the element at the front (index == 0) of the ArrayList and returns it. If the ArrayList is empty
// (ArrayList.Size() == 0), the return value will be nil.
func (l *arrayList) RemoveFirst() interface{} {
    if l.Size() > 0 {
        v, _ := l.RemoveWithIndex(0)

        return v
    }

    return nil
}

// RemoveLast removes the element at the end (index == ArrayList.Size() - 1) of the ArrayList and returns it. If the
// ArrayList is empty (List.Size() == 0), the return value will be nil.
func (l *arrayList) RemoveLast() interface{} {
    if l.Size() > 0 {
        v, _ := l.RemoveWithIndex(l.Size() - 1)

        return v
    }

    return nil
}

// RemoveWithIndex removes the element at the provided index from the ArrayList and returns it. The positions of the
// elements originally at positions index + 1 to ArrayList.Size() - 1 are decremented by 1. The returned error will be
// non-nil if the provided index is outside the bounds of the ArrayList (index < 0 || index > ArrayList.Size() - 1).
func (l *arrayList) RemoveWithIndex(index int) (interface{}, error) {
    element, err := l.ValueWithIndex(index)
    if err != nil {
        return nil, err
    }

    copy(l.elements[index:l.Size() - 1], l.elements[index + 1:l.Size()])
    l.elements = l.elements[:l.Size() - 1]

    return element, nil
}

// Filter returns a new ArrayList consisting of the elements of this ArrayList that match the given predicate.
func (l *arrayList) Filter(predicate func(element interface{}) bool) List {
    list := NewArrayList()

    l.ForEach(func(element interface{}) {
        if predicate(element) {
            _ = list.Add(element)
        }
    })

    return list
}

// Map returns a new ArrayList containing the resulting elements of applying the given function to the elements of this
// ArrayList.
func (l *arrayList) Map(mapper func(element interface{}) interface{}) List {
    list := NewArrayList()

    l.ForEach(func(element interface{}) { list.Add(mapper(element)) })

    return list
}

// ForEach performs the provided consumer function for each element of the ArrayList.
func (l *arrayList) ForEach(consumer func(element interface{})) {
    for _, v := range l.elements {
        consumer(v)
    }
}

// Size returns the number of elements in the ArrayList.
func (l *arrayList) Size() int {
    return len(l.elements)
}

// IsEmpty returns true if the ArrayList contains no elements, otherwise false is returned.
func (l *arrayList) IsEmpty() bool {
    return l.Size() == 0
}

// Clear removes all elements from the ArrayList.
func (l *arrayList) Clear() {
    l.elements = l.elements[:0]
}

// Contains returns true if an element equivalent to the provided element exists in the ArrayList, otherwise false is
// returned.
func (l *arrayList) Contains(element interface{}) bool {
    if index, err := l.IndexOf(element); err == nil && index != collection.ElementNotFound {
        return true
    }

    return false
}

// Values returns a slice containing the elements in the List in the iteration order.
func (l *arrayList) Values() []interface{} {
    elements := make([]interface{}, l.Size())
    copy(elements, l.elements)

    return elements
}

// String returns a string representation of the ArrayList in it's current state.
func (l *arrayList) String() string {
    if l.Size() == 0 {
        return "[]"
    }

    elements := make([]string, 0, l.Size())
    l.ForEach(func(element interface{}) {
        elements = append(elements, fmt.Sprintf("%v", element))
    })

    return "[" + strings.Join(elements, ", ") + "]"
}

func (l *arrayList) checkBounds(index int) error {
    if index < 0 || index > l.Size() {
        return errors.Errorf("index out of bounds [*ArrayList.Size() = %v, requested index = %v]", l.Size(), index)
    }

    return nil
}

func (l *arrayList) findFirst(element interface{}) (int, error) {
    for i, v := range l.elements {
        if reflect.DeepEqual(v, element) {
            return i, nil
        }
    }

    return -1, collection.ErrorElementNotFound
}