package main

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
const isAdministratorContextKey = contextKey("isAdministrator")
