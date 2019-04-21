package grifts

import (
	"fmt"
	"pony/base"

	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("db", func() {

	_ = grift.Desc("seed", "Seeds a database")
	_ = grift.Add("seed", func(c *grift.Context) error {
		// Add DB seeding stuff here
		return nil
	})

	_ = grift.Desc("crud", "Generate CRUD SQL for table")
	_ = grift.Add("crud", func(c *grift.Context) error {
		params := c.Args
		if len(params) != 1 {
			fmt.Println("Usage: db:curd <tablename>")
			return nil
		}

		_ = base.PrintSubCode(params[0])

		return nil
	})

})
