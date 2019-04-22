package grifts

import (
	"fmt"
	"pony/base"

	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("box", func() {
	_ = grift.Desc("list", "List all files in baseBox")
	_ = grift.Add("list", func(c *grift.Context) error {
		files := base.ABox.List()
		fmt.Println(files)

		fmt.Println("completed. ")
		return nil
	})
})
