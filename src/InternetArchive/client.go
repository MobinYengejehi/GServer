package InternetArchive

import (
	"GServer/Defaults"
	"GServer/Logger"
	"GServer/Movie"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type JsonDictionary map[string]any

type Client struct {
	BaseURL string

	AdvancedSearchEndpoint string

	TorrentURLFormat string

	Context context.Context
	Timeout time.Duration
}

func (this *Client) fetch(url *url.URL, method string, payload []byte) (JsonDictionary, error) {
	var buffer *bytes.Buffer = bytes.NewBufferString("")

	if payload != nil {
		buffer = bytes.NewBuffer(payload)
	}

	requestContext, requestContextClose := context.WithTimeout(this.Context, this.Timeout)

	defer requestContextClose()

	request, err := http.NewRequestWithContext(requestContext, method, url.String(), buffer)

	if err != nil {
		return nil, err
	}

	var httpClient *http.Client = &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	if len(bodyBytes) < 1 {
		return nil, errors.New("Server didn't send any responses")
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(bodyBytes, &responseData)

	if err != nil {
		return nil, err
	}

	responseHeaderData, exists := responseData["responseHeader"]

	if !exists {
		return nil, errors.New("Response `responseHeader` dosn't exist")
	}

	responseHeader, ok := responseHeaderData.(map[string]interface{})

	if !ok {
		return nil, errors.New("Response `responseHeader` is not an object")
	}

	statusData, exists := responseHeader["status"]

	if !exists {
		return nil, errors.New("Response `status` doesn't exist")
	}

	status, ok := statusData.(float64)

	if !ok {
		return nil, errors.New("Invalid response `status`")
	}

	if status != 0 {
		return nil, errors.New("Request failed. Returned `status` code is '" + fmt.Sprintf("%.0g", status) + "'")
	}

	dataData, exists := responseData["response"]

	if !exists {
		return nil, errors.New("Response data `response` doesn't exist")
	}

	data, ok := dataData.(map[string]interface{})

	if !ok {
		return nil, errors.New("Response `response` field is not an object")
	}

	return data, nil
}

func setMovieDetail[T any](detail *T, jsonData *map[string]interface{}, field string) {
	fieldValue, exists := (*jsonData)[field]

	if !exists {
		return
	}

	value, ok := fieldValue.(T)

	if !ok {
		return
	}

	*detail = value
}

func parseMovieDetailsFromJsonData(details *Movie.MovieDetails, jsonData *map[string]interface{}, client *Client) {
	setMovieDetail(&details.SpecialIdentifier, jsonData, "identifier")

	setMovieDetail(&details.Size, jsonData, "item_size")

	setMovieDetail(&details.DescriptionFull, jsonData, "description")

	setMovieDetail(&details.Title, jsonData, "title")

	setMovieDetail(&details.Language, jsonData, "language")

	setMovieDetail(&details.DateUploaded, jsonData, "date")

	if len(details.SpecialIdentifier) < 1 {
		return
	}

	var torrent *Movie.MovieTorrentInfo = Movie.NewMovieTorrentInfo()

	torrent.URL = fmt.Sprintf(client.TorrentURLFormat, details.SpecialIdentifier, details.SpecialIdentifier)

	err := Movie.ParseTorrentFromUrl(client.Context, torrent.URL, torrent)

	if err != nil {
		Logger.WARN("Failed to parse torrent file. [URL: " + torrent.URL + ", Message: " + err.Error() + "]")
		return
	}

	details.Torrents = append(details.Torrents, torrent)
}

func (this *Client) Search(params *SearchParameters) ([]*Movie.MovieDetails, error, float64, float64) {
	url, err := url.Parse(this.AdvancedSearchEndpoint)

	var queryParams *SearchParameters = NewSearchParameters("")

	if params != nil {
		queryParams = params
	}

	url.RawQuery = ConvertSearchParametersToURLParams(queryParams).Encode()

	if err != nil {
		return nil, err, 0, 0
	}

	responseJsonData, err := this.fetch(url, "GET", nil)

	if err != nil {
		return nil, err, 0, 0
	}

	var movieCount float64 = 0
	var start float64 = 0

	movieCountData, exists := responseJsonData["numFound"]

	if exists {
		movieCount = movieCountData.(float64)
	}

	startData, exists := responseJsonData["start"]

	if exists {
		start = startData.(float64)
	}

	moviesListData, exists := responseJsonData["docs"]

	if !exists {
		return nil, errors.New("`docs` field doesn't exist"), 0, 0
	}

	moviesList, ok := moviesListData.([]interface{})

	if !ok {
		return nil, errors.New("`docs` field must be an array"), 0, 0
	}

	var moviesListResult []*Movie.MovieDetails = make([]*Movie.MovieDetails, 0)

	for _, value := range moviesList {
		item, ok := value.(map[string]interface{})

		if !ok {
			continue
		}

		var details *Movie.MovieDetails = Movie.NewMovieDetails()

		parseMovieDetailsFromJsonData(details, &item, this)

		moviesListResult = append(moviesListResult, details)
	}

	return moviesListResult, nil, movieCount, start
}

func (this *Client) GetMovieList(params *SearchParameters, extra string) ([]*Movie.MovieDetails, error, float64, float64) {
	params.Query = "mediatype:(movies) AND (subject:\"movie\" OR subject:\"serial\" OR subject:\"animation\" OR subject:\"cartoon\" OR subject:\"anime\")"

	if len(extra) > 0 {
		params.Query += " " + extra
	}

	return this.Search(params)
}

func NewClient(ctx context.Context, timeout time.Duration) *Client {
	var client *Client = new(Client)

	var err error

	client.BaseURL = Defaults.INTERNET_ARCHIVE_BASE_URL
	client.AdvancedSearchEndpoint, err = url.JoinPath(client.BaseURL, Defaults.INTERNET_ARCHIVE_ADVANCED_SEARCH_ENDPOINT)
	client.TorrentURLFormat = Defaults.INTERNET_ARCHIVE_TORRENT_URL_FORMAT

	client.Context = ctx
	client.Timeout = timeout

	if err != nil {
		return nil
	}

	return client
}

func NewClientWithCustomURL(ctx context.Context, timeout time.Duration, baseURL string) *Client {
	var client *Client = NewClient(ctx, timeout)

	var err error

	client.BaseURL = baseURL
	client.AdvancedSearchEndpoint, err = url.JoinPath(client.BaseURL, Defaults.INTERNET_ARCHIVE_ADVANCED_SEARCH_ENDPOINT)
	client.TorrentURLFormat = Defaults.INTERNET_ARCHIVE_TORRENT_URL_FORMAT

	if err != nil {
		return nil
	}

	return client
}
