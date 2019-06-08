package actions

import "github.com/gobuffalo/buffalo"

// HomeHandler is a default handler to serve up
// a home page.
func homeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html", "layout/empty.html"))
}

func routesHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("routes.html", "layout/empty.html"))
}

func notFoundHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("notFound.html", "layout/empty.html"))
}

func loginHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("login.html", "layout/empty.html"))
}
