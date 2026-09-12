package coingecko

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Struct for working with API
type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://api.coingecko.com/api/v3",
	}
}

// Function GetPrice returns coin price in specified currency
func (c *Client) GetPrice(coin string, vsCurrency string) (float64, error) {
	coin = strings.ToLower(strings.TrimSpace(coin))
	vsCurrency = strings.ToLower(strings.TrimSpace(vsCurrency))

	// A little mapping for convenience
	coinMap := map[string]string{
		"btc":  "bitcoin",
		"eth":  "ethereum",
		"sol":  "solana",
		"ton":  "the-open-network",
		"usdt": "tether",
	}

	if fullID, exists := coinMap[coin]; exists {
		coin = fullID
	}

	// URL request generation
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=%s", c.baseURL, coin, vsCurrency)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("Sending HTTP-request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	// Answer structure CoinGecko: {"bitcoin": {"usd": 65000.50}}
	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("JSON parsing error: %w", err)
	}

	price, ok := result[coin][vsCurrency]
	if !ok {
		return 0, fmt.Errorf("Coin '%s' or currency '%s' not found", coin, vsCurrency)
	}

	return price, nil
}
