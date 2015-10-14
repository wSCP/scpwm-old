package main

//type Rule struct {
//	cause  string
//	effect string
//	reuse  bool
//	prev   *Rule
//	next   *Rule
//}

//func (r *Rule) Pop(cy Cycle) *Rule {
//	return nil
//}

//func (r *Rule) Push(cy Cycle, o *Rule) {}

type Rule struct {
	className    string
	instanceName string
	monitorDesc  string
	desktopDesc  string
	//nodeDesc     string
	splitD      string
	splitR      float64
	minWidth    uint16
	maxWidth    uint16
	minHeight   uint16
	maxHeight   uint16
	pseudoTiled bool
	floating    bool
	locked      bool
	sticky      bool
	private     bool
	center      bool
	follow      bool
	manage      bool
	focus       bool
	border      bool
}

func NewRule() *Rule {
	return &Rule{
		manage: true,
		focus:  true,
		border: true,
	}
}

//func (c *Consequence) Pending() *Pending {
/*
	pending_rule_t *make_pending_rule(int fd, xcb_window_t win, rule_consequence_t *csq)
	{
	pending_rule_t *pr = malloc(sizeof(pending_rule_t));
	pr->prev = pr->next = NULL;
	pr->fd = fd;
	pr->win = win;
	pr.consequence = c
	return pr;
	}
*/
//	return nil
//}

//func (e *Euclid) applyRules(win *Window, csq *Consequence) {}

//func (e *Euclid) scheduleRules(win *Window, csq *Consequence) bool {
//	return false
//}

//type Pending struct {
//int fd;
//window      xproto.Window
//consequence *Consequence
//prev        *Pending
//next        *Pending
//}

//func (p *Pending) Pop(cy Cycle) *Pending {
//return nil
//}

//func (p *Pending) Push(cy Cycle, o *Pending) {}
