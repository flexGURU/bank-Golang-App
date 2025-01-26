package utils

import (
	"math/rand"
	"strings"
	"time"
)


const alphabets = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}


func RandomInt(min, max int32) int32 {

	return min + rand.Int31n(max-min+1)
	
}

func RandomString(n int) string {

	var sb strings.Builder

	k := len(alphabets)

	for i := 0; i < n; i ++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()

}

func RandomOwner() string {
	
	return RandomString(6)

}

func RandomMoney() int32 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {

	currencies := []string{"KES", "USD","EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]

}

