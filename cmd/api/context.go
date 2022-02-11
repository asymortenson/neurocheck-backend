package main

import (
	"context"
	"net/http"

	"github.com/SevereCloud/vksdk/v2/api"
	"vkparser.com/internal/data"
)

type contextKey string

const (
	vkContextKey   = contextKey("vkAPI")
	userContextKey = contextKey("user")
)

func (app *application) contextSetVK(r *http.Request, vk *api.VK) *http.Request {
	ctx := context.WithValue(r.Context(), vkContextKey, vk)
	return r.WithContext(ctx)
}

func (app *application) contextGetVK(r *http.Request) *api.VK {
	vk, ok := r.Context().Value(vkContextKey).(*api.VK)
	if !ok {
		panic("missing user value in request context")
	}
	return vk
}

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
