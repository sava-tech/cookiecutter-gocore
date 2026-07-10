package infrastructure

import (
	"fmt"
	"log"
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/google"
)

func providerCallbackURL(baseURL, provider string) string {
	base := strings.TrimSuffix(baseURL, "/")
	suffix := fmt.Sprintf("/%s/callback", provider)
	if strings.HasSuffix(base, suffix) {
		return base
	}
	return base + suffix
}

// SetupGoth initializes all Goth providers
func SetupGoth(config GothProviderConfig) {

	// Google provider
	if config.GoogleKey != "" && config.GoogleSecret != "" {
		googleCallback := providerCallbackURL(config.CallbackURL, "google")
		log.Printf("Registering Google with callback: %s", googleCallback)

		provider := google.New(
			config.GoogleKey,
			config.GoogleSecret,
			googleCallback,
			"email", "profile",
		)
		goth.UseProviders(provider)
		log.Println("Google provider registered")
	} else {
		log.Printf("Google provider NOT configured (missing credentials)")
	}

	// Apple provider
	if config.AppleKey != "" && config.AppleSecret != "" {
		appleCallback := providerCallbackURL(config.CallbackURL, "apple")
		log.Printf("Registering Apple with callback: %s", appleCallback)

		provider := apple.New(
			config.AppleKey,
			config.AppleSecret,
			appleCallback,
			nil,
			apple.ScopeName, apple.ScopeEmail,
		)
		goth.UseProviders(provider)
		log.Println("Apple provider registered")
	} else {
		log.Printf("Apple provider not configured (optional)")
	}

	// Debug: List all registered providers
	log.Println("Registered providers:")
	// Goth doesn't have a built-in way to list providers, but we can check
	_, err := goth.GetProvider("google")
	if err != nil {
		fmt.Println("the provider is not working")
	}

	_, err = goth.GetProvider("apple")
	if err != nil {
		fmt.Println("the provider is not working")
	}
}


