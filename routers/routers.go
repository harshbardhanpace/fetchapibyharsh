package routers

import (
	healthController "space/controllers/api/health"
	apiControllerV1 "space/controllers/api/v1"
	blockDealController "space/controllers/api/v1"
	apiControllerV2 "space/controllers/api/v2"
	apiControllerV3 "space/controllers/api/v3"
	apiControllerV4 "space/controllers/api/v4"
	"space/middlewares"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter function will perform all route operations
func SetupRouter() *gin.Engine {

	r := gin.Default()

	//Giving access to storage folder

	//Giving access to storage
	r.Static("/storage", "storage")
	r.GET("/swagger/space/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Giving access to template folder
	r.Static("/templates", "templates")
	r.LoadHTMLGlob("templates/*")

	r.Use(middlewares.RecoveryMiddleware())

	r.Use(func(c *gin.Context) {
		// add header Access-Control-Allow-Origin
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max, DNT, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Range, p-devicetype, p-platform, clientid, Origin, Accept")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	})

	// Health check api route
	healthCheck := r.Group("api/space")
	healthCheck.Use()
	{
		healthCheck.GET("/health", healthController.GetHealthStatus)
	}

	//API route for version 1
	v1Login := r.Group("/api/space/v1/authapis")

	//If you want to pass your route through specific middlewares
	v1Login.Use(middlewares.Middleware())
	{
		// v1.POST("user-list", apiControllerV1.UserList)
		v1Login.POST("/login", apiControllerV1.Login)
		v1Login.POST("/loginByEmail", apiControllerV1.LoginByEmail)
		v1Login.POST("/validateToken", apiControllerV1.ValidateToken)
		v1Login.POST("/validateTwoFa", apiControllerV1.ValidateTwoFA)
		v1Login.POST("/setTwoFaPin", apiControllerV1.SetTwoFAPin)
		v1Login.POST("/forgotPassword", apiControllerV1.ForgotPassword)
		v1Login.POST("/forgetPasswordEmail", apiControllerV1.ForgetPasswordEmail)
		v1Login.POST("/setPassword", apiControllerV1.SetPassword)
		v1Login.POST("/forgetResetTwoFa", apiControllerV1.ForgetResetTwoFa)
		v1Login.POST("/forgetResetTwoFaEmail", apiControllerV1.ForgetResetTwoFaEmail)
		v1Login.POST("/guestUserStatus", apiControllerV1.GuestUserStatus)
		v1Login.POST("/unblockUser", apiControllerV1.UnblockUser)
	}

	v2Login := r.Group("/api/space/v2/authapis")

	v2Login.Use(middlewares.Middleware())
	{
		v2Login.POST("/login", apiControllerV2.Login)
		v2Login.POST("/validateTwofa", apiControllerV2.ValidateTwofa)
		v2Login.POST("/setupTotp", apiControllerV2.SetupTotp)
		v2Login.POST("/chooseTwofa", apiControllerV2.ChooseTwofa)
		v2Login.POST("/forgetTotp", apiControllerV2.ForgetTotp)
		v2Login.POST("/validateLoginOtp", apiControllerV2.ValidateLoginOtp)
		v2Login.POST("/setupBiometric", apiControllerV2.SetupBiometric)
		v2Login.POST("/disableBiometric", apiControllerV2.DisableBiometric)
		v2Login.POST("/forgotPassword", apiControllerV2.ForgotPasswordV2)
		v2Login.POST("/unblockUser", apiControllerV2.UnblockUserV2)
	}

	V2Logout := r.Group("/api/space/v2/user")
	V2Logout.Use(middlewares.Middleware())
	{
		V2Logout.DELETE("/logout", apiControllerV2.Logout)
	}

	v3Login := r.Group("/api/space/v3/authapis")
	v3Login.Use(middlewares.Middleware())
	{
		v3Login.PUT("/setPassword", apiControllerV3.SetPassword)
		v3Login.GET("/validateToken", apiControllerV3.ValidateToken)
		v3Login.PUT("/forgetResetTwoFa", apiControllerV3.ForgetResetTwoFa)
		v3Login.PUT("/validateLoginOtp", apiControllerV3.ValidateLoginOtp)
		v3Login.PUT("/setupBiometric", apiControllerV3.SetupBiometric)
		v3Login.DELETE("/disableBiometric", apiControllerV3.DisableBiometric)
		v3Login.POST("/loginByEmailOtp", apiControllerV3.LoginByEmailOtp)
	}

	v1Finvu := r.Group("/api/space/v1/finvu")
	v1Finvu.Use(middlewares.UserAuthentication())
	{
		v1Finvu.POST("/finvuConsentRequestPlus", apiControllerV1.FinvuConsentRequestPlus)
		v1Finvu.POST("/finvuGetBankStatement", apiControllerV1.FinvuGetBankStatement)
	}

	v1order := r.Group("/api/space/v1/orderapis")

	//If you want to pass your route through specific middlewares
	v1order.Use(middlewares.Middleware())
	{
		// v1.POST("user-list", apiControllerV1.UserList)
		v1order.POST("/placeOrder", apiControllerV1.PlaceOrder)
		v1order.POST("/modifyOrder", apiControllerV1.ModifyOrder)
		v1order.POST("/cancelOrder", apiControllerV1.CancelOrder)
		v1order.POST("/placeAMOOrder", apiControllerV1.PlaceAMOOrder)
		v1order.POST("/modifyAMOOrder", apiControllerV1.ModifyAMOOrder)
		v1order.POST("/cancelAMOOrder", apiControllerV1.CancelAMOOrder)
		v1order.POST("/pendingOrder", apiControllerV1.PendingOrder)
		v1order.POST("/completedOrder", apiControllerV1.CompletedOrder)
		v1order.POST("/tradeBook", apiControllerV1.TradeBook)
		v1order.POST("/orderHistory", apiControllerV1.OrderHistory)
		v1order.POST("/marginCalculations", apiControllerV1.MarginCalculations)

		//router for conditional orders
		v1order.POST("/placeBOOrder", apiControllerV1.PlaceBOOrder)
		v1order.POST("/modifyBOOrder", apiControllerV1.ModifyBOOrder)
		v1order.POST("/exitBOOrder", apiControllerV1.ExitBOOrder)

		v1order.POST("/placeCOOrder", apiControllerV1.PlaceCOOrder)
		v1order.POST("/modifyCOOrder", apiControllerV1.ModifyCOOrder)
		v1order.POST("/exitCOOrder", apiControllerV1.ExitCOOrder)

		v1order.POST("/placeSpreadOrder", apiControllerV1.PlaceSpreadOrder)
		v1order.POST("/modifySpreadOrder", apiControllerV1.ModifySpreadOrder)
		v1order.POST("/exitSpreadOrder", apiControllerV1.ExitSpreadOrder)

		//gtt orders routes
		v1order.POST("/createGTTOrder", apiControllerV1.CreateGTTOrder)
		v1order.POST("/modifyGTTOrder", apiControllerV1.ModifyGTTOrder)
		v1order.POST("/cancelGTTOrder", apiControllerV1.CancelGTTOrder)
		v1order.POST("/fetchGTTOrder", apiControllerV1.FetchGTTOrder)
		v1order.POST("/PlaceGTTOcoOrder", apiControllerV1.PlaceGttOCOOrder)
		v1order.POST("/lastTradedPrice", apiControllerV1.LastTradedPrice)

		//iceberg orders routes
		v1order.POST("/createIcebergOrder", apiControllerV1.CreateIcebergOrder)
		v1order.PUT("/modifyIcebergOrder", apiControllerV1.ModifyIcebergOrder)
		v1order.DELETE("/cancelIcebergOrder", apiControllerV1.CancelIcebergOrder)

	}

	v2order := r.Group("/api/space/v2/orderapis")
	v2order.Use(middlewares.Middleware())
	{
		v2order.GET("/pendingOrder", apiControllerV2.PendingOrder)
		v2order.GET("/completedOrder", apiControllerV2.CompletedOrder)
		v2order.GET("/tradeBook", apiControllerV2.TradeBook)
		v2order.GET("/orderHistory", apiControllerV2.OrderHistory)
		v2order.POST("/placeOrder", apiControllerV2.PlaceOrder)
		v2order.PATCH("/modifyOrder", apiControllerV2.ModifyOrder)
		v2order.POST("/CancelOrder", apiControllerV2.CancelOrder)
	}

	v1portfolio := r.Group("/api/space/v1/portfolioapis")
	//If you want to pass your route through specific middlewares
	v1portfolio.Use(middlewares.Middleware())
	{
		v1portfolio.POST("/fetchDematHoldings", apiControllerV1.FetchDematHoldings)
		v1portfolio.POST("/convertPositions", apiControllerV1.ConvertPositions)
		v1portfolio.POST("/getPositions", apiControllerV1.GetPositions)
	}

	v2portfolio := r.Group("/api/space/v2/portfolioapis")
	v2portfolio.Use(middlewares.Middleware())
	{
		v2portfolio.GET("/fetchDematHoldings", apiControllerV2.FetchDematHoldings)
		v2portfolio.PUT("/convertPositions", apiControllerV2.ConvertPositions)
		v2portfolio.GET("/getPositions", apiControllerV2.GetPositions)
	}

	v1optionchain := r.Group("/api/space/v1/optionchain")
	//If you want to pass your route through specific middlewares
	v1optionchain.Use(middlewares.Middleware())
	{
		v1optionchain.POST("/fetchOptionChain", apiControllerV1.FetchOptionChain)
		v1optionchain.POST("/fetchFuturesChain", apiControllerV1.FetchFuturesChain)
	}

	v2optionchain := r.Group("/api/space/v2/optionchain")
	v2optionchain.Use(middlewares.AuthCombinedMiddleware())
	{
		v2optionchain.POST("/fetchOptionChain", apiControllerV2.FetchOptionChainV2)
		v2optionchain.POST("/fetchOptionChainByExpiry", apiControllerV2.FetchOptionChainByExpiryV2)
	}

	v3optionchain := r.Group("/api/space/v3/optionchain")
	v3optionchain.Use(middlewares.Middleware())
	{
		v3optionchain.GET("/fetchOptionChain", apiControllerV3.FetchOptionChain)
		v3optionchain.GET("/fetchFuturesChain", apiControllerV3.FetchFuturesChain)
	}

	v1profile := r.Group("/api/space/v1/user/profile")

	v1profile.Use(middlewares.Middleware())
	{
		v1profile.POST("/getProfile", apiControllerV1.GetProfile)
		v1profile.POST("/sendAFOtp", apiControllerV1.SendAFOtp)
		v1profile.POST("/verifyAFOtp", apiControllerV1.VerifyAFOtp)
		v1profile.POST("/accountFreeze", apiControllerV1.AccountFreeze)
	}

	v2profile := r.Group("/api/space/v2/user/profile")
	v2profile.Use(middlewares.Middleware())
	{
		v2profile.GET("/getProfile", apiControllerV2.GetProfile)
	}

	v1funds := r.Group("/api/space/v1/funds/view")
	v1funds.Use(middlewares.Middleware())
	{
		v1funds.POST("/fetchFunds", apiControllerV1.FetchFunds)
		v1funds.POST("/cancelPayout", apiControllerV1.CancelPayout)
		v1funds.POST("/payout", apiControllerV1.Payout)
		v1funds.POST("/clientTransactions", apiControllerV1.ClientTransactions)
	}

	v2funds := r.Group("/api/space/v2/funds/view")
	v2funds.Use(middlewares.Middleware())
	{
		v2funds.GET("/fetchFunds", apiControllerV2.FetchFunds)
		v2funds.PUT("/cancelPayout", apiControllerV2.CancelPayout)
		v2funds.GET("/clientTransactions", apiControllerV2.ClientTransactions)
	}

	v3funds := r.Group("/api/space/v3/funds/view")
	v3funds.Use(middlewares.Middleware())
	{
		v3funds.PUT("/cancelPayout", apiControllerV3.CancelPayout)
		v3funds.POST("/payout", apiControllerV3.Payout)
	}

	v1Alerts := r.Group("/api/space/v1/alerts/")
	v1Alerts.Use(middlewares.Middleware())
	{
		v1Alerts.POST("/setAlerts", apiControllerV1.SetAlerts)
		v1Alerts.POST("/editAlerts", apiControllerV1.EditAlerts)
		v1Alerts.POST("/getAlerts", apiControllerV1.GetAlerts)
		v1Alerts.POST("/pauseAlerts", apiControllerV1.PauseAlerts)
		v1Alerts.POST("/deleteAlerts", apiControllerV1.DeleteAlerts)
	}

	v1contractdetails := r.Group("/api/space/v1/contractdetails")
	//If you want to pass your route through specific middlewares
	v1contractdetails.Use(middlewares.Middleware())
	{
		v1contractdetails.POST("/searchScrip", apiControllerV1.SearchScrip)
		v1contractdetails.POST("/scripInfo", apiControllerV1.ScripInfo)
	}

	v2contractdetailsTL := r.Group("/api/space/v2/contractdetails")
	v2contractdetailsTL.Use(middlewares.Middleware())
	{
		v2contractdetailsTL.GET("/searchScripTL", apiControllerV2.SearchScrip)
		v2contractdetailsTL.GET("/scripInfoTL", apiControllerV2.ScripInfo)
	}

	v1adminLogin := r.Group("/api/space/v1/adminapis/")
	v1adminLogin.Use(middlewares.Middleware())
	{
		v1adminLogin.POST("/adminLogin", apiControllerV1.AdminLogin)
	}

	v1pocketsAdmin := r.Group("/api/space/v1/adminapis/")
	v1pocketsAdmin.Use(middlewares.AdminMiddleware())
	{
		v1pocketsAdmin.POST("/createPockets", apiControllerV1.CreatePockets)
		v1pocketsAdmin.POST("/modifyPockets", apiControllerV1.ModifyPockets)
		v1pocketsAdmin.POST("/deletePockets", apiControllerV1.DeletePockets)
		v1pocketsAdmin.POST("/fetchPockets", apiControllerV1.FetchPockets)
		v1pocketsAdmin.GET("/fetchAllPockets", apiControllerV1.FetchAllPockets)

		v1pocketsAdmin.POST("/createCollections", apiControllerV1.CreateCollections)
		v1pocketsAdmin.POST("/modifyCollections", apiControllerV1.ModifyCollections)
		v1pocketsAdmin.POST("/deleteCollections", apiControllerV1.DeleteCollections)
		v1pocketsAdmin.POST("/fetchCollections", apiControllerV1.FetchCollections)
		v1pocketsAdmin.GET("/fetchAllCollections", apiControllerV1.FetchAllCollections)

	}

	v1collections := r.Group("/api/space/v1/collections")
	v1collections.Use(middlewares.AuthCombinedMiddleware())
	{
		v1collections.POST("/fetchCollections", apiControllerV1.FetchCollectionsUser)
		v1collections.GET("/fetchAllCollections", apiControllerV1.FetchAllCollectionsUser)
	}

	v1pockets := r.Group("/api/space/v1/pockets")
	v1pockets.Use(middlewares.UserAuthentication())
	{
		v1pockets.POST("/fetchPocketPortfolio", apiControllerV1.FetchPocketPortfolio)
		v1pockets.POST("/buyPocket", apiControllerV1.BuyPocket)
		v1pockets.POST("/exitPocket", apiControllerV1.ExitPocket)
		v1pockets.POST("/fetchPocketTransaction", apiControllerV1.FetchPocketTransaction)
		v1pockets.POST("/pocketsCalculations", apiControllerV1.PocketsCalculations)
		v1pockets.POST("/multipleAndIndividualStocksCalculations", apiControllerV1.MultipleAndIndividualStocksCalculations)
		v1pockets.POST("/storePocketTransaction", apiControllerV1.StorePocketTransaction)
	}
	v1pockets.Use(middlewares.AuthCombinedMiddleware())
	{
		v1pockets.GET("/fetchAllPockets", apiControllerV1.FetchAllPocketsUser)
		v1pockets.POST("/fetchPockets", apiControllerV1.FetchPocketsUser)
	}

	v1WatchList := r.Group("/api/space/v1/watchlist")
	v1WatchList.Use(middlewares.UserAuthentication())
	{
		v1WatchList.POST("/createWatchlist", apiControllerV1.CreateWatchList)
		v1WatchList.POST("/modifyWatchList", apiControllerV1.ModifyWatchList)
		v1WatchList.POST("/fetchWatchList", apiControllerV1.FetchWatchList)
		v1WatchList.POST("/deleteWatchList", apiControllerV1.DeleteWatchList)
		v1WatchList.POST("/fetchWatchListDetails", apiControllerV1.FetchWatchListDetails)
		v1WatchList.POST("/addStockToWatchList", apiControllerV1.AddStockToWatchList)
		//	v1WatchList.POST("/deleteStockInWatchList", apiControllerV1.DeleteStockInWatchList)
	}

	v1Pins := r.Group("/api/space/v1/pins")
	v1Pins.Use(middlewares.Middleware())
	{
		v1Pins.POST("/fetchPins", apiControllerV1.FetchPins)
		v1Pins.POST("/updatePins", apiControllerV1.UpdatePins)
		v1Pins.POST("/addPins", apiControllerV1.AddPins)
		v1Pins.POST("/deletePins", apiControllerV1.DeletePins)
	}

	v2Pins := r.Group("/api/space/v2/pins")
	v2Pins.Use(middlewares.AuthCombinedMiddleware())
	{
		v2Pins.GET("/fetchPins", apiControllerV2.FetchPinsV2)
	}
	v2WatchList := r.Group("/api/space/v2/watchlist")
	v2WatchList.Use(middlewares.UserAuthentication())
	{
		v2WatchList.POST("/addStockToWatchList", apiControllerV2.AddStockToWatchList)
		v2WatchList.POST("/fetchWatchList", apiControllerV2.FetchWatchList)
		v2WatchList.POST("/deleteStockInWatchList", apiControllerV2.DeleteStockInWatchList)
		v2WatchList.POST("/arrangeStocksWatchList", apiControllerV2.ArrangeStocksWatchList)
		v2WatchList.POST("/deleteStockInWatchListUpdated", apiControllerV2.DeleteStockInWatchListUpdated)
		//	v1WatchList.POST("/deleteStockInWatchList", apiControllerV1.DeleteStockInWatchList)
	}

	v1Ipo := r.Group("api/space/v1/tradeipo")
	v1Ipo.Use(middlewares.Middleware())
	{
		v1Ipo.POST("/placeIpoOrder", apiControllerV1.PlaceIpoOrder)
		v1Ipo.POST("/fetchIpoOrder", apiControllerV1.FetchIpoOrder)
		v1Ipo.POST("/cancelIpoOrder", apiControllerV1.CancelIpoOrder)
	}

	v1Ipo.Use(middlewares.AuthCombinedMiddleware())
	{
		v1Ipo.POST("/fetchIpoDataNse", apiControllerV1.FetchEIpo)
		v1Ipo.GET("/getAllIpo", apiControllerV1.GetAllIpo)
		v1Ipo.POST("/fetchIpoData", apiControllerV1.FetchIpoData)
		v1Ipo.POST("/fetchIpoGmpData", apiControllerV1.FetchIpoGmpData)
	}

	v1Screeners := r.Group("api/space/v1/screeners")
	v1Screeners.Use(middlewares.AuthCombinedMiddleware())
	{
		v1Screeners.POST("/gainersloser", apiControllerV1.GainerLoser)
		v1Screeners.POST("/mostActiveVolume", apiControllerV1.MostActiveVolume)
		v1Screeners.POST("/chartData", apiControllerV1.ChartData)
		v1Screeners.POST("/returnOnInvestment", apiControllerV1.ReturnOnInvestment)
		v1Screeners.POST("/fetchHistoricPerformance", apiControllerV1.FetchHistoricPerformance)
		v1Screeners.POST("/fetchHistoricPerformance/all", apiControllerV1.FetchAllHistoricPerformance)

	}

	v2Screeners := r.Group("api/space/v2/screeners")
	v2Screeners.Use(middlewares.AuthCombinedMiddleware())
	{
		v2Screeners.GET("/gainersloser", apiControllerV2.GainerLoser)
		v2Screeners.GET("/mostActiveVolume", apiControllerV2.MostActiveVolume)
		v2Screeners.GET("/chartData", apiControllerV2.ChartData)
		v2Screeners.GET("/returnOnInvestment", apiControllerV2.ReturnOnInvestment)
	}

	v1BasketOrder := r.Group("api/space/v1/basket")
	v1BasketOrder.Use(middlewares.Middleware())
	{
		v1BasketOrder.POST("/createBasket", apiControllerV1.CreateBasket)
		v1BasketOrder.POST("/fetchBasket", apiControllerV1.FetchBasket)
		v1BasketOrder.POST("/deleteBasket", apiControllerV1.DeleteBasket)
		v1BasketOrder.POST("/addBasketInstrument", apiControllerV1.AddBasketInstrument)
		v1BasketOrder.POST("/editBasketInstrument", apiControllerV1.EditBasketInstrument)
		v1BasketOrder.POST("/deleteBasketInstrument", apiControllerV1.DeleteBasketInstrument)
		v1BasketOrder.POST("/renameBasket", apiControllerV1.RenameBasket)
		v1BasketOrder.POST("/executeBasket", apiControllerV1.ExecuteBasket)
		v1BasketOrder.POST("/updateBasketExecutionState", apiControllerV1.UpdateBasketExecutionState)
	}

	v1Charges := r.Group("/api/space/v1/charges")
	v1Charges.Use(middlewares.UserAuthentication())
	{
		v1Charges.POST("/brokerCharges", apiControllerV1.BrokerCharges)
		v1Charges.POST("/combineBrokerCharges", apiControllerV1.CombineBrokerCharges)
		v1Charges.POST("/fundsPayout", apiControllerV1.FundsPayout)
	}

	v1SquareOff := r.Group("/api/space/v1/squareOff")
	v1SquareOff.Use(middlewares.UserAuthentication())
	{
		v1SquareOff.POST("/squareOffAll", apiControllerV1.SquareOffAll)
	}

	v1Cmots := r.Group("/api/space/v1/cmots")
	v1Cmots.Use(middlewares.AuthCombinedMiddleware())
	{
		v1Cmots.POST("/getOverview", apiControllerV1.GetOverview)
		v1Cmots.POST("/fetchFinancials", apiControllerV1.FetchFinancials)
		v1Cmots.POST("/fetchFinancialsDetailed", apiControllerV1.FetchFinancialsDetailed)
		v1Cmots.POST("/fetchPeers", apiControllerV1.FetchPeers)
		v1Cmots.POST("/shareHoldingPatterns", apiControllerV1.ShareHoldingPatterns)
		v1Cmots.POST("/ratiosCompare", apiControllerV1.RatiosCompare)
		v1Cmots.POST("/fetchTechnicalIndicators", apiControllerV1.FetchTechnicalIndicators)
		v1Cmots.POST("/stocksOnNews", apiControllerV1.StocksOnNews)
		v1Cmots.GET("/fetchSectorList", apiControllerV1.FetchSectorList)
		v1Cmots.POST("/fetchSectorWiseCompany", apiControllerV1.FetchSectorWiseCompany)
		v1Cmots.POST("/fetchCompanyCategory", apiControllerV1.FetchCompanyCategory)
		v1Cmots.POST("/stocksAnalyzer", apiControllerV1.StocksAnalyzer)
		v1Cmots.POST("/corporateActionsIndividual", apiControllerV1.CorporateActionsIndividual)
		v1Cmots.POST("/corporateActionsAll", apiControllerV1.CorporateActionsAll)
		v1Cmots.GET("/getSectorWiseStockList", apiControllerV1.GetSectorWiseStockList)
	}

	blockDealGroup := r.Group("/api/v1/blockdeals")
	{
		blockDealGroup.GET("/getallblockdeals", blockDealController.GetAllBlockDeals)
		blockDealGroup.POST("/create", blockDealController.CreateBlockDeal)
		blockDealGroup.GET("/:cocode", blockDealController.GetBlockDealByCocode)
		blockDealGroup.PUT("/:cocode", blockDealController.UpdateBlockDeal)
		blockDealGroup.DELETE("/:cocode", blockDealController.DeleteBlockDeal)
	}

	v1LoginWithQR := r.Group("/api/space/v1/qr")
	v1LoginWithQR.Use(middlewares.UserAuthentication())
	{
		v1LoginWithQR.POST("/webLogin", apiControllerV1.QRWebLogin)
	}

	v2Cmots := r.Group("/api/space/v2/cmots")
	v2Cmots.Use(middlewares.AuthCombinedMiddleware())
	{
		v2Cmots.POST("/stocksOnNewsV2", apiControllerV2.StocksOnNewsV2)
		v2Cmots.POST("/fetchFinancialsV2", apiControllerV2.FetchFinancialsV2)
		v2Cmots.POST("/fetchPeers", apiControllerV2.FetchPeersV2)
		v2Cmots.GET("/fetchSectorListV2", apiControllerV2.FetchSectorListV2)
		v2Cmots.POST("/fetchSectorWiseCompanyV2", apiControllerV2.FetchSectorWiseCompanyV2)
	}

	v1User := r.Group("/api/space/v1/userDetails")
	v1User.Use(middlewares.UserAuthentication())
	{
		v1User.POST("/getAllBankAccounts", apiControllerV1.GetAllBankAccounts)
		v1User.GET("/userNotifications", apiControllerV1.UserNotifications)
	}
	v1User.Use(middlewares.AuthCombinedMiddleware())
	{
		v1User.POST("/getUserId", apiControllerV1.GetUserId)
	}

	v1UserStatus := r.Group("/api/space/v1/userDetails/getClientStatus")
	v1UserStatus.Use(middlewares.Middleware())
	{
		v1UserStatus.GET("/", apiControllerV1.GetClientStatus)
	}

	v1PortfolioAnalyzer := r.Group("/api/space/v1/portfolioAnalyzer")
	v1PortfolioAnalyzer.Use(middlewares.UserAuthentication())
	{
		v1PortfolioAnalyzer.POST("/holdingsWeightages", apiControllerV1.HoldingsWeightages)
		v1PortfolioAnalyzer.POST("/portfolioBeta", apiControllerV1.PortfolioBeta)
		v1PortfolioAnalyzer.POST("/portfolioPE", apiControllerV1.PortfolioPE)
		v1PortfolioAnalyzer.POST("/portfolioDE", apiControllerV1.PortfolioDE)
		v1PortfolioAnalyzer.POST("/highPledgedPromoterHoldings", apiControllerV1.HighPledgedPromoterHoldings)
		v1PortfolioAnalyzer.POST("/additionalSurveillanceMeasure", apiControllerV1.AdditionalSurveillanceMeasure)
		v1PortfolioAnalyzer.POST("/gradedSurveillanceMeasure", apiControllerV1.GradedSurveillanceMeasure)
		v1PortfolioAnalyzer.POST("/highDefaultProbability", apiControllerV1.HighDefaultProbability)
		v1PortfolioAnalyzer.POST("/lowROE", apiControllerV1.LowROE)
		v1PortfolioAnalyzer.POST("/lowProfitGrowth", apiControllerV1.LowProfitGrowth)
		v1PortfolioAnalyzer.POST("/holdingStockContribution", apiControllerV1.HoldingStockContribution)
		v1PortfolioAnalyzer.POST("/investmentSector", apiControllerV1.InvestmentSector)
		v1PortfolioAnalyzer.POST("/declineInPromoterHolding", apiControllerV1.DeclineInPromoterHolding)
		v1PortfolioAnalyzer.POST("/interestCoverageRatio", apiControllerV1.InterestCoverageRatio)
		v1PortfolioAnalyzer.POST("/declineInRevenueAndProfit", apiControllerV1.DeclineInRevenueAndProfit)
		v1PortfolioAnalyzer.POST("/lowNetWorth", apiControllerV1.LowNetWorth)
		v1PortfolioAnalyzer.POST("/declineInRevenue", apiControllerV1.DeclineInRevenue)
		v1PortfolioAnalyzer.POST("/promoterPledge", apiControllerV1.PromoterPledge)
		v1PortfolioAnalyzer.POST("/pennyStocks", apiControllerV1.PennyStocks)
		v1PortfolioAnalyzer.POST("/stockReturn", apiControllerV1.StockReturn)
		v1PortfolioAnalyzer.POST("/niftyVsPortfolio", apiControllerV1.NiftyVsPortfolio)
		v1PortfolioAnalyzer.POST("/changeInInstitutionalHolding", apiControllerV1.ChangeInInstitutionalHolding)
		v1PortfolioAnalyzer.POST("/roeAndStockReturn", apiControllerV1.RoeAndStockReturn)
		v1PortfolioAnalyzer.POST("/illiquidStocks", apiControllerV1.IlliquidStocks)
	}

	v1SessionInfo := r.Group("/api/space/v1/info")
	v1SessionInfo.Use(middlewares.Middleware())
	{
		v1SessionInfo.GET("/session", apiControllerV1.SessionInfo)
	}

	v1TechnicalIndicators := r.Group("/api/space/v1/technicalIndicators")
	v1TechnicalIndicators.Use(middlewares.AuthCombinedMiddleware())
	{
		v1TechnicalIndicators.POST("/technicalIndicatorsValues", apiControllerV1.TechnicalIndicatorsValues)
	}

	v2TechnicalIndicators := r.Group("/api/space/v2/technicalIndicators")
	v2TechnicalIndicators.Use(middlewares.Middleware())
	{
		v2TechnicalIndicators.POST("/getSMA", apiControllerV2.GetSMA)
		v2TechnicalIndicators.POST("/getEMA", apiControllerV2.GetEMA)
		v2TechnicalIndicators.POST("/getHullMA", apiControllerV2.GetHullMA)
		v2TechnicalIndicators.POST("/getVWMA", apiControllerV2.GetVWMA)
		v2TechnicalIndicators.POST("/getRSI", apiControllerV2.GetRSI)
		v2TechnicalIndicators.POST("/getCCI", apiControllerV2.GetCCI)
		v2TechnicalIndicators.POST("/getADX", apiControllerV2.GetADX)
		v2TechnicalIndicators.POST("/getMACD", apiControllerV2.GetMACD)
		v2TechnicalIndicators.POST("/getStochastic", apiControllerV2.GetStochastic)
		v2TechnicalIndicators.POST("/getIchimokuBaseLine", apiControllerV2.GetIchimokuBaseLine)
		v2TechnicalIndicators.POST("/getAwesomeOscillator", apiControllerV2.GetAwesomeOscillator)
		v2TechnicalIndicators.POST("/getMomentum", apiControllerV2.GetMomentum)
		v2TechnicalIndicators.POST("/getStochRSIFast", apiControllerV2.GetStochRSIFast)
		v2TechnicalIndicators.POST("/getWilliamsRange", apiControllerV2.GetWilliamsRange)
		v2TechnicalIndicators.POST("/getUltimateOscillator", apiControllerV2.GetUltimateOscillator)
		v2TechnicalIndicators.POST("/getAllTechnicalIndicators", apiControllerV2.GetAllTechnicalIndicators)
	}

	v1Notifcations := r.Group("/api/space/v1/notifications")
	v1Notifcations.Use(middlewares.Middleware())
	{
		v1Notifcations.POST("/adminMessages", apiControllerV1.FetchAdminMessages)
		v1Notifcations.GET("/notificationUpdates", apiControllerV1.NotificationUpdates)
	}

	v2Notifcations := r.Group("/api/space/v2/notifications")
	v2Notifcations.Use(middlewares.Middleware())
	{
		v2Notifcations.GET("/adminMessages", apiControllerV2.FetchAdminMessages)
	}

	v1BondEtf := r.Group("/api/space/v1/bondEtf")
	v1BondEtf.Use(middlewares.UserAuthentication())
	{
		v1BondEtf.POST("/fetchBondData", apiControllerV1.FetchBondData)
	}

	v2Group := r.Group("/api/space/v2/contractdetails")
	v2Group.Use()
	{
		v2Group.GET("/searchScrip", apiControllerV2.SearchScript)
	}

	v1Warning := r.Group("/api/space/v1/warning")
	v1BondEtf.Use(middlewares.UserAuthentication())
	{
		v1Warning.POST("/nudgeAlert", apiControllerV1.NudgeAlert)
	}

	v3Cmots := r.Group("/api/space/v3/cmots")
	v3Cmots.Use(middlewares.AuthCombinedMiddleware())
	{
		v3Cmots.POST("/fetchFinancials", apiControllerV3.FetchFinancialsV3)
	}

	v1Testing := r.Group("/api/space/v1/test")
	// v1Testing.Use(middlewares.TLAuth())
	{
		v1Testing.POST("/testingRes", apiControllerV1.TestingApi)
	}

	v1Backoffice := r.Group("/api/space/v1/shilpi")
	v1Backoffice.Use(middlewares.UserAuthentication())
	{
		v1Backoffice.POST("/tradeConfirmationDateRange", apiControllerV1.TradeConfirmationDateRange)
		v1Backoffice.POST("/getBillDetailsCdsl", apiControllerV1.GetBillDetailsCdsl)
		v1Backoffice.POST("/longTermShortTerm", apiControllerV1.LongTermShortTerm)
		v1Backoffice.POST("/fetchProfile", apiControllerV1.FetchProfile)
		v1Backoffice.POST("/tradeConfirmationOnDate", apiControllerV1.TradeConfirmationOnDate)
		v1Backoffice.POST("/openPositions", apiControllerV1.OpenPositions)
		v1Backoffice.POST("/getHolding", apiControllerV1.GetHolding)
		v1Backoffice.POST("/getMarginOnDate", apiControllerV1.GetMarginOnDate)
		v1Backoffice.POST("/financialLedgerBalanceOnDate", apiControllerV1.FinancialLedgerBalanceOnDate)
		v1Backoffice.POST("/getFinancial", apiControllerV1.GetFinancial)
	}

	v1Edis := r.Group("/api/space/v1/edis")
	v1Edis.Use(middlewares.Middleware())
	{
		v1Edis.POST("/edisRequest", apiControllerV1.EdisRequest)
		v1Edis.POST("/generateTpin", apiControllerV1.GenerateTpin)

	}

	v1Epledge := r.Group("/api/space/v1/epledge")
	v1Epledge.Use(middlewares.Middleware())
	{
		v1Epledge.POST("/epledgeRequest", apiControllerV1.EpledgeRequest)
		v1Epledge.POST("/unpledge", apiControllerV1.UnpledgeRequest)
		v1Epledge.POST("/mtfEpledgeRequest", apiControllerV1.MTFEpledgeRequest)
		v1Epledge.GET("/getPledgeList", apiControllerV1.GetPledgeList)
		v1Epledge.POST("/getCTDQuantityList", apiControllerV1.GetCTDQuantityList)
		v1Epledge.GET("/getPledgeTransactions", apiControllerV1.GetPledgeTransactions)
		v1Epledge.POST("/mtfCtd", apiControllerV1.MTFCTD)
	}

	v2Alerts := r.Group("/api/space/v2/alerts")
	v2Alerts.Use(middlewares.Middleware())
	{
		v2Alerts.PUT("/editAlerts", apiControllerV2.EditAlerts)
		v2Alerts.PUT("/pauseAlerts", apiControllerV2.PauseAlerts)
		v2Alerts.DELETE("/deleteAlerts", apiControllerV2.DeleteAlerts)
		v2Alerts.GET("/getAlerts", apiControllerV2.GetAlerts)
	}

	v4Cmots := r.Group("/api/space/v4/cmots")
	v4Cmots.Use(middlewares.AuthCombinedMiddleware())
	{
		v4Cmots.POST("/fetchFinancials", apiControllerV4.FetchFinancialsV4)
	}

	v1Reports := r.Group("/api/space/v1/reports")
	v1Reports.Use(middlewares.TLAuthReports())
	{
		v1Reports.GET("/viewDPCharges", apiControllerV1.ViewDPCharges)
		v1Reports.GET("/downloadDPCharges", apiControllerV1.DownloadDPCharges)
		v1Reports.GET("/viewTradebook", apiControllerV1.ViewTradebook)
		v1Reports.GET("/downloadTradebook", apiControllerV1.DownloadTradebook)
		v1Reports.GET("/viewLedger", apiControllerV1.ViewLedger)
		v1Reports.GET("/downloadLedger", apiControllerV1.DownloadLedger)
		v1Reports.GET("/viewOpenPosition", apiControllerV1.ViewOpenPosition)
		v1Reports.GET("/downloadOpenPosition", apiControllerV1.DownloadOpenPosition)
		v1Reports.GET("/viewFnoPnl", apiControllerV1.ViewFnoPnl)
		v1Reports.GET("/downloadFnoPnl", apiControllerV1.DownloadFnoPnl)
		v1Reports.GET("/viewHoldingFinancial", apiControllerV1.ViewHoldingFinancial)
		v1Reports.GET("/downloadHoldingFinancial", apiControllerV1.DownloadHoldingFinancial)
		v1Reports.GET("/sendEmailLedger", apiControllerV1.SendEmailLedger)
		v1Reports.GET("/viewCommodityTradebook", apiControllerV1.ViewCommodityTradebook)
		v1Reports.GET("/downloadCommodityTradebook", apiControllerV1.DownloadCommodityTradebook)
		v1Reports.GET("/sendEmailCommodityTradebook", apiControllerV1.SendEmailCommodityTradebook)
		v1Reports.GET("/viewFnoTradebook", apiControllerV1.ViewFnoTradebook)
		v1Reports.GET("/downloadFnoTradebook", apiControllerV1.DownloadFnoTradebook)
		v1Reports.GET("/sendEmailFnoTradebook", apiControllerV1.SendEmailFnoTradebook)
		v1Reports.GET("/sendEmailDPCharges", apiControllerV1.SendEmailDPCharges)
		v1Reports.GET("/sendEmailHoldingFinancial", apiControllerV1.SendEmailHoldingFinancial)
	}

	v2pockets := r.Group("/api/space/v2/pockets")
	v2pockets.Use(middlewares.UserAuthentication())
	{
		v2pockets.POST("/buyPocket", apiControllerV2.BuyPocket)
		v2pockets.POST("/exitPocket", apiControllerV2.ExitPocket)
		v2pockets.POST("/fetchPocketPortfolio", apiControllerV2.FetchPocketPortfolio)
	}

	v1Upi := r.Group("/api/space/v1/upi")
	v1Upi.Use(middlewares.UserAuthentication())
	{
		v1Upi.POST("/setUpiPreference", apiControllerV1.SetUpiPreference)
		v1Upi.GET("/fetchUpiPreference", apiControllerV1.FetchUpiPreference)
		v1Upi.DELETE("/deleteUpiPreference", apiControllerV1.DeleteUpiPreference)
	}

	v1oauth2Login := r.Group("/api/space/v1/authapis")
	v1oauth2Login.Use(middlewares.Middleware())
	{
		v1oauth2Login.GET("/handleAuthCode", apiControllerV1.HandleAuthCode)
		v1oauth2Login.POST("/getAccessToken", apiControllerV1.GetAccessToken)
	}

	v2oauth2Login := r.Group("/api/space/v2/authapis")
	v2oauth2Login.Use(middlewares.Middleware())
	{
		v2oauth2Login.POST("/getAccessToken", apiControllerV2.GetAccessTokenV2)
	}

	v1oauth2Login.Use(middlewares.UserAuthentication())
	{
		v1oauth2Login.POST("/createApp", apiControllerV1.CreateApp)
		v1oauth2Login.GET("/fetchApps/:clientId", apiControllerV1.FetchApps)
		v1oauth2Login.DELETE("/deleteApp/:appId", apiControllerV1.DeleteApp)
	}

	v3WatchList := r.Group("/api/space/v3/watchlist")
	v3WatchList.Use(middlewares.UserAuthentication())
	{
		v3WatchList.POST("/addStockToWatchList", apiControllerV3.AddStockToWatchList)
		v3WatchList.GET("/fetchWatchList", apiControllerV3.FetchWatchList)
		v3WatchList.DELETE("/deleteStockInWatchList", apiControllerV3.DeleteStockInWatchList)
		v3WatchList.POST("/arrangeStocksWatchList", apiControllerV3.ArrangeStocksWatchList)
		v3WatchList.DELETE("/deleteStockInWatchListUpdated", apiControllerV3.DeleteStockInWatchListUpdated)
	}

	freshdesk := r.Group("/api/space/v1/support")

	tickets := freshdesk.Group("/ticket")
	tickets.Use(middlewares.UserAuthentication())
	{
		tickets.POST("/create", apiControllerV1.CreateFreshdeskTicket)
	}

	// pocket-actions
	v3Pockets := r.Group("/api/space/v3/pockets")
	v3Pockets.Use(middlewares.UserAuthentication())
	{
		v3Pockets.POST("/buyPocket", apiControllerV3.BuyPocketV3)
		v3Pockets.GET("/fetchAllPockets", apiControllerV3.FetchAllPocketsV3)
		v3Pockets.POST("/fetchPocketPortfolio", apiControllerV3.FetchPocketPortfolioV3)
		v3Pockets.GET("/fetchUsersPockets", apiControllerV3.FetchUsersPockets)
		v3Pockets.POST("/sellPocket", apiControllerV3.SellPocketV3)
		v3Pockets.POST("/exitPocket", apiControllerV3.ExitPocketV3)
		v3Pockets.GET("/getPocketDetails", apiControllerV3.GetPocketDetails)

		v3PocketActions := v3Pockets.Group("/action")
		{
			v3PocketActions.POST("/checkActionRequired", apiControllerV3.CheckActionRequired)
			v3PocketActions.POST("/adjustStocks", apiControllerV3.AdjustStocksForPocket)
		}

	}

	v1BondsDetails := r.Group("/api/space/v1/bondsDetails")
	v1BondsDetails.Use(middlewares.UserAuthentication())
	{
		v1BondsDetails.GET("/fetchBondDataByIsin", apiControllerV1.FetchBondDataByIsin)
	}

	stockSip := r.Group("/api/space/v1/sip")
	stockSip.Use(middlewares.Middleware())
	{
		stockSip.GET("/fetchStockSipOrder", apiControllerV1.FetchStockSips)
		stockSip.POST("/placeSipOrder", apiControllerV1.PlaceSipOrder)
		stockSip.DELETE("/deleteSipOrder/:clientId/:sipId", apiControllerV1.DeleteSipOrder)
		stockSip.PUT("/modifySipOrder", apiControllerV1.ModifySipOrder)
		stockSip.PUT("/updateSipStatus", apiControllerV1.UpdateSipStatus)
	}

	//API route for version 2+
	// v2 := r.Group("/api/space/v2")

	// v2.POST("user-list", apiControllerV2.UserList)

	// Add v3 search route
	v3Search := r.Group("/api/space/v3/contractdetails")
	v3Search.Use(middlewares.UserAuthentication())
	{
		v3Search.GET("/searchScrip", apiControllerV3.SearchScrip)
	}

	return r

}
