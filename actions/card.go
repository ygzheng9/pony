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
