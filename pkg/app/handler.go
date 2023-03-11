package app

import (
	"github.com/gocopper/copper/csql"
	"net/http"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clogger"
)

type NewHTTPHandlerParams struct {
	DatabaseTxMW    *csql.TxMiddleware
	RequestLoggerMW *chttp.RequestLoggerMiddleware

	App  *Router
	HTML *chttp.HTMLRouter

	Logger clogger.Logger
}

func NewHTTPHandler(p NewHTTPHandlerParams) http.Handler {
	return chttp.NewHandler(chttp.NewHandlerParams{
		GlobalMiddlewares: []chttp.Middleware{

			p.DatabaseTxMW,
			p.RequestLoggerMW,
		},

		Routers: []chttp.Router{
			p.HTML,

			p.App,
		},

		Logger: p.Logger,
	})
}
