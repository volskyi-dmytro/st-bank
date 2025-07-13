package util

import (
	"math/rand"
	"strings"
	"time"
)

// alphabet contains all lowercase letters for generating random strings
const alphabet = "abcdefghijklmnopqrstuvwxyz"

// init initializes the random seed with current time
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
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

// RandomMoney generates a random money amount
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency returns a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "UAH"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomEmail generates a random email address
func RandomEmail() string {
	return RandomString(6) + "@" + RandomString(4) + ".com"
}
