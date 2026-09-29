package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const alphabet = "asdfghjklqwertyuiopzxcvbnm"

func init() {
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))
}

// RandomInt generate a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Randon string generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency
func RandomCurrency() string {
	currencies := []string{"USD", "NGN", "GHS"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomCountry generates a random country
func RandomCountry() string {
	countries := []string{"Nigeria", "Ghana", "India", "South Africa", "USA"}
	n := len(countries)
	return countries[rand.Intn(n)]
}

func RandomGender() string {
	gender := []string{"male", "female", "prefer not to say"}
	n := len(gender)
	return gender[rand.Intn(n)]
}

// random email addresss
func RandomEmail() string {
	// random email provider
	email_providers := []string{"gmail", "yahoo", "private"}
	e := len(email_providers)
	provider := email_providers[rand.Intn(e)]

	// generate random string
	ran_str := RandomString(8)

	return fmt.Sprintf("%s@%s.com", ran_str, provider)
}

func RandomNewUUID() pgtype.UUID {
    u := uuid.New()
    return pgtype.UUID{
        Bytes: u,   // uuid.UUID is a [16]byte, matches directly
        Valid: true,
    }
}
