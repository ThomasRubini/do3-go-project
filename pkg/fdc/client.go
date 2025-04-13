package fdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://api.nal.usda.gov/fdc/v1"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

type SearchRequest struct {
	Query      string `json:"query"`
	DataType   []string `json:"dataType,omitempty"`
	PageSize   int    `json:"pageSize"`
	PageNumber int    `json:"pageNumber"`
}

type SearchResponse struct {
	Foods     []FoodItem `json:"foods"`
	TotalHits int        `json:"totalHits"`
}

type FoodItem struct {
	FdcId       int        `json:"fdcId"`
	Description string     `json:"description"`
	DataType    string     `json:"dataType"`
	Nutrients   []Nutrient `json:"foodNutrients"`
}

type Nutrient struct {
	Number   string  `json:"number"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	UnitName string  `json:"unitName"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func toStr(rc io.ReadCloser) string {
	defer rc.Close()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(rc)
	if err != nil {
		fmt.Println("error reading response body:", err)
		return ""
	}

	return buf.String()
}

func (c *Client) SearchFoods(query string, dataType string) (*SearchResponse, error) {
	reqBody := SearchRequest{
		Query:      query,
		DataType:   []string{dataType},
		PageSize:   25,
		PageNumber: 1,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/foods/search?api_key=%s", baseURL, c.apiKey), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s. Body: %s", resp.Status, toStr(resp.Body))
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &searchResp, nil
}

func (c *Client) GetFoodDetails(fdcId int) (*FoodItem, error) {
	url := fmt.Sprintf("%s/food/%d?api_key=%s", baseURL, fdcId, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %s. Body: %s", resp.Status, toStr(resp.Body))
	}

	var food FoodItem
	if err := json.NewDecoder(resp.Body).Decode(&food); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &food, nil
}
