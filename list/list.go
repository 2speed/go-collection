package list

import "github.com/2speed/go-collection"

// List defines the behavior for a container the represents a Collection of elements that are accessed via their
// position much like that of an array or slice. Unlike an array or slice however, elements of a list are more "generic"
// and can be of any type. As a result, it is the responsibility of the caller to perform any necessary type assertions
// for elements returned from accessor methods (e.g. List.ValueWithIndex(index)) if a specific type is required.
type List interface {
    collection.Collection

    // AddFirst inserts the provided element at the front (index == 0) of the List. The positions of the existing
    // elements are increased by one. The returned error will be non-nil for bounded List implementations that have
    // reached capacity and cannot hold any further elements.
    AddFirst(element interface{}) error

    // AddLast inserts the provided element at the end of the List (index == List.Size()). The returned error will be
    // non-nil for bounded List implementations that have reached capacity and cannot hold any further elements.
    AddLast(element interface{}) error

    // AddWithIndex inserts the provided element into the List specified by index. The position of the elements that
    // were at positions index to List.Size() - 1 increase by one. The returned error will be non-nil if the provided
    // index is outside the current bounds of the List (index < 0 || index > List.Size() - 1).
    AddWithIndex(index int, element interface{}) error

    // ValueWithIndex returns the element at the position specified by the provided index. The returned error will be
    // non-nil if the provided index is outside the current bounds of the List (index < 0 || index > List.Size() - 1).
    ValueWithIndex(index int) (interface{}, error)

    // IndexOf returns the position of the first occurrence (if any) of an element equivalent to the provided element.
    // The returned error will be non-nil if provided element is not found in the List, and the returned index will be
    // -1.
    IndexOf(element interface{}) (int, error)

    // RemoveFirst removes the element at the front (index == 0) of the List and returns it. If the List is empty
    // (List.Size() == 0), the return value will be nil.
    RemoveFirst() interface{}

    // RemoveLast removes the element at the end (index == List.Size() - 1) of the List and returns it. If the
    // List is empty (List.Size() == 0), the return value will be nil.
    RemoveLast() interface{}

    // RemoveWithIndex removes the element at the provided index from the List and returns it. The positions of the
    // elements originally at positions index + 1 to List.Size() - 1 are decremented by 1. The returned error will be
    // non-nil if the provided index is outside the bounds of the List (index < 0 || index > List.Size() - 1).
    RemoveWithIndex(index int) (interface{}, error)

    // Filter returns a new List consisting of the elements of this List that match the given predicate.
    Filter(predicate func(element interface{}) bool) List

    // Map returns a new List containing the resulting elements of applying the given function to the elements of this
    // List.
    Map(mapper func(element interface{}) interface{}) List

    // ForEach performs the provided consumer function for each element of the List.
    ForEach(consumer func(element interface{}))
}