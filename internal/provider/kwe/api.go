package kwe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	cfg        *config
	httpClient *http.Client
}

func NewClient() (*Client, error) {
	cfg, err := newConfig()
	if err != nil {
		return nil, err
	}

	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{},
	}, nil
}

type req struct {
	Language          string `json:"language"`
	NumberOfKeywords  int    `json:"number_of_keywords"`
	MaxNgramSize      int    `json:"max_ngram_size"`
	WindowSize        int    `json:"window_size"`
	DeduplicationAlgo string `json:"deduplication_algo"`
	Text              string `json:"text"`
}

type resNgram struct {
	Score float64 `json:"-"`
	Ngram string  `json:"ngram"`
}

func newReq(text string) io.Reader {
	r := req{
		Language:          "ru",
		NumberOfKeywords:  10,
		MaxNgramSize:      1,
		WindowSize:        1,
		DeduplicationAlgo: "seqm",
		Text:              text,
	}

	data, _ := json.Marshal(r)

	return bytes.NewReader(data)
}

func (c *Client) Extract(_ context.Context, text string) ([]string, error) {
	res, err := c.httpClient.Post(c.cfg.apiURL, "application/json", newReq(text))
	if err != nil {
		return nil, fmt.Errorf("post request to yake: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response from yake: %w", err)
	}

	var result []resNgram

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("parse response from yake: %w", err)
	}

	var keywords []string

	for _, v := range result {
		keywords = append(keywords, v.Ngram)
	}

	return keywords, nil
}
