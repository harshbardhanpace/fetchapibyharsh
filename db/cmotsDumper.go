package db

// commenting this code for future reference

// func GetAllCoCode() {
// 	var dbResponse []models.CompanyDetails

// 	dbResponse, err := GetPgObj().FetchCompanyMaster()
// 	if err != nil {
// 		loggerconfig.Error("Error in getting the cmpany master details error: ", err)
// 	}

// 	for i := 0; i < len(dbResponse); i++ {
// 		var req models.FetchFinancialsReq
// 		req.Isin = dbResponse[i].Isin
// 		res, err := GetPgObj().FetchFinancialsDataV4(req)
// 		if err != nil {
// 			loggerconfig.Error("Here is the error in the code, from fetching data from the db error :", err)
// 		}

// 		mongoFetchFinancials := &models.MongoFetchFinancials{
// 			CoCode:     dbResponse[i].CoCode,
// 			Isin:       dbResponse[i].Isin,
// 			Bsecode:    dbResponse[i].Bsecode,
// 			Nsesymbol:  dbResponse[i].Nsesymbol,
// 			Financials: res,
// 		}

// 		collection := models.GetMongoCollection(constants.FETCHFINANCIALS)
// 		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 		filter := bson.D{{"isin", req.Isin}}
// 		update := bson.D{{"$set", mongoFetchFinancials}}
// 		opts := options.Update().SetUpsert(true)
// 		_, err = collection.UpdateOne(ctx, filter, update, opts)
// 		if err != nil {
// 			loggerconfig.Error("Intertion in mongo failed cocode: ", dbResponse[i].CoCode, "error :", err)
// 			return
// 		}
// 		loggerconfig.Info("the Insertion in mongo is successsfull for this cocode", dbResponse[i].CoCode)
// 	}
// }
