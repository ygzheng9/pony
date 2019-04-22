package base

import (
	packr "github.com/gobuffalo/packr/v2"
)

// ABox for base, will packr to binary file
var ABox = packr.New("app:baseBox", "../templates/abox")
