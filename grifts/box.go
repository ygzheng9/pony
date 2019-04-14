package grifts

import (
	"fmt"
	"pony/base"

	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("box", func() {
	grift.Desc("list", "List all files in baseBox")
	grift.Add("list", func(c *grift.Context) error {
		files := base.Box.List()
		fmt.Println(files)

		fmt.Println("completed. ")
		return nil
	})
})
