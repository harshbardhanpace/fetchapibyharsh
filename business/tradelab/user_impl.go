package tradelab

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers"
	"space/helpers/cache"
	"space/models"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type LogoutObj struct {
	tradeLabURL string
}

func InitLogoutProvider() LogoutObj {

	logoutObj := LogoutObj{
		tradeLabURL: constants.TLURL,
	}

	return logoutObj
}

func (obj LogoutObj) LogoutSingleDevice(reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + LOGOUTURL

	//make payload
	payload := new(bytes.Buffer)

	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "LogoutSingleDevice", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		logrus.Error("LogoutSingleDevice call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		logrus.Error("LogoutSingleDevice tl res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		apiRes.Data = tlErrorRes.Data
		return res.StatusCode, apiRes
	}

	if res.StatusCode != http.StatusOK {
		logrus.Error("LogoutSingleDevice tl status not ok =", res.StatusCode, " uccId:", reqH.ClientId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		apiRes.Data = tlErrorRes.Data
		return res.StatusCode, apiRes
	}

	logrus.Info("LogoutSingleDevice tl success uccId:", reqH.ClientId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	key := strings.ToUpper(reqH.ClientId) + "_" + reqH.Authorization[7:]
	delStatus := dbops.RedisRepo.Delete(key)
	logrus.Info("LogoutSingleDevice auth key removed from redis: ", delStatus, " uccid: ", reqH.ClientId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	if reqH.ClientId != "" && reqH.FCMToken != "" {
		// delete FCM Token
		err = db.GetPgObj().DeleteFCMToken(reqH.ClientId, reqH.FCMToken)
		if err != nil {
			logrus.Error("Alert Severity:P1-Critical, Error (LogoutSingleDevice) while deleting FCMToken: ", err.Error(), " uccId: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		}

		if constants.KafkaEnable {

			redisCliObj := cache.GetRedisClientObj()

			err = redisCliObj.SRem(constants.ClientMembers+strings.ToUpper(reqH.ClientId), reqH.FCMToken)
			if err != nil && err != redis.Nil {
				logrus.Error("Alert Severity:P1-Critical, LogoutSingleDevice Error (Logout) removing FCM token in Redis:", err, " uccId: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
			}

			logrus.Info("LogoutSingleDevice FCM token removed from redis: ", reqH.FCMToken, " uccid: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

		}
	}

	apiRes.Message = tlErrorRes.Message
	apiRes.Status = true
	apiRes.Data = tlErrorRes.Data

	return http.StatusOK, apiRes
}
