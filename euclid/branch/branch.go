// Package branch implements a custom doubly linked list from
// https://golang.org/src/container/list/list.go for creating
// a limited tree structure for scpwm-euclid containing monitor/desktop/client
// iterations
package branch

type Element struct {
	next, prev *Element
	branch     *Branch
	Value      interface{}
}

// Next returns the next branch element or nil.
func (e *Element) Next() *Element {
	if n := e.next; e.branch != nil && n != &e.branch.root {
		return n
	}
	return nil
}

// NextContinuous returns the next *Element continuously, looping to the front
// when the root is reached.
func (e *Element) NextContinuous() *Element {
	n := e.next
	if e.branch != nil {
		if n == &e.branch.root {
			n = e.branch.Front()
		}
		return n
	}
	return nil
}

// Prev returns the previous branch element or nil.
func (e *Element) Prev() *Element {
	if p := e.prev; e.branch != nil && p != &e.branch.root {
		return p
	}
	return nil
}

// PrevContinuous returns the previous *Element continuously, looping to the back
// when the root is reached.
func (e *Element) PrevContinuous() *Element {
	p := e.prev
	if e.branch != nil {
		if p == &e.branch.root {
			p = e.branch.Back()
		}
		return p
	}
	return nil
}

type Branch struct {
	id   string
	root Element
	len  int
}

// Init initializes or clears branch b.
func (b *Branch) Init(id string) *Branch {
	b.id = id
	b.root.next = &b.root
	b.root.prev = &b.root
	b.len = 0
	return b
}

// New returns an initialized branch.
func New(id string) *Branch { return new(Branch).Init(id) }

// Len returns the number of elements of branch b.
// The complexity is O(1).
func (b *Branch) Len() int { return b.len }

// Front returns the first element of branch b or nil.
func (b *Branch) Front() *Element {
	if b.len == 0 {
		return nil
	}
	return b.root.next
}

// Back returns the last element of branch l or nil.
func (b *Branch) Back() *Element {
	if b.len == 0 {
		return nil
	}
	return b.root.prev
}

func (b *Branch) lazyInit() {
	if b.root.next == nil {
		b.Init(b.id)
	}
}

func (b *Branch) insert(e, at *Element) *Element {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.branch = b
	b.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (b *Branch) insertValue(v interface{}, at *Element) *Element {
	return b.insert(&Element{Value: v}, at)
}

func (b *Branch) remove(e *Element) *Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.branch = nil
	b.len--
	return e
}

// Remove removes e from b if e is an element of branch b.
// It returns the element value e.Value.
func (b *Branch) Remove(e *Element) interface{} {
	if e.branch == b {
		// if e.branch == b, b must have been initialized when e was inserted
		// in b or b == nil (e is a zero Element) and b.remove will crash
		b.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of branch b and returns e.
func (b *Branch) PushFront(v interface{}) *Element {
	b.lazyInit()
	return b.insertValue(v, &b.root)
}

// PushBack inserts a new element e with value v at the back of branch b and returns e.
func (b *Branch) PushBack(v interface{}) *Element {
	b.lazyInit()
	return b.insertValue(v, b.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of b, the branch is not modified.
func (b *Branch) InsertBefore(v interface{}, mark *Element) *Element {
	if mark.branch != b {
		return nil
	}
	// see comment in Branch.Remove about initialization of b
	return b.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of b, the branch is not modified.
func (b *Branch) InsertAfter(v interface{}, mark *Element) *Element {
	if mark.branch != b {
		return nil
	}
	// see comment in Branch.Remove about initialization of l
	return b.insertValue(v, mark)
}

// MoveToFront moves element e to the front of branch l.
// If e is not an element of l, the branch is not modified.
func (b *Branch) MoveToFront(e *Element) {
	if e.branch != b || b.root.next == e {
		return
	}
	// see comment in Branch.Remove about initialization of l
	b.insert(b.remove(e), &b.root)
}

// MoveToBack moves element e to the back of branch b.
// If e is not an element of b, the branch is not modified.
func (b *Branch) MoveToBack(e *Element) {
	if e.branch != b || b.root.prev == e {
		return
	}
	// see comment in Branch.Remove about initialization of b
	b.insert(b.remove(e), b.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of b, or e == mark, the branch is not modified.
func (b *Branch) MoveBefore(e, mark *Element) {
	if e.branch != b || e == mark || mark.branch != b {
		return
	}
	b.insert(b.remove(e), mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of b, or e == mark, the branch is not modified.
func (b *Branch) MoveAfter(e, mark *Element) {
	if e.branch != b || e == mark || mark.branch != b {
		return
	}
	b.insert(b.remove(e), mark)
}

// PushBackBranch inserts a copy of an other branch at the back of branch b.
// The branchs b and other may be the same.
func (b *Branch) PushBackBranch(other *Branch) {
	b.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		b.insertValue(e.Value, b.root.prev)
	}
}

// PushFrontBranch inserts a copy of an other branch at the front of branch b.
// The branchs b and other may be the same.
func (b *Branch) PushFrontBranch(other *Branch) {
	b.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		b.insertValue(e.Value, &b.root)
	}
}
