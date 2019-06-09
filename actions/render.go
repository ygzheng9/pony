package actions

import (
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

func Init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			// uncomment for non-Bootstrap form helpers:
			// "form":     plush.FormHelper,
			// "form_for": plush.FormForHelper,
			"isCurrentPathName": func(current buffalo.RouteInfo, name string) string {
				// 子菜单选中时，父菜单也得选中
				parts := strings.Split(name, " ")

				for _, s := range parts {
					f := "tabler" + s + "Path"

					// fmt.Printf("a=%s, b=%s;", current.PathName, f)

					if strings.EqualFold(current.PathName, f) {
						return "active"
					}
				}

				return ""
			},
		},
	})
}
