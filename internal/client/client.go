package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alevern/pokedexapi/internal/cache"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

type Client struct {
	httpClient http.Client
	cache      cache.Cache
}

func NewClient(timeout time.Duration) Client {
	pokeCache := cache.NewCache(20 * time.Second)
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokeCache,
	}
}

type RespShallowLocations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) ListLocations(pageUrl *string) (RespShallowLocations, error) {
	url := baseURL + "/location-area"
	if pageUrl != nil {
		url = *pageUrl
	}
	locationsResp := RespShallowLocations{}
	savedLocations, found := c.cache.Get(url)
	if found {
		err := json.Unmarshal(savedLocations, &locationsResp)
		if err != nil {
			return RespShallowLocations{}, err
		}
		return locationsResp, nil
	}

	if locations, found := c.loadLocations(url); found == true {
		fmt.Println("Cached!")
		return locations, nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespShallowLocations{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RespShallowLocations{}, err
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespShallowLocations{}, err
	}
	err = json.Unmarshal(dat, &locationsResp)
	if err != nil {
		return RespShallowLocations{}, err
	}
	c.cache.Add(url, dat)
	return locationsResp, nil
}

func (c *Client) loadLocations(url string) (RespShallowLocations, bool) {
	savedLocations, found := c.cache.Get(url)
	if !found {
		return RespShallowLocations{}, false
	}

	locations := RespShallowLocations{}
	err := json.Unmarshal(savedLocations, &locations)
	if err != nil {
		return RespShallowLocations{}, false
	}
	return locations, true

}
