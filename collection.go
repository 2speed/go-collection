package collection

const ElementNotFound = -1

const (
    ErrorElementNotFound = CollectionError("the requested element could not be found")
)

type CollectionError string

func (e CollectionError) Error() string {
    return string(e)
}

// Collection defines the behavior for maintaining a collection of elements.
type Collection interface {

    // Add inserts the provided element into the Collection. The returned error will be non-nil for bounded Collection
    // implementations that have reached capacity and cannot hold any further elements.
    Add(element interface{}) error

    // AddAll inserts all elements from the provided collection into the Collection. The returned error will be non-nil
    // for bounded Collection implementations that have reached capacity and cannot hold any further elements.
    AddAll(collection Collection) error

    // Remove removes the first occurrence (if any) of an element equivalent to the provided element. If an element was
    // removed, the return value will be true, otherwise false will be returned.
    Remove(element interface{}) bool

    // Size returns the number of elements in the Collection.
    Size() int

    // IsEmpty returns true if the Collection contains no elements, otherwise false is returned.
    IsEmpty() bool

    // Clear removes all elements from the Collection.
    Clear()

    // Contains returns true if an element equivalent to the provided element exists in the Collection, otherwise false
    // is returned.
    Contains(value interface{}) bool

    // Values returns a slice containing the elements in the Collection in the iteration order.
    Values() []interface{}
}

// Ordered defines the behavior for a Collection whose elements are algorithmically positioned.
type Ordered interface {
    Collection

    // Min returns the element with the lowest position in the Collection. More specifically, the first element in the
    // iteration order is returned.
    Min() interface{}

    // Max returns the element with the highest position in the Collection. More specifically, the last element in the
    // iteration order is returned.
    Max() interface{}

    // Predecessor returns the element (if any) from the Collection that is less than the provided element. More
    // specifically, the element before the first occurrence of the provided element in iteration order is returned.
    Predecessor(element interface{}) interface{}

    // Successor returns the element (if any) from the Collection that is greater than the provided element. More
    // specifically, the element after the first occurrence of the provided element in iteration order is returned.
    Successor(element interface{}) interface{}
}