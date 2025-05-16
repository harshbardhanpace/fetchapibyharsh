package tradelab

const (
	LOGINURL                   = "/api/v1/user/login"
	VERIFYTWOFA                = "/api/v1/user/twofa"
	SETTWOFA                   = "/api/v1/choose_twofa"
	FORGOTPASSWORD             = "/api/v1/user/password/forgot"
	SETPASSWORD                = "/api/v1/user/password"
	FORGETRESETTWOFA           = "/api/v1/user/force_reset_twofa"
	PLACEORDERURL              = "/api/v1/orders"
	PENDINGORDERURL            = "/api/v1/orders"
	COMPLETEDORDERURL          = "/api/v1/orders"
	CANCELGTTURL               = "/api/v1/event/gtt"
	CONDATIONALORDERSURL       = "/api/v1/orders/kart"
	TRADEURL                   = "/api/v1/trades"
	FETCHDEMATHOLDINGSURL      = "/api/v1/holdings"
	CONVERTPOSITIONSURL        = "/api/v1/position/convert"
	GETPOSITIONURL             = "/api/v1/positions"
	FETCHOPTIONCHAINURL        = "/api/v1/optionchain"
	FETCHFUTURESCHAINURL       = "/api/v1/futureschain/NFO"
	PROFILEURL                 = "/api/v1/user/profile"
	FETCHFUNDSURL              = "/api/v2/funds/view"
	CANCELPAYOUTFUNDSURL       = "/api/v1/funds/transactions"
	SEARCHSCRIPURL             = "/api/v1/search"
	SCRIPINFOURL               = "/api/v1/contract/"
	ORDERHISTORY               = "/api/v1/order/"
	GTTURL                     = "/api/v1/event/gtt"
	GETIPOURL                  = "/api/v1/tradeipo/get-all-ipo"
	PLACEIPOORDER              = "/api/v1/tradeipo/place-order"
	FETCHIPOORDER              = "/api/v1/tradeipo/fetch-order"
	CANCELIPOORDER             = "/api/v1/tradeipo/cancel-order"
	MARGINCALCULATION          = "/api/v1/calc"
	GainerLoserNse             = "/api/v1/screeners/gainerslosers"
	Screeners                  = "/api/v1/screeners"
	Charts                     = "/api/v1/charts"
	LastTradedPrice            = "/api/v1/marketdata"
	BasketURL                  = "/api/v1/basket"
	BasketInstrumentURL        = "/api/v1/basket/order"
	ExecuteBasket              = "/api/v1/orders/kart"
	Payout                     = "/api/v1/pg/atom/payout"
	ClientTransactions         = "/api/v1/funds/transactions"
	ALERTSURL                  = "/api/v1/alerts"
	UPDATEBASKETEXECUTIONSTATE = "/api/v1/update_basket_execution_state"
	SESSIONINFOURL             = "/api/v1/exchange/session/info"
	LOGINV2URL                 = "/api/v3/user/login"
	TWOFAV2URL                 = "/api/v3/user/twofa"
	SETTOTPV2URL               = "/api/v3/setup_totp"
	CHOOSETWOFAV2URL           = "/api/v3/choose_twofa"
	RETURNONINVESTMENTURL      = "/api/v1/screeners/roi"
	ALLBANKACCOUNTURL          = "/api/v1/user/bank/accounts"
	ADMINMESSAGE               = "/api/v1/updates"
	ADMINMESSAGELATEST         = "/api/v1/latest/updates"
	FORGETTOTPV2URL            = "/api/v3/user/forgot_totp"
	VALIDATELOGINOTP           = "/api/v3/user/login/otp"
	BIOMETRIC                  = "/api/v3/user/biometric"
	LOGOUTURL                  = "/api/v1/user/logout"
	UNBLOCKUSER                = "/api/v1/management/user/unblock_user"
	CREATEAPP                  = "/api/v1/app"
	FETCHAPPS                  = "/api/v1/apps"
	GENERATEACCESSTOKEN        = "/oauth2/token"
	ACCOUNTFREEZE              = "/api/v1/user/account/freeze"
	MTFEPLEDGE                 = "/api/v1/mpi/mtf/epledge"
	MTFPLEDGELIST              = "/api/v1/mtf/pledge/list"
	CTDQUANTITYLIST            = "/api/v1/mtf/ctd/list"
	SIPUPDATESTATUSURL         = "/api/v1/events/sip"
	SIPURL                     = "/api/v1/event/sip"
	PLACEICEBERGORDERURL       = "/api/v1/event/iceberg"
	MTFCTD                     = "/api/v1/mtf/ctd/multiple/isin"
)

const (
	TLERROR = "error"
)

const (
	USERORDERIDKEY = "userOrderId"
)

const (
	EDISREQ     = "/api/v1/edis/instrument_details"
	TPINREQ     = "/api/v1/edis/bopin_generation"
	EPLEDGEREQ  = "/api/v1/epledge"
	UNPLEDGEURL = "/api/v1/pledge"

	SIPREVOCATION = "/revocation"
)

const (
	UNPLEDGE = "UNPLEDGE"
	CAPITAL  = "CAPITAL"
	Capital  = "Capital"
)

const (
	FetchIpoOrderDataRedisKeyPrefix = "FetchIpoOrderDataRedisKeyPrefix"
	FetchAllIPODataRedisKey         = "FetchAllIPODataRedisKey"
)

const (
	AccountFrozen           = "You can't login as account is in freeze state."
	RequestDataNotFound     = "Unable to read request Data."
	InvalidToken            = "Invalid UnderLyingtoken"
	PledgeHours             = "Pledge is only available between 8am (08:00) to 3:30pm (15:30)"
	PledgeTimeError         = "[400] - pledge not allowed at this time"
	PledgeTimeErrorNew      = "[95010] - Pledge not allowed at this time"
	PledgeTimeErrorOnlyText = "pledge not allowed at this time"
)
