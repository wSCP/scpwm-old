package client

type Shiftr interface{}

type shiftr struct {
	rotate int
	splitR float64 // 0.01 - 0.99
	splitM string  //splitMode // Automatic, Manual
	splitT string  //Orientation // Horizontal, Vertical
	splitD string  //Direction   // Right, Down, Left, Up
}

//type splitMode int

//const (
//	Automatic splitMode = iota
//	Manual
//)
