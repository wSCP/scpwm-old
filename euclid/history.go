package main

/*
type Historys struct {
	e *Euclid
	*History
}

type History struct {
	anchor bool
	loc    coordinate
	prev   *History
	next   *History
}

func (h *Historys) Add(m *Monitor, d *Desktop, n *Node) {
	nh := &History{
		loc:    Coordinates(h.e, m, d, n),
		latest: true,
	}
	if h.History == nil {
		nh.anchor = true
		h.History = nh
	} else {
		h.Push(nh, Youngest)
	}
}

func (h *History) Pop(a Age) *History {
	switch a {
	case Youngest:
		curr := h
		var tail *History
		for curr != nil {
			tail = curr
			curr = curr.next
		}
		return last
	case Younger:
		return h.next
	case Older:
		return h.prev
	case Oldest:
		curr := h
		var head *History
		for curr != nil {
			head = curr
			curr = curr.prev
		}
		return head
	}
}

func (h *History) Push(o *History, a Age) {
	switch a {
	case Youngest:
		h.Pop(Youngest).next = o
	case Younger:
		//-->h.next = o
		//next = h.next
		//prev = h.prev

		//next.prev = o
		//o.next = next
		//o.prev = prev
		//prev.next = o
	case Older:
		//-->h.prev = o
		//next = h.next
		//prev = h.prev
	case Oldest:
		old := h.Pop(Oldest)
		o.next = old
		old.prev = o
	}
}

/*
func (h *History) Detach() {
	next = h.next
	prev = h.prev

	next.prev = prev
	prev.next = next

	h.next, h.prev = nil, nil
}

func (h *History) transferNode(m *Monitor, d *Desktop, n *Node) {
	head := h.Head()
	for head != nil {
		if head.loc.n == n {
			head.loc.m = m
			head.loc.d = d
		}
		head = head.next
	}
}

func (h *history) transferDesktop(m *Monitor, d *Desktop) {
	//void history_transfer_Desktop(Monitor_t *m, Desktop_t *d);
}

func (h *history) swapNodes(n1, n2 *Node) {
	//void history_swap_nodes(Monitor_t *m1, Desktop_t *d1, node_t *n1, Monitor_t *m2, Desktop_t *d2, node_t *n2);
}

func (h *history) swapDesktops(d1, d2 *Desktop) {
	//void history_swap_Desktops(Monitor_t *m1, Desktop_t *d1, Monitor_t *m2, Desktop_t *d2);
}

func (h *history) remove(d *Desktop, n *Node) {
		   //removing from the newest to the oldest is required
		   //for maintaining the *latest* attribute
			history_t *b = history_tail;
			while (b != NULL) {
				if ((n != NULL && n == b->loc.node) || (n == NULL && d == b->loc.Desktop)) {
					history_t *a = b->next;
					history_t *c = b->prev;
					if (a != NULL) {
						// remove duplicate entries
						while (c != NULL && ((a->loc.node != NULL && a->loc.node == c->loc.node) ||
						       (a->loc.node == NULL && a->loc.Desktop == c->loc.Desktop))) {
							history_t *d = c->prev;
							if (history_head == c)
								history_head = history_tail;
							if (history_needle == c)
								history_needle = history_tail;
							free(c);
							c = d;
						}
						a->prev = c;
					}
					if (c != NULL)
						c->next = a;
					if (history_tail == b)
						history_tail = c;
					if (history_head == b)
						history_head = a;
					if (history_needle == b)
						history_needle = c;
					free(b);
					b = c;
				} else {
					b = b->prev;
				}
			}
}

func (h *history) empty() {
	head := h.Head()
	var next *history
	for head != nil {
		next = head.next
		head = nil
		head = next
	}
}

func (h *history) getNode(d *Desktop, n *Node) *Node {
	curr := h
	for curr != nil {
		if curr.latest && curr.loc.n != nil && h.loc.n != n && h.loc.d == d {
			return curr.loc.n
		}
		h = h.prev
	}
	return nil
}

func (h *history) findNode(a Age, ref, dst *coordinates, sel ClientSelect) bool {
	//bool history_find_node(history_dir_t hdi, coordinates_t *ref, coordinates_t *dst, client_select_t sel);
	return false
}

func (h *history) getDesktop(m *Monitor, d *Desktop) *Desktop {
	curr := h.Tail()
	for curr != nil {
		if curr.latest && curr.loc.d != d && curr.loc.m == m {
			return curr.loc.d
		}
		curr = curr.prev
	}
	return nil
}

func (h *history) findDesktop(a Age, ref, dst *coordinates, sel DesktopSelect) bool {
	//bool history_find_Desktop(history_dir_t hdi, coordinates_t *ref, coordinates_t *dst, Desktop_select_t sel);
	return false
}

func (h *history) getMonitor(m *Monitor) *Monitor {
	curr := h
	for curr != nil {
		if curr.latest && curr.loc.m != m {
			return curr.loc.m
		}
		curr = curr.prev
	}
	return nil
}

func (h *history) findMonitor(a Age, ref, dst *coordinates, sel DesktopSelect) {
	//bool history_find_Monitor(history_dir_t hdi, coordinates_t *ref, coordinates_t *dst, Desktop_select_t sel);
}

func rank(d *Desktop, n *Node) int {
	//int history_rank(Desktop_t *d, node_t *n);
	return 0
}
*/
