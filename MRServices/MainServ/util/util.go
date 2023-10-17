package util

import (
	"math"
	"math/rand"
	"regexp"
	"strconv"

	com "github.com/aureontu/MRWebServer/mr_services/common"
)

func CheckEmailAddr(emailAddr string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$`)
	return regex.MatchString(emailAddr)
}

func GenerateRandomCode(length int) string {
	codeInt := rand.Intn(int(math.Pow(float64(10), float64(length))))
	code := strconv.Itoa(codeInt)
	for len(code) < com.VCodeLen {
		code = "0" + code
	}
	return code
}
