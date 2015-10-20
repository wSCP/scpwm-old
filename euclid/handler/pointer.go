package handler

/*
import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type PointerAction int

const (
	NoAction PointerAction = iota
	AFocus
	AMove
	AResizeSide
	AResizeCorner
)

type corner int

const (
	TopLeft corner = iota
	TopRight
	BottomRight
	BottomLeft
)

type side int

const (
	LeftSide side = iota
	TopSide
	RightSide
	BottomSide
)

type Pointer interface {
	Grab(PointerAction)
	Ungrab()
	Center(xproto.Rectangle)
	Track(int16, int16)
}

type pointer struct {
	loc       coordinate
	window    *Window
	position  xproto.Point
	action    PointerAction
	rectangle xproto.Rectangle
	vFence    *Client
	hFence    *Client
	IsTiled   bool
	vRatio    float64
	hRatio    float64
	corner    corner
	side      side
	*MotionRecorder
}

func NewPointer(m *MotionRecorder) Pointer {
	return &pointer{
		window:         &NoWindow,
		action:         NoAction,
		MotionRecorder: m,
	}
}

func (p *pointer) Center(r xproto.Rectangle) {
	cx := r.X + int16(r.Width/2)
	cy := r.Y + int16(r.Height/2)
	p.Lower()
	c := p.window.Conn
	xproto.WarpPointer(c, xproto.WindowNone, p.root, 0, 0, 0, 0, cx, cy)
	p.Raise()
}

func (p *pointer) query(w *Window, pt *xproto.Point) {
	p.Lower()

	qpr, _ := xproto.QueryPointer(p.window.Conn, p.root).Reply()

	if qpr != nil {
		if w != nil {
			w.Window = qpr.Child
		}
		if pt != nil {
			pt.X = qpr.RootX
			pt.Y = qpr.RootY
		}
	}

	p.Raise()
}

func (p *pointer) Grab(a PointerAction) {
	win := &NoWindow

	var pos xproto.Point

	p.query(win, &pos)

	//var loc coordinate

	//if _, exists := locateWindow(p.loc.e, win.Window); exists {
	//var c *Client
	//p.position = pos
	//p.action = a
	//c = loc.n.Client
	//p.monitor = loc.m
	//p.desktop = loc.d
	//p.node = loc.n
	//p.client = c
	//p.window = c.Window
	//p.hFence = nil
	//p.vFence = nil

	//switch a {
	//case AFocus:
	//	if (loc.node != mon->desk->focus) {
	//		bool backup = pointer_follows_monitor;
	//		pointer_follows_monitor = false;
	//		focus_node(loc.monitor, loc.desktop, loc.node);
	//		pointer_follows_monitor = backup;
	//	} else if (focus_follows_pointer) {
	//		stack(loc.node, STACK_ABOVE);
	//	}
	//	frozen_pointer->action = ACTION_NONE;
	//	break;
	//case AMove, AResizeSide, AResizeCorner:
	//if (is_floating(c)) {
	//	frozen_pointer->rectangle = c->floating_rectangle;
	//	frozen_pointer->is_tiled = false;
	//} else if (is_tiled(c)) {
	//	frozen_pointer->rectangle = c->tiled_rectangle;
	//	frozen_pointer->is_tiled = (pac == ACTION_MOVE || !c->pseudo_tiled);
	//} else {
	//	frozen_pointer->action = ACTION_NONE;
	//	return;
	//}
	//if (pac == ACTION_RESIZE_SIDE) {
	//	float W = frozen_pointer->rectangle.width;
	//	float H = frozen_pointer->rectangle.height;
	//	float ratio = W / H;
	//	float x = pos.x - frozen_pointer->rectangle.x;
	//	float y = pos.y - frozen_pointer->rectangle.y;
	//	float diag_a = ratio * y;
	//	float diag_b = W - diag_a;
	//	if (x < diag_a) {
	//		if (x < diag_b)
	//			frozen_pointer->side = SIDE_LEFT;
	//		else
	//			frozen_pointer->side = SIDE_BOTTOM;
	//	} else {
	//		if (x < diag_b)
	//			frozen_pointer->side = SIDE_TOP;
	//		else
	//			frozen_pointer->side = SIDE_RIGHT;
	//	}
	//} else if (pac == ACTION_RESIZE_CORNER) {
	//	int16_t mid_x = frozen_pointer->rectangle.x + (frozen_pointer->rectangle.width / 2);
	//	int16_t mid_y = frozen_pointer->rectangle.y + (frozen_pointer->rectangle.height / 2);
	//	if (pos.x > mid_x) {
	//		if (pos.y > mid_y)
	//			frozen_pointer->corner = CORNER_BOTTOM_RIGHT;
	//		else
	//			frozen_pointer->corner = CORNER_TOP_RIGHT;
	//	} else {
	//		if (pos.y > mid_y)
	//			frozen_pointer->corner = CORNER_BOTTOM_LEFT;
	//		else
	//			frozen_pointer->corner = CORNER_TOP_LEFT;
	//	}
	//}
	//if (frozen_pointer->is_tiled) {
	//	if (pac == ACTION_RESIZE_SIDE) {
	//		switch (frozen_pointer->side) {
	//			case SIDE_TOP:
	//				frozen_pointer->horizontal_fence = find_fence(loc.node, DIR_UP);
	//				break;
	//			case SIDE_RIGHT:
	//				frozen_pointer->vertical_fence = find_fence(loc.node, DIR_RIGHT);
	//				break;
	//			case SIDE_BOTTOM:
	//				frozen_pointer->horizontal_fence = find_fence(loc.node, DIR_DOWN);
	//				break;
	//			case SIDE_LEFT:
	//				frozen_pointer->vertical_fence = find_fence(loc.node, DIR_LEFT);
	//				break;
	//		}
	//	} else if (pac == ACTION_RESIZE_CORNER) {
	//		switch (frozen_pointer->corner) {
	//			case CORNER_TOP_LEFT:
	///				frozen_pointer->horizontal_fence = find_fence(loc.node, DIR_UP);
	//				frozen_pointer->vertical_fence = find_fence(loc.node, DIR_LEFT);
	//				break;
	//			case CORNER_TOP_RIGHT:
	//				frozen_pointer->horizontal_fence = find_fence(loc.node, DIR_UP);
	//				frozen_pointer->vertical_fence = find_fence(loc.node, DIR_RIGHT);
	//				break;
	//			case CORNER_BOTTOM_RIGHT:
	//				frozen_pointer->horizontal_fence = find_fence(loc.node, DIR_DOWN);
	//				frozen_pointer->vertical_fence = find_fence(loc.node, DIR_RIGHT);
	//				break;
	//			case CORNER_BOTTOM_LEFT:
	//				frozen_pointer->horizontal_fence = find_fence(loc.node, DIR_DOWN);
	//				frozen_pointer->vertical_fence = find_fence(loc.node, DIR_LEFT);
	//				break;
	//		}
	//	}
	//	if (frozen_pointer->horizontal_fence != NULL)
	//		frozen_pointer->horizontal_ratio = frozen_pointer->horizontal_fence->split_ratio;
	//	if (frozen_pointer->vertical_fence != NULL)
	//		frozen_pointer->vertical_ratio = frozen_pointer->vertical_fence->split_ratio;
	//}
	//}
	//} else {
	//if a == AFocus {
	//	monitor_t *m = monitor_from_point(pos);
	//	if (m != NULL && m != mon)
	//		focus_node(m, m->desk, m->desk->focus);
	//}
	//p.action = NoAction
	//}
}

func (p *pointer) Ungrab() {
	p.action = NoAction
}

func (p *pointer) Track(x, y int16) {
	if p.action != NoAction {
		//var dx, dy, x, y int16
		//var w int16 = 1
		//var h int16 = 1
		//dx = x - p.position.X
		//dy = y - p.position.Y

		//switch p.action {
		//case AMove:
		//	if p.IsTiled {
		//		pwin := &NoWindow
		//		p.query(pwin, nil)
		//		if pwin != p.window {
		//			var loc Coordinate
		//			var isManaged bool
		//			if pwin.Window != xproto.WindowNone {
		//				isManaged = LocateWindow(p.e, pwin, &loc)
		//			}
		//			if isManaged && loc.n.IsTiled() && loc.m == p.monitor {
		//				//SwapNodes(p.monitor, p.monitor, p.desktop, p.desktop, p.node, loc.n)
		//				p.desktop.Arrange()
		//			} else {
		//				if isManaged && loc.m == p.monitor {
		//					return
		//				} else if !isManaged {
		//					//pmon := p.monitor.FromPoint(xproto.Point{x, y})
		//					//if pmon == nil || pmon == p.monitor {
		//					//	return
		//					//} else {
		//					//	loc.m = pmon
		//					//	loc.d = pmon.desk
		//					//}
		//				}
		//				//focused := (p.node == mon.desk.focus)
		//				//TransferNode(p.monitor, loc.m, p.desktop, loc.d, p.node, loc.d.focus)
		//				//if focused {
		//				//focus_node(loc.monitor, loc.desktop, n)
		//				//}
		//				p.monitor = loc.m
		//				p.desktop = loc.d
		//			}
		//		}
		//	} else {
		//		nx := p.rectangle.X + dx
		//		ny := p.rectangle.Y + dy
		//		p.window.Move(nx, ny)
		//		p.client.fRectangle.X = nx
		//		p.client.fRectangle.Y = ny
		//		//pmon := fromPoint(xproto.Point{x, y})
		//		//if pmon == nil || pmon == p.monitor {
		//		//	return
		//		//}
		//		//focused := (p.node == mon.desk.focus)
		//		//transferNode(p.monitor, pmon, p.desktop, pmon.desk, p.node, pmon.desk.focus)
		//		//if focused {
		//		//focus_node(pmon, pmon->desk, n)
		//		//}
		//		//p.monitor = pmon
		//		//p.desktop = pmon.desk
		//	}
		//case AResizeSide, AResizeCorner:
		//	if p.IsTiled {
		//		if p.vFence != nil {
		//			sr := p.vRatio + float64(dx/int16(p.vFence.Rectangle.Width))
		//			sr = fmax(0, sr)
		//			sr = fmin(1, sr)
		//			p.vFence.splitR = sr
		//		}
		//		if p.hFence != nil {
		//			sr := p.hRatio + float64(dy/int16(p.hFence.Rectangle.Height))
		//			sr = fmax(0, sr)
		//			sr = fmin(1, sr)
		//			p.hFence.splitR = sr
		//		}
		//		p.desktop.Arrange()
		//	} else {
		//		if p.action == AResizeSide {
		//			switch p.side {
		//			case TopSide:
		//				x = p.rectangle.X
		//				y = p.rectangle.Y + dy
		//				w = int16(p.rectangle.Width)
		//				h = int16(p.rectangle.Height) - dy
		//			case RightSide:
		//				x = p.rectangle.X
		//				y = p.rectangle.Y
		//				w = int16(p.rectangle.Width) + dx
		//				h = int16(p.rectangle.Height)
		//			case BottomSide:
		//				x = p.rectangle.X
		//				y = p.rectangle.Y
		//				w = int16(p.rectangle.Width)
		//				h = int16(p.rectangle.Height) + dy
		//			case LeftSide:
		//				x = p.rectangle.X + dx
		//				y = p.rectangle.Y
		//				w = int16(p.rectangle.Width) - dx
		//				h = int16(p.rectangle.Height)
		//			}
		//		} else if p.action == AResizeCorner {
		//			switch p.corner {
		//			case TopLeft:
		//				x = p.rectangle.X + dx
		//				y = p.rectangle.Y + dy
		//				w = int16(p.rectangle.Width) - dx
		//				h = int16(p.rectangle.Height) - dy
		//			case TopRight:
		//				x = p.rectangle.X
		//				y = p.rectangle.Y + dy
		//				w = int16(p.rectangle.Width) + dx
		//				h = int16(p.rectangle.Height) - dy
		//			case BottomLeft:
		//				x = p.rectangle.X + dx
		//				y = p.rectangle.Y
		//				w = int16(p.rectangle.Width) - dx
		//				h = int16(p.rectangle.Height) + dy
		//			case BottomRight:
		//				x = p.rectangle.X
		//				y = p.rectangle.Y
		//				w = int16(p.rectangle.Width) + dx
		//				h = int16(p.rectangle.Height) + dy
		//			}
		//		}

		//		oldw := w
		//		oldh := h

		//restrainFloatingSize(p.client, &w, &h)

		//		if p.client.pseudoTiled {
		//			p.client.fRectangle.Width = uint16(w)
		//			p.client.fRectangle.Height = uint16(h)
		//			p.desktop.Arrange()
		//		} else {
		//			if oldw == w {
		//				p.client.fRectangle.X = x
		//				p.client.fRectangle.Width = uint16(w)
		//			}
		//			if oldh == h {
		//				p.client.fRectangle.Y = y
		//				p.client.fRectangle.Height = uint16(h)
		//			}
		//			p.window.MoveResize(
		//				p.client.fRectangle.X,
		//				p.client.fRectangle.Y,
		//				p.client.fRectangle.Width,
		//				p.client.fRectangle.Height)
		//		}
		//	}
		//}
	}
}

type MotionRecorder struct {
	*Window
}

func NewMotionRecorder(conn *xgb.Conn, r, w xproto.Window) *MotionRecorder {
	return &MotionRecorder{
		Window: &Window{conn, w, r},
	}
}

func (m *MotionRecorder) enable() {
	m.Raise()
	m.Show()
}

func (m *MotionRecorder) disable() {
	m.Hide()
}

func (m *MotionRecorder) update() {
	geo, _ := xproto.GetGeometry(m.Window.Conn, xproto.Drawable(m.root)).Reply()

	if geo != nil {
		m.Resize(geo.Width, geo.Height)
	}
}
*/
