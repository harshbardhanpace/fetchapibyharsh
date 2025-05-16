package models

type SquareOffAllReq struct {
	ClientID string `json:"clientId"`
}

type SquareOffAllRes struct {
	SquareOffAll []PlaceOrderResponse `json:"squareOffAll"`
}
