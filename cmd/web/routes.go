package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create a new httprouter instance.
	router := httprouter.New()

	// Wrap the http.NotFound() function in a http.HandlerFunc so that it
	// returns our own custom 404 Not Found response.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	dynamic := alice.New(
		app.sessionManager.LoadAndSave,
		//app.noSurf,
		app.authenticate,
	)

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.signupUserForm))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.signupUserPost))
	router.Handler(http.MethodGet, "/user/signin", dynamic.ThenFunc(app.signinUserForm))
	router.Handler(http.MethodPost, "/user/signin", dynamic.ThenFunc(app.signinUserPost))

	protected := dynamic.Append(app.requireAuthentication)
	protectedAdmin := dynamic.Append(app.requireAuthentication, app.requireAdmin)

	router.Handler(http.MethodGet, "/admin", protectedAdmin.ThenFunc(app.admin))
	router.Handler(http.MethodGet, "/admin/api-keys", protectedAdmin.ThenFunc(app.adminAPIKeys))
	router.Handler(http.MethodPost, "/admin/api-keys/add", protectedAdmin.ThenFunc(app.addAPIKeyPost))
	router.Handler(http.MethodPost, "/admin/api-keys/delete/:name", protectedAdmin.ThenFunc(app.deleteAPIKeyPost))

	router.Handler(http.MethodGet, "/admin/users", protectedAdmin.ThenFunc(app.adminListUsers))
	router.Handler(http.MethodPost, "/admin/users/delete/:id", protectedAdmin.ThenFunc(app.adminDeleteUser))

	router.Handler(http.MethodGet, "/", protected.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/template/create", protected.ThenFunc(app.templateCreateForm))
	router.Handler(http.MethodPost, "/template/create", protected.ThenFunc(app.templateCreatePost))
	router.Handler(http.MethodGet, "/template/edit/:id", protected.ThenFunc(app.templateEditForm))
	router.Handler(http.MethodPost, "/template/edit/:id", protected.ThenFunc(app.templateEditPost))
	router.Handler(http.MethodGet, "/template/view/:id", protected.ThenFunc(app.templateViewForm))
	router.Handler(http.MethodPost, "/template/email/:id", protected.ThenFunc(app.templateEmailPost))
	router.Handler(http.MethodPost, "/template/delete/:id", protected.ThenFunc(app.templateDeletePost))
	router.Handler(http.MethodGet, "/template/loading/:id", protected.ThenFunc(app.showLoading))
	router.Handler(http.MethodGet, "/template/status/:id", protected.ThenFunc(app.checkStatus))
	router.Handler(http.MethodGet, "/template/preview/:id", protected.ThenFunc(app.previewEmail))
	router.Handler(http.MethodGet, "/template/logs/:id", protected.ThenFunc(app.getTemplateLogs))
	router.Handler(http.MethodPost, "/template/title/:id", protected.ThenFunc(app.templateTitlePost))
	router.Handler(http.MethodPost, "/split", protected.ThenFunc(app.split))
	router.Handler(http.MethodPost, "/user/signout", protected.ThenFunc(app.signoutUser))

	// Create a middleware chain containing our 'standard' middleware
	// which will be used by every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
