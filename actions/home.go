package actions

import "github.com/gobuffalo/buffalo"

// HomeHandler is a default handler to serve up
// a home page.
func homeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html", "layout/empty.html"))
}

// HomeHandler is a default handler to serve up
// a home page.
func valiHandler(c buffalo.Context) error {
	c.Set("task-title", "Dashboard")
	c.Set("task-icon", "fa fa-dashboard")
	c.Set("task-desc", "A free and open source Bootstrap 4 admin template")

	return c.Render(200, r.HTML("vali/dashboard.html"))
}

// ChartsHandler testing dynamic content handling of vali-admin theme
func ChartsHandler(c buffalo.Context) error {
	c.Set("task-title", "Charts")
	c.Set("task-icon", "fa fa-pie-chart")
	c.Set("task-desc", "Various type of charts for your project")

	return c.Render(200, r.HTML("vali/charts.html"))
}

// UIElementHandler bootstrap UI elements demo
func UIElementHandler(c buffalo.Context) error {
	c.Set("task-title", "Bootstrap Elements")
	c.Set("task-icon", "fa fa-laptop")
	c.Set("task-desc", "Bootstrap Components")

	return c.Render(200, r.HTML("vali/bootstrap-components.html"))
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
