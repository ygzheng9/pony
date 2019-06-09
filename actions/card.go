package actions

import "github.com/gobuffalo/buffalo"

func tailwindDemo(c buffalo.Context) error {
	return c.Render(200, r.HTML("card/tailwind_index.html", "card/tailwind_layout.html"))
}

func tailwindAdmin(c buffalo.Context) error {
	return c.Render(200, r.HTML("card/tailwind_admin.html", "card/tailwind_layout.html"))
}

func tailwindAdminDay(c buffalo.Context) error {
	return c.Render(200, r.HTML("card/tailwind_admin_day.html", "card/tailwind_layout.html"))
}

func tailwindProfile(c buffalo.Context) error {
	return c.Render(200, r.HTML("card/tailwind_profile.html", "card/tailwind_layout.html"))
}

func tailwindLanding(c buffalo.Context) error {
	return c.Render(200, r.HTML("card/tailwind_landing.html", "card/tailwind_layout.html"))
}

func tailwindForm(c buffalo.Context) error {
	return c.Render(200, r.HTML("card/tailwind_form.html", "card/tailwind_layout.html"))
}

func bootstrapDemo(c buffalo.Context) error {
	return c.Render(200, r.HTML("card/bootstrap_index.html", "card/bootstrap_layout.html"))
}

//////////////////////////////////////////////

// valiRender for all files under vali dir
func valiRender(c buffalo.Context, page string) error {
	const valiTemplate = "vali/layout/main.html"
	fullPage := "vali/" + page

	return c.Render(200, r.HTML(fullPage, valiTemplate))
}

// a home page.
func valiHandler(c buffalo.Context) error {
	c.Set("task-title", "Dashboard")
	c.Set("task-icon", "fa fa-dashboard")
	c.Set("task-desc", "A free and open source Bootstrap 4 admin template")

	return valiRender(c, "dashboard.html")
}

// ChartsHandler testing dynamic content handling of vali-admin theme
func ChartsHandler(c buffalo.Context) error {
	c.Set("task-title", "Charts")
	c.Set("task-icon", "fa fa-pie-chart")
	c.Set("task-desc", "Various type of charts for your project")

	return valiRender(c, "charts.html")
}

// UIElementHandler bootstrap UI elements demo
func UIElementHandler(c buffalo.Context) error {
	c.Set("task-title", "Bootstrap Elements")
	c.Set("task-icon", "fa fa-laptop")
	c.Set("task-desc", "Bootstrap Components")

	return valiRender(c, "bootstrap-components.html")
}

//////////////////////////////////////////////
func tablerRender(c buffalo.Context, page string) error {
	const valiTemplate = "tabler/layout/main.html"
	fullPage := "tabler/" + page

	return c.Render(200, r.HTML(fullPage, valiTemplate))
}

func tablerIndex(c buffalo.Context) error {
	return tablerRender(c, "index.html")
}

func tablerCard(c buffalo.Context) error {
	return tablerRender(c, "cards.html")
}

func tablerPricingCard(c buffalo.Context) error {
	return tablerRender(c, "pricing-cards.html")
}

func tablerCharts(c buffalo.Context) error {
	return tablerRender(c, "charts.html")
}

func tablerProfile(c buffalo.Context) error {
	return tablerRender(c, "profile.html")
}

func tablerEmpty(c buffalo.Context) error {
	return tablerRender(c, "empty.html")
}

func tablerEmail(c buffalo.Context) error {
	return tablerRender(c, "email.html")
}

func tablerFormElements(c buffalo.Context) error {
	return tablerRender(c, "form-elements.html")
}

func tablerStore(c buffalo.Context) error {
	return tablerRender(c, "store.html")
}

func tablerBlog(c buffalo.Context) error {
	return tablerRender(c, "blog.html")
}
