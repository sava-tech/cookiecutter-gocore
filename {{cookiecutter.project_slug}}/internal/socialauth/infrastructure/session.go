package infrastructure

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

type GothProviderConfig struct {
	GoogleKey    string
	GoogleSecret string
	AppleKey     string
	AppleSecret  string
	CallbackURL  string
}

// SetupSession initializes the session store for Goth
func SetupSession(secret string) {
	// Create a cookie store with the secret
	store := sessions.NewCookieStore([]byte(secret))

	// Configure the store
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,     // Set to true in production with HTTPS
		MaxAge:   86400 * 7, // 7 days
	}

	// Set the store for gothic
	gothic.Store = store
}
