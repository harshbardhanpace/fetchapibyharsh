package constants

// Error Keys
const (
	InternalServerError          = "P11001"
	InvalidParameters            = "P11002"
	MethodNotAllowed             = "P11003"
	InvalidIp                    = "P11004"
	USERIDINVALID                = "P11006"
	UserBlocked                  = "P11007"
	UserBlockedForTrading        = "P11008"
	MobileOrEmailAlreadyExists   = "P11010"
	InvalidUserId                = "P11011"
	InvalidToken                 = "P11012"
	TokenMissing                 = "P11013"
	IdDoesNotExists              = "P11014"
	TipDoesNotExists             = "P11015"
	IdAlreadyExists              = "P11016"
	InvalidRequest               = "P11017"
	InvalidFrontPart             = "P11018"
	DetailsDoesNotExsists        = "P11019"
	AdminInvalidCreds            = "P11020"
	InvalidUserVideo             = "P11021"
	InvalidUserSelfie            = "P11022"
	InvalidMobileNo              = "P11030"
	InvalidEmailId               = "P11031"
	InvalidEmailOtp              = "P11032"
	InvalidMobileOtp             = "P11033"
	InvalidDeviceType            = "P11034"
	InvalidType                  = "P11040"
	InvalidClient                = "P11041"
	PocketAlreadyExists          = "P11042"
	PocketDoesNotExists          = "P11043"
	InvalidUserIdOrPass          = "P11044"
	CollectionAlreadyExists      = "P11045"
	CollectionDoesNotExists      = "P11046"
	WatchListsAlreadyExists      = "P11047"
	WatchListsDoesNotExists      = "P11048"
	PinsDoesNotExists            = "P11049"
	PinsCapacityFull             = "P11050"
	TLChartDataFetchFailed       = "P11051"
	AuthenticationFailed         = "P11052"
	EmptyIsin                    = "P11053"
	WatchListsStockAlreadyExists = "P11054"
	InvalidExchange              = "P11055"
	DecodingHeaderError          = "P11056"
	EmptyProfileResponse         = "P11057"
	EmptyIpAddress               = "P11058"
	InvalidDeviceId              = "P11059"
	TokenExpired                 = "P11060"
	UpiDontExist                 = "P11061"
	DuplicateUpi                 = "P11062"
	InvalidDate                  = "P11063"
	InvalidIsin                  = "P11064"
	MismatchAuthClient           = "P11065"
	InvalidDisplayName           = "P11066"
	InvalidPayoutRequest         = "P11067"
	PinsSizeExceed               = "P11068"
	PinsIndexInvalid             = "P11069"
	InvalidNameLength            = "P11070"
	ExistPayoutRequest           = "P11071"
	EmptyCredentials             = "P11072"
	InvalidOtp                   = "P11073"
	AccountFreezeInvalidRequest  = "P11074"
	InvalidHeader                = "P11075"
	AppDoesNotExists             = "P11076"
	LotSizeExceeds               = "P11077"
	InvalidPage                  = "P11078"
)

// Errors Code Map
var ErrorCodeMap = map[string]string{
	"P11001": "INTERNAL SERVER ERROR",
	"P11002": "INVALID PARAMETERS",
	"P11003": "METHOD NOT ALLOWED",
	"P11004": "INVALID IP",
	"P11006": "YOUR USER ID IS INVALID",
	"P11007": "USER IS BLOCKED",
	"P11008": "USER IS BLOCKED FOR TRADING",
	"P11010": "MOBILE OR EMAIL ALREADY EXISTS",
	"P11011": "INVALID USERID",
	"P11012": "INVALID TOKEN",
	"P11013": "TOKEN MISSING",
	"P11014": "ID DOES NOT EXISTS",
	"P11015": "TIP DOES NOT EXISTS",
	"P11016": "ID ALREADY EXISTS",
	"P11017": "INVALID REQUEST",
	"P11018": "INVALID FRONT PART",
	"P11019": "DETAILS DOES NOT EXISTS",
	"P11020": "ADMIN INVALID CREDS",
	"P11021": "USER VIDEO IS INVALID",
	"P11022": "USER SELFIE IS INVALID",
	"P11030": "INVALID MOBILE NO",
	"P11031": "INVALID EMAIL ID",
	"P11032": "INVALID EMAIL OTP",
	"P11033": "INVALID MOBILE OTP",
	"P11034": "INVALID DEVICE TYPE",
	"P11040": "INVALID TYPE",
	"P11041": "INVALID CLIENT",
	"P11042": "POCKET ALREADY EXISTS",
	"P11043": "POCKET DOES NOT EXISTS",
	"P11044": "INVALID USERID OR PASSWORD",
	"P11045": "COLLECTION ALREADY EXISTS",
	"P11046": "COLLECTION DOES NOT EXISTS",
	"P11047": "WATCHLISTS ALREADY EXISTS",
	"P11048": "WATCHLISTS DOES NOT EXISTS",
	"P11049": "PIN DOES NOT EXISTS",
	"P11050": "NO More Pins can be added! Capacity full!",
	"P11051": "Failed to fetch chart Data from TradeLab chart data api!",
	"P11052": "AUTHENTICATION FAILED",
	"P11053": "ISIN IS EMPTY",
	"P11054": "WatchLists Stock Already Exists",
	"P11055": "Invalid Exchange",
	"P11056": "Decoding Header Error",
	"P11057": "Empty Profile Response",
	"P11058": "Empty Ip Address",
	"P11059": "Invalid Device ID",
	"P11060": "Token Expired",
	"P11061": "Upi Don't Exist",
	"P11062": "Duplicate Upi",
	"P11063": "Invalid Date",
	"P11064": "Invalid Isin",
	"P11065": "Mismatch Auth Client",
	"P11066": "Invalid Display name",
	"P11067": "Your request is already proceed, try to request another payout",
	"P11068": "Pins size exceeds the available capacity",
	"P11069": "Invalid Pin Index",
	"P11070": "Name of Basket must be less than 30 characters",
	"P11071": "You cannot place another payout request while a previous request is processing.",
	"P11072": "empty client id and email id",
	"P11073": "Invalid Otp",
	"P11074": "Account Freeze Invalid Request",
	"P11075": "Invalid Header",
	"P11076": "App Does Not Exists",
	"P11077": "Lot Size exceeds the available quantity",
	"P11078": "Invalid Page",
}

const (
	SECRET_KEY = "eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTY0OTMzMzc4MiwiaWF0IjoxNjQ5MzMzNzgyfQ.XqmyGaUYYAToWylEWfD26EN52rBMivW7Zkt3O2u1cSo"
)
