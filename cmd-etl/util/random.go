package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomFloat generates a random float between min and max
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
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

// RandomUser generates a random user name
func RandomUser() string {
	return RandomString(6)
}

// RandomStore generates a random store name
func RandomStore() string {
	return RandomString(12)
}

// RandomMoney generates a random amount of money
func RandomMoney() float64 {
	return RandomFloat(0, 5000.)
}

// RandomDate generates a random timestamp
func RandomTimestamp() time.Time {
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000

	return time.Unix(randomTime, 0).UTC()
}
