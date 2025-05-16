package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	apihelpers "space/apiHelpers"
	"space/business/tradelab"
	"space/constants"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				loggerconfig.Error("Alert Severity:P0-Critical, PanicRecover RECOVERED from: %v\n", r)
				stackTrace := make([]byte, 1024)
				runtime.Stack(stackTrace, true)
				loggerconfig.Error("Alert Severity:P0-Critical, PanicRecover Stack trace: %s\n", string(stackTrace))

				code, apiRes := apihelpers.SendInternalServerError()
				apihelpers.CustomResponse(c, code, apiRes)

				// Abort further request processing
				c.Abort()
			}
		}()
		c.Next()
	}
}

/*
UserMiddlewares function to add auth
*/
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Code for middlewares
		var reqH models.ReqHeader
		if err := c.ShouldBindHeader(&reqH); err != nil {
			loggerconfig.Warn("invalid  headers:", err)
		}

		reqH.RequestId = uuid.New().String()

		if reqH.ClientPublicIP == "" {
			loggerconfig.Info("ClientPublicIP is empty, RequestId: ", reqH.RequestId)
		}

		if reqH.DeviceId == "" {
			loggerconfig.Info("DeviceId is empty, RequestId: ", reqH.RequestId)
		}

		if reqH.ClientType != "" && strings.ToLower(reqH.ClientType) == constants.GUESTUSERTYPE {
			authDaoValidate := GuestDaoMiddleware(reqH)
			if !authDaoValidate.Status {
				c.JSON(http.StatusForbidden, authDaoValidate)
				c.Abort()
			}
		}

		c.Set("reqH", reqH)
		c.Next()
	}
}

func GuestDaoMiddleware(reqH models.ReqHeader) apihelpers.APIRes {
	//Code for middlewares

	var resJS apihelpers.APIRes

	if reqH.ClientId == "" {
		resJS.Status = false
		resJS.Message = constants.ErrorCodeMap[constants.InvalidClient]
		resJS.ErrorCode = constants.InvalidClient
		return resJS
	}

	if len(reqH.Authorization) <= 7 {
		resJS.Status = false
		resJS.Message = constants.ErrorCodeMap[constants.TokenMissing]
		resJS.ErrorCode = constants.TokenMissing
		return resJS
	}

	sub, err := helpers.ValidateToken(reqH.Authorization[7:])
	if err != nil {
		loggerconfig.Error("GuestDaoMiddleware error validating auth token, err:", err)
		resJS.Status = false
		resJS.Message = constants.ErrorCodeMap[constants.InvalidToken]
		resJS.ErrorCode = constants.InvalidToken
		return resJS
	}

	//check for validity in redis
	redisCli := cache.GetRedisClientObj()

	if reqH.ClientId != sub {
		loggerconfig.Error("GuestDaoMiddleware error validating userId using guest auth token, clientId:", reqH.ClientId, " sub:", sub)
		resJS.Status = false
		resJS.Message = constants.ErrorCodeMap[constants.MismatchAuthClient]
		resJS.ErrorCode = constants.MismatchAuthClient
		return resJS
	}

	keyExists, _ := redisCli.Exists("auth|" + sub).Result()
	if keyExists != 1 {
		loggerconfig.Error("GuestDaoMiddleware error in finding token in redis keyExists != 1, keyExists", keyExists)
		resJS.Status = false
		resJS.Message = constants.ErrorCodeMap[constants.InvalidToken]
		resJS.ErrorCode = constants.InvalidToken
		return resJS
	}

	resJS.Status = true
	return resJS
}

func TLAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Code for middlewares
		var reqH models.ReqHeader
		if err := c.ShouldBindHeader(&reqH); err != nil {
			loggerconfig.Warn("invalid  headers=", err)
		}
		reqH.RequestId = uuid.New().String()

		var resJS apihelpers.APIRes
		if len(reqH.Authorization) <= 7 {
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.TokenMissing]
			resJS.ErrorCode = constants.TokenMissing
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}

		redisCli := cache.GetRedisClientObj()

		resRedis := redisCli.GetRedis(reqH.ClientId)
		decrypt, err := resRedis.Result()
		if decrypt != "" && decrypt == reqH.Authorization {
			loggerconfig.Info("Auth passed for client:", reqH.ClientId, " by cache.")
			c.Set("reqH", reqH)
			c.Next()
			return
		}

		url := constants.TLURL + tradelab.FETCHFUNDSURL + "?type=" + constants.FetchFundsTypeTL + "&client_id=" + url.QueryEscape(reqH.ClientId)
		payload := new(bytes.Buffer)
		//call api
		res, err := apihelpers.CallApiTradeLab(http.MethodGet, url, payload, reqH.Authorization)
		if err != nil {
			loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " fetchFundsRes call api error =", err, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.AuthenticationFailed]
			resJS.ErrorCode = constants.AuthenticationFailed
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		tlErrorRes := models.Auth{}
		err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
		tlAuthRes := models.Auth{}
		json.Unmarshal([]byte(string(body)), &tlAuthRes)
		if tlAuthRes.Status != "success" {
			loggerconfig.Error("Alert Severity:P0-Critical, TLAuth (controller) tl response != success fetchFundsRes call api error =", err, " tl response status: ", tlAuthRes, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.AuthenticationFailed]
			resJS.ErrorCode = constants.AuthenticationFailed
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}

		if res.StatusCode != http.StatusOK {
			loggerconfig.Error("TLAuth tl status not ok StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		}

		cacheTime := time.Duration(constants.TokenCacheTime) * time.Minute

		redisCli.SetRedis(reqH.ClientId, reqH.Authorization, cacheTime)
		loggerconfig.Info("Token set in cache for clientid:", reqH.ClientId)

		c.Set("reqH", reqH)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Code for middlewares

		var resJS apihelpers.APIRes
		var reqH models.ReqHeader
		if err := c.ShouldBindHeader(&reqH); err != nil {
			loggerconfig.Warn("invalid  headers=", err)
		}

		if len(reqH.Authorization) <= 7 {
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.TokenMissing]
			resJS.ErrorCode = constants.TokenMissing
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}

		sub, err := helpers.ValidateToken(reqH.Authorization[7:])
		fmt.Printf("error validating auth token =%v\n", err)
		if err != nil {
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.InvalidToken]
			resJS.ErrorCode = constants.InvalidToken
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}

		reqH.RequestId = uuid.New().String()

		//check for token expiry
		// val := cache.Exists(reqH.Authorization)
		// if val.Val() != 0 {
		// 	fmt.Println("I am here")
		// 	resJS.Status = false
		// 	resJS.Message = constants.ErrorCodeMap[constants.InvalidToken]
		// 	resJS.ErrorCode = constants.InvalidToken
		// 	c.JSON(http.StatusForbidden, resJS)
		// 	c.Abort()
		// 	return
		// }

		fmt.Printf("subject=%v\n", sub)
		c.Set("reqH", reqH)
		c.Next()
	}
}

func TLAuthReports() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Code for middlewares
		var reqH models.ReqHeader
		if err := c.ShouldBindHeader(&reqH); err != nil {
			loggerconfig.Warn("invalid  headers=", err)
		}
		reqH.RequestId = uuid.New().String()

		var resJS apihelpers.APIRes
		if len(reqH.Authorization) <= 7 {
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.TokenMissing]
			resJS.ErrorCode = constants.TokenMissing
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}

		_, err := helpers.ExtractTokenHeader(reqH.Authorization[7:])
		if err != nil {
			// Invalid authtoken
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.AuthenticationFailed]
			resJS.ErrorCode = constants.AuthenticationFailed
			c.JSON(http.StatusUnauthorized, resJS)
			c.Abort()
			return
		}

		url := constants.TLURL + constants.FetchProfileUrlTL + "?client_id=" + url.QueryEscape(reqH.ClientId)
		payload := new(bytes.Buffer)
		//call api
		res, err := apihelpers.CallApiTradeLab(http.MethodGet, url, payload, reqH.Authorization)
		if err != nil {
			loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " TLAuthReports fetchFundsRes call api error =", err, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.AuthenticationFailed]
			resJS.ErrorCode = constants.AuthenticationFailed
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		tlErrorRes := models.Auth{}
		err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
		tlAuthRes := models.Auth{}
		json.Unmarshal([]byte(string(body)), &tlAuthRes)
		if tlAuthRes.Status != "success" {
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.AuthenticationFailed]
			resJS.ErrorCode = constants.AuthenticationFailed
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}

		matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqH.ClientId, reqH.Authorization)
		if !tokenValidStatus {
			loggerconfig.Error("TLAuthReports CheckAuthWithClient invalid authtoken", " clientId: ", reqH.ClientId, " requestId:", reqH.RequestId)
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.InvalidToken]
			resJS.ErrorCode = constants.InvalidToken
			c.JSON(http.StatusUnauthorized, resJS)
			c.Abort()
			return
		}

		if !matchStatus {
			loggerconfig.Error("TLAuthReports CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqH.ClientId, " requestId:", reqH.RequestId)
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.MismatchAuthClient]
			resJS.ErrorCode = constants.MismatchAuthClient
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}

		tlGetProfileResponse := tradelab.TradeLabProfileResponse{}
		json.Unmarshal([]byte(string(body)), &tlGetProfileResponse)

		if res.StatusCode != http.StatusOK {
			loggerconfig.Error("TLAuthReports tl status not ok StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
			resJS.Status = false
			resJS.Message = tlGetProfileResponse.Message
			c.JSON(res.StatusCode, resJS)
			c.Abort()
			return
		}

		var profileData models.ProfileDataResp
		profileData.Name = tlGetProfileResponse.Data.Name
		if len(profileData.BoID) > 0 {
			profileData.BoID = tlGetProfileResponse.Data.BoID[0]
		}
		profileData.EmailID = tlGetProfileResponse.Data.EmailID
		profileData.PanNumber = tlGetProfileResponse.Data.PanNumber
		profileData.ClientID = reqH.ClientId

		c.Set("profileData", profileData)
		c.Set("reqH", reqH)
		c.Next()
	}
}

func UserAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqH models.ReqHeader
		if err := c.ShouldBindHeader(&reqH); err != nil {
			loggerconfig.Warn("invalid  headers=", err)
		}
		reqH.RequestId = uuid.New().String()

		var resJS apihelpers.APIRes

		if reqH.ClientId == "" {
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.InvalidClient]
			resJS.ErrorCode = constants.InvalidClient
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}

		if len(reqH.Authorization) <= 7 {
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.TokenMissing]
			resJS.ErrorCode = constants.TokenMissing
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}

		// auth for admin calls
		adminRequestKey := reqH.AdminRequestKey

		if reqH.ClientType == constants.ADMIN && constants.AdminSecretKey == adminRequestKey {
			loggerconfig.Info("Admin access granted with secret key")
			c.Set("reqH", reqH)
			c.Next()
			return
		}

		redisCli := cache.GetRedisClientObj()

		reqH.ClientId = strings.ToUpper(reqH.ClientId)

		key := reqH.ClientId + "_" + reqH.Authorization[7:]

		exists, err := redisCli.Exists(key).Result()
		if err != nil {
			loggerconfig.Error("UserAuthenticationv Error checking key existence:", err, " clientId: ", reqH.ClientId)
		}

		if exists > 0 {
			loggerconfig.Info("UserAuthentication token found in redis for client:", reqH.ClientId)
			c.Set("reqH", reqH)
			c.Next()
			return
		} else {
			loggerconfig.Info("UserAuthentication token not found in redis for client:", reqH.ClientId)
		}

		tokenHeader, err := helpers.ExtractTokenHeader(reqH.Authorization[7:])
		if err != nil {
			// Invalid authtoken
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.AuthenticationFailed]
			resJS.ErrorCode = constants.AuthenticationFailed
			c.JSON(http.StatusUnauthorized, resJS)
			c.Abort()
			return
		}

		url := constants.TLURL + tradelab.FETCHFUNDSURL + "?type=" + constants.FetchFundsTypeTL + "&client_id=" + url.QueryEscape(strings.ToUpper(tokenHeader.ClientID))

		payload := new(bytes.Buffer)
		//call api
		res, err := apihelpers.CallApiTradeLab(http.MethodGet, url, payload, reqH.Authorization)
		if err != nil {
			loggerconfig.Error("Alert Severity:P0-Critical, UserAuthentication fetchFundsRes call api error =", err, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.AuthenticationFailed]
			resJS.ErrorCode = constants.AuthenticationFailed
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		tlAuthRes := models.Auth{}
		json.Unmarshal([]byte(string(body)), &tlAuthRes)
		if res.StatusCode != http.StatusOK || tlAuthRes.Status != "success" {
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.AuthenticationFailed]
			resJS.ErrorCode = constants.AuthenticationFailed
			c.JSON(http.StatusUnauthorized, resJS)
			c.Abort()
			return
		}

		matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqH.ClientId, reqH.Authorization)
		if !tokenValidStatus {
			loggerconfig.Error("UserAuthentication CheckAuthWithClient invalid authtoken", " clientId: ", reqH.ClientId, " requestId:", reqH.RequestId)
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.InvalidToken]
			resJS.ErrorCode = constants.InvalidToken
			c.JSON(http.StatusUnauthorized, resJS)
			c.Abort()
			return
		}

		if !matchStatus {
			loggerconfig.Error("UserAuthentication CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqH.ClientId, " requestId:", reqH.RequestId)
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.MismatchAuthClient]
			resJS.ErrorCode = constants.MismatchAuthClient
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}

		currentUnixTime := helpers.GetCurrentTimeInIST().Unix()
		var cacheTime time.Duration
		remainTokenExpiry := tokenHeader.Exp - currentUnixTime
		if remainTokenExpiry < (int64)(constants.TokenCacheTime) {
			cacheTime = time.Duration(remainTokenExpiry)
		} else {
			cacheTime = time.Duration(constants.TokenCacheTime)
		}

		// cacheTime
		errRedisSet := redisCli.SetRedis(key, "", cacheTime)
		if errRedisSet != nil {
			loggerconfig.Error("UserAuthentication Error setting auth from redis:", errRedisSet, " clientId: ", reqH.ClientId)
		}
		loggerconfig.Info("UserAuthentication token set in cache for clientid:", reqH.ClientId)

		c.Set("reqH", reqH)
		c.Next()

	}
}

func AuthCombinedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqH models.ReqHeader
		if err := c.ShouldBindHeader(&reqH); err != nil {
			loggerconfig.Warn("AuthCombineMiddleware invalid  headers=", err)
			var resJS apihelpers.APIRes
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.InvalidHeader]
			resJS.ErrorCode = constants.InvalidHeader
			c.JSON(http.StatusBadRequest, resJS)
			c.Abort()
			return
		}

		if strings.EqualFold(strings.TrimSpace(reqH.ClientType), constants.GUESTUSERTYPE) {
			Middleware()(c) // guest dao middleware
		} else {
			UserAuthentication()(c)
		}
	}
}
