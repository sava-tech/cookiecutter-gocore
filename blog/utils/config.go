package utils

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	DBSource                 string        `mapstructure:"DB_SOURCE"`
	ServerAddress            string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey        string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration      time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration     time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	VerificationCodeDuration time.Duration `mapstructure:"VERIFICATION_CODE_DURATION"`
	RedisAddr                string        `mapstructure:"REDIS_ADDR"`
	RedisHost                string        `mapstructure:"REDISHOST"`
	RedisPort                string        `mapstructure:"REDISPORT"`
	RedisUser                string        `mapstructure:"REDISUSER"`
	RedisPassword            string        `mapstructure:"REDISPASSWORD"`
	ApiAccessKey             string        `mapstructure:"API_ACCESS_KEY"`
	Provider                 string        `mapstructure:"EMAIL_PROVIDER"`
	MailtrapURL              string        `mapstructure:"MAILTRAP_URL"`
	MailtrapAuthToken        string        `mapstructure:"MAILTRAP_AUTH_TOKEN"`
	SendGridURL              string        `mapstructure:"SENDGRID_URL"`
	SendGridAuthToken        string        `mapstructure:"SENDGRID_AUTH_TOKEN"`
	DefaultFromEmail         string        `mapstructure:"DEFAULT_FROM_EMAIL"`
	EmailSubjectPrefit       string        `mapstructure:"EMAIL_SUBJECT_PREFIX"`
	TermiiApiKey             string        `mapstructure:"TERMII_API_KEY"`
	SENDER                   string        `mapstructure:"SENDER"`
	SMSProvider              string        `mapstructure:"SMSPROVIDER"`
	UseMailpit               bool          `mapstructure:"USE_MAILPIT"`
	MailpitHost              string        `mapstructure:"MAILPIT_HOST"`
	MailpitPort              string        `mapstructure:"MAILPIT_PORT"`
	EmailFrom                string        `mapstructure:"EMAIL_FROM"`
	// Twilio Config
	TwilioSID         string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken   string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioPhoneNumber string `mapstructure:"TWILIO_PHONE_NUMBER"`
	TwilioServiceSID  string `mapstructure:"TWILIO_MESSAGING_SERVICE_SID"`
	// PaymentGatewayConfig
	// This is used to configure the payment gateway
	// It can be used to set the payment provider, secret keys, etc.
	TestPayment          bool   `mapstructure:"TEST_PAYMENT"`
	PaymentProvider      string `mapstructure:"PAYMENT_PROVIDER"`
	FlutterWaveSecretKey string `mapstructure:"FLUTTERWAVE_SECRET_KEY"`
	FlutterWaveBaseURL   string `mapstructure:"FLUTTERWAVE_BASE_URL"`
	PaystackSecretKey    string `mapstructure:"PAYSTACK_SECRET_KEY"`
	PaystackBaseURL      string `mapstructure:"PAYSTACK_BASE_URL"`

	// Social Auth (Goth)
	GoogleKey    string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	AppleKey     string `mapstructure:"APPLE_CLIENT_ID"`
	AppleSecret  string `mapstructure:"APPLE_CLIENT_SECRET"`
	SocialCallbackURL string `mapstructure:"SOCIAL_CALLBACK_URL"`
	// Session secret for Goth
	SessionSecret string `mapstructure:"SESSION_SECRET"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AutomaticEnv()

	// 🔑 Required for ENV mapping
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// ======================
	// Core
	// ======================
	viper.BindEnv("DB_SOURCE")
	viper.BindEnv("SERVER_ADDRESS")
	viper.BindEnv("TOKEN_SYMMETRIC_KEY")
	viper.BindEnv("USE_MAILPIT")
	viper.BindEnv("MAILPIT_HOST")
	viper.BindEnv("MAILPIT_PORT")
	viper.BindEnv("EMAIL_FROM")

	// ======================
	// Auth / Tokens
	// ======================
	viper.BindEnv("ACCESS_TOKEN_DURATION")
	viper.BindEnv("REFRESH_TOKEN_DURATION")
	viper.BindEnv("VERIFICATION_CODE_DURATION")

	// ======================
	// SocialAuth / Tokens
	// ======================
	viper.BindEnv("GOOGLE_CLIENT_ID")
	viper.BindEnv("GOOGLE_CLIENT_SECRET")
	viper.BindEnv("APPLE_CLIENT_ID")
	viper.BindEnv("APPLE_CLIENT_SECRET")
	viper.BindEnv("APPLE_TEAM_ID")
	viper.BindEnv("APPLE_KEY_ID")
	viper.BindEnv("APPLE_PRIVATE_KEY")
	viper.BindEnv("SOCIAL_CALLBACK_URL")
	viper.BindEnv("SESSION_SECRET")

	// ======================
	// Redis
	// ======================
	viper.BindEnv("REDIS_ADDR")
	viper.BindEnv("REDISHOST")
	viper.BindEnv("REDISPORT")
	viper.BindEnv("REDISUSER")
	viper.BindEnv("REDISPASSWORD")

	// ======================
	// API / Email
	// ======================
	viper.BindEnv("API_ACCESS_KEY")
	viper.BindEnv("EMAIL_PROVIDER")
	viper.BindEnv("MAILTRAP_URL")
	viper.BindEnv("MAILTRAP_AUTH_TOKEN")
	viper.BindEnv("SENDGRID_URL")
	viper.BindEnv("SENDGRID_AUTH_TOKEN")
	viper.BindEnv("DEFAULT_FROM_EMAIL")
	viper.BindEnv("EMAIL_SUBJECT_PREFIX")

	// ======================
	// SMS
	// ======================
	viper.BindEnv("TERMII_API_KEY")
	viper.BindEnv("SENDER")
	viper.BindEnv("SMSPROVIDER")

	// ======================
	// Twilio
	// ======================
	viper.BindEnv("TWILIO_ACCOUNT_SID")
	viper.BindEnv("TWILIO_AUTH_TOKEN")
	viper.BindEnv("TWILIO_PHONE_NUMBER")
	viper.BindEnv("TWILIO_MESSAGING_SERVICE_SID")

	// ======================
	// Payments
	// ======================
	viper.BindEnv("TEST_PAYMENT")
	viper.BindEnv("PAYMENT_PROVIDER")
	viper.BindEnv("FLUTTERWAVE_SECRET_KEY")
	viper.BindEnv("FLUTTERWAVE_BASE_URL")
	viper.BindEnv("PAYSTACK_SECRET_KEY")
	viper.BindEnv("PAYSTACK_BASE_URL")

	// ======================
	// Unmarshal
	// ======================
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	// ✂️ Sanitize secrets
	config.TokenSymmetricKey = strings.TrimSpace(config.TokenSymmetricKey)

	// 🔎 Debug (optional)
	log.Printf(
		"TOKEN_SYMMETRIC_KEY=%q len=%d",
		config.TokenSymmetricKey,
		len(config.TokenSymmetricKey),
	)

	// ✅ Validate critical secret
	if len(config.TokenSymmetricKey) != 32 {
		return config, fmt.Errorf(
			"TOKEN_SYMMETRIC_KEY must be exactly 32 chars, got %d",
			len(config.TokenSymmetricKey),
		)
	}

	return config, nil
}
