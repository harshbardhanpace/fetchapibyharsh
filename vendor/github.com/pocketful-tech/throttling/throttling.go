package throttling

import (
	"bytes"
	"io"
	"time"
)

func (c *APIClient) MakeAPIRequest(method string, apiURL string, payload *bytes.Buffer, authorization string) ([]byte, error) {
	c.Throttler.Throttle()

	var resByte []byte

	resp, err := CallApi(method, apiURL, payload, authorization)

	if err != nil {
		return resByte, err
	}
	defer resp.Body.Close()

	// Read and print the response body (you can handle the response as needed)
	resByte, err = io.ReadAll(resp.Body)
	if err != nil {
		return resByte, err
	}

	return resByte, nil
}

func NewAPIThrottler(requestsPerSecond int, clientName, vendorName string) *APIThrottler {
	return &APIThrottler{
		RequestsPerSecond: requestsPerSecond,
		LastRequestTime:   time.Now(),
		ClientName:        clientName,
		VendorName:        vendorName,
		Count:             0,
	}
}
