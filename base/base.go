package base

import (
	packr "github.com/gobuffalo/packr/v2"
)

// Box for base, will packr to binary file
var Box = packr.New("app:baseBox", "../baseBox")
