package helpers

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/models"

	"github.com/sirupsen/logrus"
)

// Int64ToString function convert a float number to a string
func Int64ToString(inputNum int64) string {
	return strconv.FormatInt(inputNum, 10)
}

func LogStructAsJSON(data interface{}) interface{} {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Error("LogStructAsJSON Error marshaling JSON: ", err)
		// if there is error in marshalling then atleast return the exact same packet
		return data
	}

	return string(jsonData)
}

func IsAllDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func GenerateOTP(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func SendSms(otp string, mobNo string) {

	var msg91SendOtp models.MSGSendSmsRequest
	msg91SendOtp.FlowID = constants.Msg91FlowID // constants.OTPTEMPLATE
	msg91SendOtp.Sender = constants.OTPSENDER
	msg91SendOtp.ShortURL = "0"
	msg91SendOtp.Mobiles = "91" + mobNo
	msg91SendOtp.Var = otp

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(msg91SendOtp)

	res, err := apihelpers.CallApiMsg91(http.MethodPost, constants.Msg91Url, payload, constants.AuthKeyMsg91)
	if err != nil {
		log.Printf("sendSMS call api error =%v\n", err)
		apihelpers.SendInternalServerError()
		return
	}
	defer res.Body.Close()
}

func GetCurrentTimeInIST() time.Time {
	return time.Now().In(constants.LocationKolkata)
}
