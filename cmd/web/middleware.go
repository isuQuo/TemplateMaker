package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement correct CSP headers
		//w.Header().Set("Content-Security-Policy",
		//	"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Conent-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response before calling
				// serverError. This tells the client to close the connection when
				// it receives the response. This is good practice when dealing
				// with panics, because Go's HTTP server doesn't automatically close
				// connections if a handler panics.
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
			return
		}

		// Add Cache-Control header to all authenticated pages.
		// This will prevent the browser from caching any pages that require
		// authentication.
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (app *application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	// TODO - change this to true in production
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.sessionManager.GetString(r.Context(), "authenticatedUserID")
		if id == "" {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.users.Get(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if user != nil {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			ctx = context.WithValue(ctx, isAdministratorContextKey, user.IsAdmin)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAdmin(r) {
			// Redirect to forbidden or show error
			app.clientError(w, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
