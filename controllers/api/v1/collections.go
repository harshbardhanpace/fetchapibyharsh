package v1

import (
	"encoding/json"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
)

var collectionProvider models.CollectionsProvider

func InitCollectionProvider(provider models.CollectionsProvider) {
	defer models.HandlePanic()
	collectionProvider = provider
}

// CreateCollections
// @Tags space admin collections V1
// @Description Create Collections
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.CreateCollectionsRequest true "Collections"
// @Success 200 {object} apihelpers.APIRes{data=models.CreateCollectionsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/adminapis/createCollections [POST]
func CreateCollections(c *gin.Context) {
	var CollectionsReq models.CreateCollectionsRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("CreateCollections (controller), error parsing header, error:", err, "clientID: ", reqH.ClientId)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&CollectionsReq)
	if err != nil {
		loggerconfig.Error("CreateCollections (controller), error decoding body, error:", err, "clientID: ", reqH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	//call service
	code, resp := collectionProvider.CreateCollections(CollectionsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: CreateCollections requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ModifyCollections
// @Tags space admin collections V1
// @Description Modify Collections
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ModifyCollectionsRequest true "Collections"
// @Success 200 {object} apihelpers.APIRes{data=models.ModifyCollectionsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/adminapis/modifyCollections [POST]
func ModifyCollections(c *gin.Context) {
	var CollectionsReq models.ModifyCollectionsRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("ModifyCollections (controller), error parsing header, error:", err, "clientID: ", reqH.ClientId)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&CollectionsReq)
	if err != nil {
		loggerconfig.Error("ModifyCollections (controller), error decoding body, error:", err, "clientID: ", reqH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	//call service
	code, resp := collectionProvider.ModifyCollections(CollectionsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: ModifyCollections requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Collections
// @Tags space admin collections V1
// @Description Collections
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchCollectionsDetailsRequest true "Collections"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchCollectionsDetailsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/adminapis/fetchCollections [POST]
func FetchCollections(c *gin.Context) {
	var CollectionsReq models.FetchCollectionsDetailsRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchCollections (controller), error parsing header, error:", err, "clientID: ", reqH.ClientId)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&CollectionsReq)
	if err != nil {
		loggerconfig.Error("FetchCollections (controller), error decoding body, error:", err, "clientID: ", reqH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	//call service
	code, resp := collectionProvider.FetchCollections(CollectionsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchCollections requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Collections User
// @Tags space user collections V1
// @Description Collections User
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchCollectionsDetailsRequest true "Collections"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchCollectionsDetailsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/collections/fetchCollections [POST]
func FetchCollectionsUser(c *gin.Context) {
	var CollectionsReq models.FetchCollectionsDetailsRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchCollectionsUser (controller), error parsing header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&CollectionsReq)
	if err != nil {
		loggerconfig.Error("FetchCollectionsUser (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	//call service
	code, resp := collectionProvider.FetchCollections(CollectionsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchCollectionsUser requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Collections
// @Tags space admin collections V1
// @Description Collections
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.DeleteCollectionsRequest true "Collections"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/adminapis/deleteCollections [POST]
func DeleteCollections(c *gin.Context) {
	var CollectionsReq models.DeleteCollectionsRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("DeleteCollections (controller), error parsing header, error:", err, "clientID: ", reqH.ClientId)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&CollectionsReq)
	if err != nil {
		loggerconfig.Error("DeleteCollections (controller), error decoding body, error:", err, "clientID: ", reqH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	//call service
	code, resp := collectionProvider.DeleteCollections(CollectionsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: DeleteCollections requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Collections
// @Tags space admin collections V1
// @Description Collections
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchAllCollectionsDetailsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/adminapis/fetchAllCollections [GET]
func FetchAllCollections(c *gin.Context) {

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchAllCollections (controller), error parsing header, error:", err, "clientID: ", reqH.ClientId)
	}

	//call service
	code, resp := collectionProvider.FetchAllCollections(reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchAllCollections requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Collections User
// @Tags space user collections V1
// @Description Collections User
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchAllCollectionsDetailsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/collections/fetchAllCollections [GET]
func FetchAllCollectionsUser(c *gin.Context) {

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchAllCollectionsUser (controller), error parsing header, error:", err, "clientID: ", reqH.ClientId)
	}

	//call service
	code, resp := collectionProvider.FetchAllCollections(reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchAllCollections requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
