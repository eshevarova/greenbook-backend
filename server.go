package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"net/http"
)

var (
	corsAllowHeaders     = "*"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func createRequestHandler() fasthttp.RequestHandler {
	metric := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())

	return func(ctx *fasthttp.RequestCtx) {
		path, method := ctx.Path(), ctx.Method()
		switch {
		case string(ctx.Method()) == "OPTIONS":
			writeCors(ctx)
			ctx.SetStatusCode(fasthttp.StatusOK)

		case string(path) == "/v1/private":
			writeCors(ctx)
			processGetPrivatePage(ctx)

		case string(path) == "/v1/common":
			writeCors(ctx)
			processGetCommonPage(ctx)

		case string(path) == "/v1/friends":
			writeCors(ctx)
			processGetFriends(ctx)

		case string(path) == "/v1/login" && string(method) == "POST":
			writeCors(ctx)
			processLogin(ctx)

		case string(path) == "/v1/registration" && string(method) == "POST":
			writeCors(ctx)
			processRegistration(ctx)

		case string(path) == "/v1/friends/add" && string(method) == "POST":
			writeCors(ctx)
			processAddFriend(ctx)

		case string(path) == "/v1/friends/remove" && string(method) == "POST":
			writeCors(ctx)
			processRemoveFriend(ctx)

		case string(path) == "/health":
			ctx.SetStatusCode(http.StatusOK)

		case string(path) == "/metrics":
			metric(ctx)

		default:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
	}
}


func writeCors(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
	ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
	ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
}

func writeError(ctx *fasthttp.RequestCtx, msg string, statusCode int) {
	ctx.SetBodyString(msg)
	ctx.SetStatusCode(statusCode)
}
