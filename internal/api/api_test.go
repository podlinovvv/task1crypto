package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"task1crypto/internal/api"
)

func TestGetCryptoRates(t *testing.T) {
	testCases := []struct {
		name           string
		responseStatus int
		responseBody   string
		expectedError  bool
	}{
		{
			name:           "successful_request",
			responseStatus: http.StatusOK,
			responseBody: `{
				"BTC": {"USD": 1000, "EUR": 900, "RUB": 60000},
				"ETH": {"USD": 500, "EUR": 450, "RUB": 30000}
			}`,
			expectedError: false,
		},
		{
			name:           "unsuccessful_request",
			responseStatus: http.StatusInternalServerError,
			responseBody:   "Internal Server Error",
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.responseStatus)
				_, _ = w.Write([]byte(tc.responseBody))
			}))
			defer mockServer.Close()

			client := api.NewCryptoCompareClient("test_api_key")

			// Replace the default API URL with the mock server URL
			api.APIURL = mockServer.URL + "/data/pricemulti?fsyms=BTC,ETH&tsyms=USD,EUR,RUB&api_key=%s"

			rates, err := client.GetCryptoRates()

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rates)
				assert.Equal(t, 2, len(rates))
			}
		})
	}
}
