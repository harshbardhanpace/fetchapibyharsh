package apihelpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"space/constants"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceId, userAgent, remoteAddr, authToken string) (*http.Response, error) {
	return CallApi(methodType, url, payload, deviceType, deviceId, userAgent, remoteAddr, authToken)
}

var CallAPIFuncV2 = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceId, userAgent, remoteAddr, authToken string) (*http.Response, error) {
	return CallApiV2(methodType, url, payload, deviceType, deviceId, userAgent, remoteAddr, authToken)
}

// ResponseData structure
type ResponseData struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

type ReturnValue struct {
	StatusCode int    `json:"statuscode"`
	ApiRes     APIRes `json:"apires"`
}

type APIRes struct {
	Status    bool        `json:"status"`
	Message   string      `json:"message"`
	ErrorCode string      `json:"errorcode"`
	Data      interface{} `json:"data"`
}

func logJSON(data interface{}) interface{} {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Error("LogJSON Error marshaling JSON: ", err)
		// if there is error in marshalling then atleast return the exact same packet
		return data
	}

	return string(jsonData)
}

// Message returns map data
func Message(status int, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func ErrorMessage(c *gin.Context, ErrorCode string, optionalParams ...interface{}) {
	var apiRes APIRes
	apiRes.ErrorCode = ErrorCode
	apiRes.Message = constants.ErrorCodeMap[ErrorCode]
	apiRes.Status = false
	CustomResponse(c, http.StatusBadRequest, apiRes, optionalParams...)
}

func SendErrorController(c *gin.Context, status bool, errorCode string, httpCode int, optionalParams ...interface{}) {
	var apiRes APIRes
	apiRes.Status = status
	apiRes.Message = constants.ErrorCodeMap[errorCode]
	apiRes.ErrorCode = errorCode
	CustomResponse(c, httpCode, apiRes, optionalParams...)
}

// Respond returns basic response structure
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func CustomResponse(c *gin.Context, code int, data interface{}, optionalParams ...interface{}) {
	// Format optional parameters into a single string
	var optionalParamsStr string
	if len(optionalParams) > 0 {
		optionalParamsStr = fmt.Sprintf(" optionalParams: %v", optionalParams)
	}

	logrus.Info("CustomResponse ", optionalParamsStr, " code: ", code, " data: ", logJSON(data))

	// Send the JSON response
	c.JSON(code, data)
}

func CallApi(methodType, url string, payload *bytes.Buffer, deviceType, deviceId, userAgent, remoteAddr, authToken string) (*http.Response, error) {
	req, _ := http.NewRequest(methodType, url, payload)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-device-type", deviceType)
	req.Header.Set("x-device-id", deviceId)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("remote-address", remoteAddr)
	req.Header.Set("Authorization", authToken)
	if authToken != "" && len(authToken) >= 8 {
		req.Header.Set("x-authorization-token", authToken[7:])
	}

	if authToken != "" && len(authToken) >= 8 {
		req.Header.Set("x-authorization-token", authToken[7:])
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func CallApiV2(methodType, url string, payload *bytes.Buffer, deviceType, deviceId, userAgent, remoteAddr, authToken string) (*http.Response, error) {
	req, _ := http.NewRequest(methodType, url, payload)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-device-type", deviceType)
	req.Header.Set("x-device-id", deviceId)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("remote-address", remoteAddr)
	req.Header.Set("Authorization", authToken)
	if authToken != "" && len(authToken) >= 8 {
		req.Header.Set("x-authorization-token", authToken[7:])
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func CallCmotsApi(methodType, url string, payload *bytes.Buffer, authToken string) (*http.Response, error) {
	req, _ := http.NewRequest(methodType, url, payload)
	req.Header.Set("Authorization", authToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func CallApiShilpi(methodType, url string, payload *bytes.Buffer) (*http.Response, error) {
	req, _ := http.NewRequest(methodType, url, payload)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func SendInternalServerError() (int, APIRes) {
	var apiRes APIRes
	apiRes.Status = false
	apiRes.Message = constants.ErrorCodeMap[constants.InternalServerError]
	apiRes.ErrorCode = constants.InternalServerError
	return http.StatusInternalServerError, apiRes
}

func SendErrorResponse(status bool, errorCode string, httpCode int) (int, APIRes) {
	var apiRes APIRes
	apiRes.Status = status
	apiRes.Message = constants.ErrorCodeMap[errorCode]
	apiRes.ErrorCode = errorCode
	return httpCode, apiRes
}

func CallApiTradeLab(methodType, url string, payload *bytes.Buffer, authToken string) (*http.Response, error) {
	req, _ := http.NewRequest(methodType, url, payload)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-device-type", "web")
	req.Header.Set("User-Agent", "")
	req.Header.Set("remote-address", "")
	//req.Header.Set("Authorization", authToken)
	if authToken != "" && len(authToken) >= 8 {
		req.Header.Set("x-authorization-token", authToken[7:])
	}

	//req.Header.Set("Authorization", authToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func CallApiFinvu(methodType, url string, payload *bytes.Buffer, authToken string) (*http.Response, error) {
	req, _ := http.NewRequest(methodType, url, payload)
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func CallApiFinvuBankStatementPdf(methodType, url string, payload *bytes.Buffer, authToken string) (*http.Response, error) {
	req, _ := http.NewRequest(methodType, url, payload)
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}
	req.Header.Set("Accept", "application/pdf")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func CallApiOauth(methodType, url string, payload *bytes.Buffer, Authorization string) (*http.Response, error) {
	req, _ := http.NewRequest(methodType, url, payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", Authorization)
	req.Header.Set("CacheControl", "no-cache")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

var CallAPIFuncOauth = func(methodType, url string, payload *bytes.Buffer, Authorization string) (*http.Response, error) {
	return CallApiOauth(methodType, url, payload, Authorization)
}

func CallApiFreshDesk(methodType, url string, payload *bytes.Buffer, apiKey, pass string) (*http.Response, error) {
	req, _ := http.NewRequest(methodType, url, payload)

	req.SetBasicAuth(apiKey, pass)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func CallApiMsg91(methodType, url string, payload *bytes.Buffer, authKey string) (*http.Response, error) {
	req, _ := http.NewRequest(methodType, url, payload)
	req.Header.Set("authKey", authKey)
	req.Header.Set("content-type", "application/JSON")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}
