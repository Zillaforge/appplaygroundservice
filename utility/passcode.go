package utility

import (
	"encoding/base32"
	"fmt"
	"math"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
)

func GetPassCode(userAccount, secretKey string) (passCode string, err error) {
	// period 表示每 30 秒更換一次驗證碼
	period := 30
	// secret 表示驗證碼內文。這邊定義格式為：<User ID:secretKey>
	secret := fmt.Sprintf("%s:%s", userAccount, secretKey)

	// 將 secret 作 base32 編碼後儲存在 byteSecret 變數中
	byteSecret := base32.StdEncoding.EncodeToString([]byte(secret))

	// t 為當前時間
	t := time.Now()

	// counter 表示每個 30 秒的累加流水號。舉例來說 162927900 / 30 = 5430930，
	// 其表示從 1970-01-01 00:00:01 開始已經過了幾個 30 秒，每個 30 秒視為一個時間單位，每個時間單位對應一組驗證碼
	counter := uint64(math.Floor(float64(t.Unix()) / float64(period)))

	// Digits 表示驗證碼位數。舉例來說，otp.DigitsSix 表示驗證碼為六位數，其套件還支援八位數驗證碼
	// Algorithm 表示驗證碼產生的演算法。舉例來說，otp.AlgorithmSHA1 表示使用 SHA1 演算法產生驗證碼，其套件還支援 AlgorithmSHA256, AlgorithmSHA512 等等
	code, err := hotp.GenerateCodeCustom(byteSecret, counter, hotp.ValidateOpts{
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", err
	}

	return code, nil
}
