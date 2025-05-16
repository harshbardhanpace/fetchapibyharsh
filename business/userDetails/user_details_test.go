package userdetails

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/loggerconfig"
	"space/models"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestUserNotifications(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.UserNotificationsReq
		reqH models.ReqHeader
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	req1 := models.UserNotificationsReq{
		ClientId: "TEST",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	field1 := fields{
		tradeLabURL: "http://test",
	}

	db.CallFindAllMongo = func(db db.MongoDatabase, collectionName string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
		var ctx *mongo.Cursor
		return ctx, errors.New("Something went wrong")
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceId, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UserDetailsObj{}
			got, got1 := obj.UserNotifications(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("UserDetailsObj.UserNotifications() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("UserDetailsObj.UserNotifications() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	db.CallFindAllMongo = func(db db.MongoDatabase, collectionName string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
		var ctx *mongo.Cursor
		return ctx, nil
	}

	readValueMongoCurs = func(req models.UserNotificationsReq, ctx context.Context, cursor *mongo.Cursor, mongoNotificationStore *[]models.MongoNotificationStore, reqH models.ReqHeader) error {
		return nil
	}

	var mongoNotificationStore []models.MongoNotificationStore
	// mongoNotificationStore

	res2 := apihelpers.APIRes{
		Data:    mongoNotificationStore,
		Status:  true,
		Message: "SUCCESS",
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field1, arg1, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UserDetailsObj{}
			got, got1 := obj.UserNotifications(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("UserDetailsObj.UserNotifications() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("UserDetailsObj.UserNotifications() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}
