package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (c *Client) GetPokemonInfos(name string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + name
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}
	pokemonResp := Pokemon{}
	savedPokemon, found := c.cache.Get(url)
	if found {
		err := json.Unmarshal(savedPokemon, &pokemonResp)
		if err != nil {
			return Pokemon{}, err
		}
		return pokemonResp, nil
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}
	err = json.Unmarshal(dat, &pokemonResp)
	if err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			fmt.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("err response: %+v", pokemonResp)
		return Pokemon{}, err
	}
	c.cache.Add(url, dat)
	return pokemonResp, nil
}
func (c *Client) ListPokemonsEncounters(location string) (RespLocationArea, error) {
	url := baseURL + "/location-area/" + location
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespLocationArea{}, err
	}
	encountersResp := RespLocationArea{}
	savedEncounters, found := c.cache.Get(url)
	if found {
		err := json.Unmarshal(savedEncounters, &encountersResp)
		if err != nil {
			return RespLocationArea{}, err
		}
		return encountersResp, nil
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RespLocationArea{}, err
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespLocationArea{}, err
	}
	err = json.Unmarshal(dat, &encountersResp)
	if err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			fmt.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("sakura response: %q", encountersResp)
		return RespLocationArea{}, err
	}
	c.cache.Add(url, dat)
	return encountersResp, nil
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
