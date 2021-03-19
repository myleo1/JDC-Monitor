package cryptokit

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cast"
	"time"
)

var hmacKey = "706390cef611241d57573ca601eb3c061e174948"

func EncodeAuthorization(body, accessKey string) string {
	deviceKey := encodeMD5(fmt.Sprintf("ios6.5.5iPhone11,813.7:%s", cast.ToString(getDayOfYear())))
	t := getTime()
	return fmt.Sprintf("smart %s:::%s:::%s", accessKey, encodeBase64(encodeHmacSHA1(hmacKey, fmt.Sprintf("%spostjson_body%s%s%s%s", deviceKey, body, t, accessKey, deviceKey))), t)
}

func encodeMD5(deviceKey string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(deviceKey)))
}

func encodeHmacSHA1(key, data string) []byte {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func getTime() string {
	return time.Now().Format("2006-01-02T15:04:05.000Z")
}

func getDayOfYear() int {
	t := time.Now()
	return cast.ToInt(t.Sub(time.Date(t.Year(), 01, 01, 0, 0, 0, 0, time.Local)).Hours()/24) + 1
}
