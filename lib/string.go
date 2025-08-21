package lib

import (
	"adire-apparel/config"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

func CreateSku(prefix string) string {
	randomInt := int(time.Now().UnixNano()%900000) + 100000
	return prefix + string(rune(randomInt))
}

func ToKebabCase(str string) string {
	return strings.ReplaceAll(strings.ToLower(str), " ", "-")
}

func GenerateOtp() string {
	const charset = "0123456789"
	otp := make([]byte, 6)
	for i := range otp {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		otp[i] = charset[randomIndex.Int64()]
	}
	return string(otp)
}

func GenerateUrl() string {
	clientUrl := config.AppConfig.ClientUrl
	otp := GenerateOtp()
	url := clientUrl + "/reset-password/" + otp
	return url
}

func GenerateTransactionReference() string {
	timestamp := time.Now().Unix()
	randomStr := GenerateOtp()
	return fmt.Sprintf("TXN-%d%s", timestamp, randomStr)
}

func GenerateTransactionReferenceWithType(txType string) string {
	timestamp := time.Now().Unix()
	randomStr := GenerateOtp()
	return fmt.Sprintf("TXN-%s-%d%s", txType, timestamp, randomStr)
}

func GenerateTrackingNumber() string {
	timestamp := time.Now().Unix()
	randomStr := GenerateOtp()
	return fmt.Sprintf("TRN-%d%s", timestamp, randomStr)
}
