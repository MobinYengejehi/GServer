package Crawler

import (
	"GServer/Logger"
	"GServer/Movie"
	"context"
)

type SearchResultFunction func(*Client) []*Movie.MovieDetails
type ServiceTotalLengthFunction func(*Client) float64

type Client struct {
	Name string

	Rows int32

	StartPage   int32
	CurrentPage int32

	TotalMovies int64

	Started bool

	GetSearchResult    SearchResultFunction
	GetTotalMovieCount ServiceTotalLengthFunction

	Context context.Context

	ServiceClient interface{}
}

func (this *Client) Stop() {
	if !this.Started {
		return
	}

	this.Started = false

	this.StartPage = 0
	this.CurrentPage = 0
}

func (this *Client) Start() {
	if this.Started {
		return
	}

	this.Started = true

	Logger.INFO("crawler started : ", this.Name)
}

func NewClient(ctx context.Context, name string, rows int32, startPage int32) *Client {
	var client *Client = new(Client)

	client.Name = name

	client.Rows = rows

	client.StartPage = startPage
	client.CurrentPage = 0

	client.TotalMovies = 0

	client.Started = false

	client.GetSearchResult = func(c *Client) []*Movie.MovieDetails { return []*Movie.MovieDetails{} }
	client.GetTotalMovieCount = func(c *Client) float64 { return 0 }

	client.Context = ctx

	client.ServiceClient = nil

	return client
}
