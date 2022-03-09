package orderedmap

// Element is an element of a doubly linked list that contains the key of the correspondent element in the ordered map too.
type Element[K comparable, V any] struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element[K, V]

	// The list to which this element belongs.
	list *list[K, V]

	// The key that corresponds to this element in the ordered map.
	Key K

	// The value stored with this element.
	Value V
}

// Next returns the next list element or nil.
func (e *Element[K, V]) Next() *Element[K, V] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Element[K, V]) Prev() *Element[K, V] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// list represents a doubly linked list.
type list[K comparable, V any] struct {
	root Element[K, V] // sentinel list element, only &root, root.prev, and root.next are used
}

// Init initializes or clears list l.
func (l *list[K, V]) Init() {
	l.root.next = &l.root
	l.root.prev = &l.root
}

func (l *list[K, V]) IsEmpty() bool {
	return l.root.next == &l.root
}

// Front returns the first element of list l or nil if the list is empty.
func (l *list[K, V]) Front() *Element[K, V] {
	if l.IsEmpty() {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *list[K, V]) Back() *Element[K, V] {
	if l.IsEmpty() {
		return nil
	}
	return l.root.prev
}

// insert inserts e after at, and returns e.
func (l *list[K, V]) insert(e, at *Element[K, V]) *Element[K, V] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	return e
}

// Remove removes e from its list
func (l *list[K, V]) Remove(e *Element[K, V]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *list[K, V]) PushFront(key K, value V) *Element[K, V] {
	return l.insert(&Element[K, V]{Key: key, Value: value}, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *list[K, V]) PushBack(key K, value V) *Element[K, V] {
	return l.insert(&Element[K, V]{Key: key, Value: value}, l.root.prev)
}
