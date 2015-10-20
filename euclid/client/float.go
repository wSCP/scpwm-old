package client

import "github.com/BurntSushi/xgb/xproto"

type Floatr interface{}

type floatr struct {
	rectangle xproto.Rectangle
}

func (f *floatr) UpdateRectangle() {
	//void update_floating_rectangle(client_t *c);
}

func (f *floatr) UpdateWidth(w int) {
	//void restrain_floating_width(client_t *c, int *width);
}

func (f *floatr) UpdateHeight(h int) {
	//void restrain_floating_height(client_t *c, int *height);
}

func (f *floatr) UpdateSize(w, h int) {
	//void restrain_floating_size(client_t *c, int *width, int *height);
}
