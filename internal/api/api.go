package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var APIURL = "https://min-api.cryptocompare.com/data/pricemulti?fsyms=BTC,ETH&tsyms=USD,EUR,RUB&api_key=%s"

type CryptoClient struct {
	apiKey string
}

func NewCryptoCompareClient(apiKey string) *CryptoClient {
	return &CryptoClient{apiKey: apiKey}
}

func (c *CryptoClient) GetCryptoRates() (map[string]map[string]float64, error) {
	url := fmt.Sprintf(APIURL, c.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from API: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var rawCurrencies map[string]map[string]float64
	err = json.Unmarshal(body, &rawCurrencies)
	if err != nil {
		return nil, err
	}
	return rawCurrencies, nil
}
