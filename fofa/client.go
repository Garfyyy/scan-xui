package fofa

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Garfyyy/scan-xui/utils"
)

var FOFA_API_URL string = "https://fofa.info/api/v1/search/all"

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

type SearchParams struct {
	Query  string
	Size   int
	Page   int
	Fields string
}

type FOFAResponse struct {
	Error   bool       `json:"error"`
	Size    int        `json:"size"`
	Results [][]string `json:"results"`
}

func (c *Client) Search(params *SearchParams) (*FOFAResponse, error) {
	url := fmt.Sprintf("%s?key=%s&qbase64=%s&size=%d&page=%d&fields=%s&full=false", FOFA_API_URL, c.apiKey, utils.EncodeBase64(params.Query), params.Size, params.Page, params.Fields)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var fofaResp FOFAResponse
	if err := json.NewDecoder(resp.Body).Decode(&fofaResp); err != nil {
		return nil, err
	}

	if fofaResp.Error {
		return nil, fmt.Errorf("FOFA API error")
	}

	return &fofaResp, nil
}
