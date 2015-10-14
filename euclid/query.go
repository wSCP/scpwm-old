package main

import "bytes"

/*type domain int

const (
	DMONITOR domain = iota
	DDESKTOP
	DWINDOW
	DTREE
	DHISTORY
	DSTACK
)*/

type query struct {
	*bytes.Buffer
}

func Query() *query {
	return &query{new(bytes.Buffer)}
}

//func (q *query) monitors(loc *Coordinate) *query {
//void query_monitors(coordinates_t loc, domain_t dom, FILE *rsp);
//	return q
//}

//func (q *query) desktops(m *Monitor, loc *Coordinate, depth uint, dom domain) *query {
//void query_desktops(monitor_t *m, domain_t dom, coordinates_t loc, unsigned int depth, FILE *rsp);
//	return q
//}

//func (q *query) Node(d *Desktop, n *Node, depth uint) *query {
//void query_tree(desktop_t *d, node_t *n, FILE *rsp, unsigned int depth);
//return q
//}

//func (q *query) History(loc *Coordinate) *query {
//void query_history(coordinates_t loc, FILE *rsp);
//return q
//}

//func (q *query) Stack() *query {
//void query_stack(FILE *rsp);
//return q
//}

//func (q *query) Windows(loc *Coordinate) *query {
//void query_windows(coordinates_t loc, FILE *rsp);
//	return q
//}

//func (q *query) Pointer(ptr *Pointer) *query {
//	return q
//}
