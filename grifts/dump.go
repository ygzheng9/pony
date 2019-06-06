package grifts

import (
	"fmt"
	"time"

	"github.com/markbates/grift/grift"

	"pony/base"
)

var _ = grift.Namespace("dump", func() {

	_ = grift.Desc("po", "Import PO items ")
	_ = grift.Add("po", func(c *grift.Context) error {
		// Add DB seeding stuff here
		base.ProcessPOItems()

		return nil
	})

	_ = grift.Desc("mo", "Import MO items ")
	_ = grift.Add("mo", func(c *grift.Context) error {
		// Add DB seeding stuff here
		base.ProcessMOItems()

		return nil
	})

	_ = grift.Desc("inv", "Import inv items ")
	_ = grift.Add("inv", func(c *grift.Context) error {
		// Add DB seeding stuff here
		base.ProcessInvItems()

		return nil
	})

	_ = grift.Desc("so", "Import so items ")
	_ = grift.Add("so", func(c *grift.Context) error {
		// Add DB seeding stuff here
		base.ProcessSoItems()

		return nil
	})

	_ = grift.Add("test", func(c *grift.Context) error {
		s := "4/26/18 08:34"
		t, err := time.Parse("1/2/06 15:04", s)

		// 43471
		if err != nil {
			fmt.Printf("error: %+v\n", err)
		} else {
			fmt.Printf("ok: %+v\n", t)
		}
		return nil
	})
})
