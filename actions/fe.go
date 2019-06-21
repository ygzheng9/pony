package actions

import "github.com/gobuffalo/buffalo"

// feRender for all files under fe dir
func feRender(c buffalo.Context, page string) error {
	const layout = "fe/layout.html"
	content := "fe/" + page

	return c.Render(200, r.HTML(content, layout))
}

func redomHandler(c buffalo.Context) error {
	return feRender(c, "redom.html")
}

func reactHandler(c buffalo.Context) error {
	return feRender(c, "react.html")
}

func preactHandler(c buffalo.Context) error {
	return feRender(c, "preact.html")
}

func infernoHandler(c buffalo.Context) error {
	return feRender(c, "inferno.html")
}
