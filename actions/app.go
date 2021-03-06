package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	csrf "github.com/gobuffalo/mw-csrf"
	forcessl "github.com/gobuffalo/mw-forcessl"
	i18n "github.com/gobuffalo/mw-i18n"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/packr/v2"
	"github.com/unrolled/secure"

	// _ "github.com/konsorten/go-windows-terminal-sequences"

	"pony/models"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// T for i18
var T *i18n.Translator

type H map[string]interface{}

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		Init()

		models.Init()

		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_pony_session",
		})

		// Automatically redirect to SSL
		// app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Setup and use translations:
		app.Use(translations())

		app.GET("/", homeHandler)
		app.GET("/login", loginHandler)

		// all routes
		app.GET("/routes", routesHandler)

		// biz logic
		app.GET("/surveys/open", SurveysOpen)
		app.POST("/surveys/submit", SurveysSubmit)

		app.GET("/matrix/open", MatrixOpen)
		app.POST("/matrix/open", MatrixSubmit)
		app.GET("/matrix/openTypeA", MatrixOpenTypeA)
		app.POST("/matrix/openTypeA", MatrixSubmitTypeA)

		app.GET("/chart/first", ChartFirst)
		app.GET("/chart/get_wordcloud", WordCloudHandle)
		app.GET("/chart/get_wordfreq", WordFreqHandle)
		app.GET("/chart/get_worddist", wordDistHandle)

		// 决策论
		app.GET("/games/index", gamesIndex)

		// app.GET("/games/byid", gamesByID)
		app.POST("/games/byid", gamesByID)

		app.POST("/games/create", gamesCreate)

		app.POST("/games/saveCriterion", gamesSaveCriterion)
		app.POST("/games/saveCriterionPairs", gamesSaveCriterionPairs)

		app.POST("/games/saveOptions", gamesSaveOptions)
		app.POST("/games/loadOptionPairs", gamesLoadOptionPairs)
		app.POST("/games/saveOptionPairs", gamesSaveOptionPairs)

		app.POST("/games/calcFinal", gamesCalcFinal)

		// 两个维度的散点图
		app.GET("/scatter/show", scatterShow)
		app.GET("/scatter/data/{name}", scatterData)

		// tailwind css
		app.GET("/card/tailwind", tailwindDemo)
		app.GET("/card/tailwind_admin", tailwindAdmin)
		app.GET("/card/tailwind_admin_day", tailwindAdminDay)
		app.GET("/card/tailwind_profile", tailwindProfile)
		app.GET("/card/tailwind_landing", tailwindLanding)
		app.GET("/card/tailwind_form", tailwindForm)

		// bootstrap css
		app.GET("/card/bootstrap", bootstrapDemo)

		// vali admin
		app.GET("/vali/index", valiHandler)
		app.GET("/vali/charts", ChartsHandler)
		app.GET("/vali/bootstrap-components", UIElementHandler)

		// tabler
		g := app.Group("/tabler")
		g.GET("/index", tablerIndex)
		g.GET("/cards", tablerCard)
		g.GET("/pricing-cards", tablerPricingCard)
		g.GET("/charts", tablerCharts)
		g.GET("/profile", tablerProfile)
		g.GET("/empty", tablerEmpty)
		g.GET("/email", tablerEmail)
		g.GET("/form-elements", tablerFormElements)
		g.GET("/store", tablerStore)
		g.GET("/blog", tablerBlog)

		// redomHandler
		app.GET("/fe/redom", redomHandler)
		app.GET("/fe/react", reactHandler)
		app.GET("/fe/preact", preactHandler)
		app.GET("/fe/inferno", infernoHandler)

		// 在生产环境里，页面不存在时，重定向到统一的页面
		app.GET("/notFound", notFoundHandler)
		// if ENV == "production" {
		app.ErrorHandlers[404] = pageNotFound
		// }

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.New("app:locales", "../locales"), "en-US"); err != nil {
		_ = app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

func pageNotFound(status int, err error, c buffalo.Context) error {
	// res := c.Response()
	// res.WriteHeader(404)
	// res.Write([]byte(fmt.Sprintf("Oops!! There was an error %s", err.Error())))

	// c.Redirect(307, "notFoundPath()")
	// c.Render(200, r.HTML("notFound.html", "layout/empty.html"))
	tmpl := `
	<!DOCTYPE html>
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta charset="utf-8">
    <title>Survey</title>

		<link rel="stylesheet" href="/assets/vendors/main.css">

		<link rel="icon" href="/assets/images/favicon.ico">

  </head>

  <body class="app">
	<main>
		<div class="page-error tile">
			<h1><i class="fa fa-exclamation-circle"></i> 信息不存在 </h1>
			<p>The page you have requested is not found.</p>
			<p><a class="btn btn-primary" href="javascript:window.history.back();">返回</a></p>
		</div>
	</main>
  </body>
</html>
	`
	res := c.Response()
	res.WriteHeader(404)
	_, _ = res.Write([]byte(tmpl))

	return nil
}
