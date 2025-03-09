// https://yts.mx/api

package YTS

import (
	"GServer/Config"
	"GServer/Defaults"
	"GServer/Logger"
	"GServer/Movie"
	"GServer/TaskManager"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
	"unsafe"
)

type JsonDictionary map[string]any

type Client struct {
	BaseURL string

	ListMoviesEndpoint       string
	MovieDetailsEndpoint     string
	MovieSuggestionsEndpoint string

	Context context.Context
	Timeout time.Duration

	HttpClient *http.Client
}

func (this *Client) fetch(url *url.URL, method string, payload []byte) (JsonDictionary, error) {
	var buffer *bytes.Buffer = bytes.NewBufferString("")

	if payload != nil {
		buffer = bytes.NewBuffer(payload)
	}

	requestContext, requestContextCancel := context.WithTimeout(this.Context, this.Timeout)

	defer requestContextCancel()

	request, err := http.NewRequestWithContext(requestContext, method, url.String(), buffer)

	if err != nil {
		return nil, err
	}

	response, err := this.HttpClient.Do(request)

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

	statusData, exists := responseData["status"]

	if !exists {
		return nil, errors.New("Response `status` doesn't exist")
	}

	status, ok := statusData.(string)

	if !ok {
		return nil, errors.New("Invalid `status` type. It must be an string")
	}

	statusMessageData, exists := responseData["status_message"]

	if !exists {
		statusMessageData = "UNKNOWN"
	}

	statusMessage, ok := statusMessageData.(string)

	if status != "ok" {
		return nil, errors.New("Request failed. Returned `status` is '" + status + "' [Message: " + statusMessage + "]")
	}

	dataData, exists := responseData["data"]

	if !exists {
		return nil, errors.New("Request failed. Server didn't sent any datas [Message: " + statusMessage + "]")
	}

	data, ok := dataData.(map[string]interface{})

	if !ok {
		return nil, errors.New("Couldn't parse `data`")
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

func parseMovieDetailsFromJsonData(client *Client, details *Movie.MovieDetails, jsonData *map[string]interface{}) {
	setMovieDetail(&details.Id, jsonData, "id")

	setMovieDetail(&details.URL, jsonData, "url")

	setMovieDetail(&details.IMDBCode, jsonData, "imdb_code")

	setMovieDetail(&details.Title, jsonData, "title")
	setMovieDetail(&details.TitleEnglish, jsonData, "title_english")
	setMovieDetail(&details.TitleLong, jsonData, "title_long")
	setMovieDetail(&details.Slug, jsonData, "slug")

	setMovieDetail(&details.Year, jsonData, "year")
	setMovieDetail(&details.Rating, jsonData, "rating")
	setMovieDetail(&details.Runtime, jsonData, "runtime")

	genresData, exists := (*jsonData)["genres"]

	if exists {
		genresList, ok := genresData.([]interface{})

		if ok {
			var genres []string

			for _, g := range genresList {
				if genre, ok := g.(string); ok {
					genres = append(genres, genre)
				}
			}

			details.Genres = genres
		}
	}

	setMovieDetail(&details.LikeCount, jsonData, "like_count")

	setMovieDetail(&details.Summery, jsonData, "summary")
	setMovieDetail(&details.DescriptionIntro, jsonData, "description_intro")
	setMovieDetail(&details.DescriptionFull, jsonData, "description_full")
	setMovieDetail(&details.Synopsis, jsonData, "synopsis")

	setMovieDetail(&details.YTTrailerCode, jsonData, "yt_trailer_code")

	setMovieDetail(&details.Language, jsonData, "language")

	setMovieDetail(&details.MPARating, jsonData, "mpa_rating")

	setMovieDetail(&details.BackgroundImage, jsonData, "background_image")
	setMovieDetail(&details.BackgroundImageOriginal, jsonData, "background_image_original")
	setMovieDetail(&details.SmallCoverImage, jsonData, "small_cover_image")
	setMovieDetail(&details.MediumCoverImage, jsonData, "medium_cover_image")
	setMovieDetail(&details.LargeCoverImage, jsonData, "large_cover_image")

	setMovieDetail(&details.State, jsonData, "state")

	torrentsData, exists := (*jsonData)["torrents"]

	if exists {
		torrentsList, ok := torrentsData.([]interface{})

		if ok {
			var torrents []*Movie.MovieTorrentInfo

			var appendListMutex sync.Mutex

			tmContext, _ := context.WithTimeout(client.Context, client.Timeout)

			var taskManager *TaskManager.TaskManager = TaskManager.CreateTaskManagerWithContext(tmContext, "YTS_TORRENT_PARSER_"+fmt.Sprintf("%d", (uintptr)(unsafe.Pointer(details))), Config.Main.TasksMaxThreads.YTS_TORRENT_PARSER)

			taskManager.Start()

			for _, t := range torrentsList {
				if torrentInfo, ok := t.(map[string]interface{}); ok {
					var torrent *Movie.MovieTorrentInfo = Movie.NewMovieTorrentInfo()

					setMovieDetail(&torrent.URL, &torrentInfo, "url")

					setMovieDetail(&torrent.Hash, &torrentInfo, "hash")

					setMovieDetail(&torrent.Quality, &torrentInfo, "quality")
					setMovieDetail(&torrent.Type, &torrentInfo, "type")

					setMovieDetail(&torrent.IsRepack, &torrentInfo, "is_repack")

					setMovieDetail(&torrent.VideoCodec, &torrentInfo, "video_codec")

					setMovieDetail(&torrent.BitDepth, &torrentInfo, "bit_depth")
					setMovieDetail(&torrent.AudioChannels, &torrentInfo, "audio_channels")

					setMovieDetail(&torrent.Seeds, &torrentInfo, "seeds")
					setMovieDetail(&torrent.Peers, &torrentInfo, "peers")

					setMovieDetail(&torrent.SizeString, &torrentInfo, "size")
					setMovieDetail(&torrent.Size, &torrentInfo, "size_bytes")

					setMovieDetail(&torrent.DateUploaded, &torrentInfo, "date_uploaded")
					setMovieDetail(&torrent.DateUploadedUnix, &torrentInfo, "date_uploaded_unix")

					taskManager.AddTaskWithDelay(func(t *TaskManager.Task) {
						err := Movie.ParseTorrentFromUrl(tmContext, torrent.URL, torrent)

						if err != nil {
							Logger.WARN("Failed to parse torrent file. [URL: " + torrent.URL + ", Message: " + err.Error() + "]")
							return
						}

						appendListMutex.Lock()
						torrents = append(torrents, torrent)
						appendListMutex.Unlock()
					}, Config.Main.TasksExecutionDelay.YTS_TORRENT_PARSER)
				}
			}

			taskManager.WaitForTasks()

			TaskManager.DeleteTaskManager(taskManager.Name)

			details.Torrents = torrents
		}
	}

	setMovieDetail(&details.DateUploaded, jsonData, "date_uploaded")
	setMovieDetail(&details.DateUploadedUnix, jsonData, "date_uploaded_unix")
}

func (this *Client) GetMovieList(params *MoviesListParameters) ([]*Movie.MovieDetails, error, float64) {
	url, err := url.Parse(this.ListMoviesEndpoint)

	var queryParams *MoviesListParameters = NewMoviesListParameters()

	if params != nil {
		queryParams = params
	}

	url.RawQuery = ConvertMoviesListParametersToURLParams(queryParams).Encode()

	if err != nil {
		return nil, err, 0
	}

	moviesJsonData, err := this.fetch(url, "GET", nil)

	if err != nil {
		return nil, err, 0
	}

	var movieCount float64 = 0

	movieCountData, exists := moviesJsonData["movie_count"]

	if exists {
		movieCount = movieCountData.(float64)
	}

	moviesListData, exists := moviesJsonData["movies"]

	if !exists {
		return nil, errors.New("`movies` field doesn't exist"), 0
	}

	moviesList, ok := moviesListData.([]interface{})

	if !ok {
		return nil, errors.New("`movies` field must be an array"), 0
	}

	var moviesListResult []*Movie.MovieDetails = make([]*Movie.MovieDetails, 0)

	var appendListMutex sync.Mutex

	tmContext, _ := context.WithTimeout(this.Context, this.Timeout)

	var movieParserTaskManager *TaskManager.TaskManager = TaskManager.CreateTaskManagerWithContext(tmContext, "YTS_MOVIE_PARSER_"+fmt.Sprintf("%d", (uintptr)(unsafe.Pointer(queryParams))), Config.Main.TasksMaxThreads.YTS_MOVIE_PARSER)

	movieParserTaskManager.Start()

	for _, value := range moviesList {
		item, ok := value.(map[string]interface{})

		if !ok {
			continue
		}

		var details *Movie.MovieDetails = Movie.NewMovieDetails()

		movieParserTaskManager.AddTaskWithDelay(func(t *TaskManager.Task) {
			parseMovieDetailsFromJsonData(this, details, &item)

			appendListMutex.Lock()
			moviesListResult = append(moviesListResult, details)
			appendListMutex.Unlock()
		}, Config.Main.TasksExecutionDelay.YTS_MOVIE_PARSER)
	}

	movieParserTaskManager.WaitForTasks()

	TaskManager.DeleteTaskManager(movieParserTaskManager.Name)

	return moviesListResult, nil, movieCount
}

func (this *Client) GetMovieCount(params *MoviesListParameters) (float64, error) {
	url, err := url.Parse(this.ListMoviesEndpoint)

	var queryParams *MoviesListParameters = NewMoviesListParameters()

	if params != nil {
		queryParams = params
	}

	url.RawQuery = ConvertMoviesListParametersToURLParams(queryParams).Encode()

	if err != nil {
		return 0, err
	}

	moviesJsonData, err := this.fetch(url, "GET", nil)

	if err != nil {
		return 0, err
	}

	var movieCount float64 = 0

	movieCountData, exists := moviesJsonData["movie_count"]

	if exists {
		movieCount = movieCountData.(float64)
	}

	return movieCount, nil
}

func (this *Client) GetMovieDetails(params *MovieDetailsParameters) (*Movie.MovieDetails, error) {
	url, err := url.Parse(this.MovieDetailsEndpoint)

	var queryParams *MovieDetailsParameters = NewMovieDetailsParameters(INVALID_MOVIE_DETAILS_PARAMETERS_ID)

	if params != nil {
		queryParams = params
	}

	url.RawQuery = ConvertMovieDetailsParametersToURLParams(queryParams).Encode()

	if err != nil {
		return nil, err
	}

	dataJson, err := this.fetch(url, "GET", nil)

	if err != nil {
		return nil, err
	}

	movieDetailsJson, exists := dataJson["movie"]

	if !exists {
		return nil, errors.New("`movie` field doesn't exists")
	}

	movieDetails, ok := movieDetailsJson.(map[string]interface{})

	if !ok {
		return nil, errors.New("`movie` must be an object")
	}

	var details *Movie.MovieDetails = Movie.NewMovieDetails()

	parseMovieDetailsFromJsonData(this, details, &movieDetails)

	return details, nil
}

func (this *Client) GetMovieSuggestions(params *MovieSuggestionsParameters) ([]*Movie.MovieDetails, error, float64) {
	url, err := url.Parse(this.MovieSuggestionsEndpoint)

	var queryParams *MovieSuggestionsParameters = NewMovieSuggestionsParameters(INVALID_MOVIE_DETAILS_PARAMETERS_ID)

	if params != nil {
		queryParams = params
	}

	url.RawQuery = ConvertMovieSuggestionsParametersToURLParams(queryParams).Encode()

	if err != nil {
		return nil, err, 0
	}

	moviesJsonData, err := this.fetch(url, "GET", nil)

	if err != nil {
		return nil, err, 0
	}

	var movieCount float64 = 0

	movieCountData, exists := moviesJsonData["movie_count"]

	if exists {
		movieCount = movieCountData.(float64)
	}

	moviesListData, exists := moviesJsonData["movies"]

	if !exists {
		return nil, errors.New("`movies` field doesn't exist"), 0
	}

	moviesList, ok := moviesListData.([]interface{})

	if !ok {
		return nil, errors.New("`movies` field must be an array"), 0
	}

	var moviesListResult []*Movie.MovieDetails = make([]*Movie.MovieDetails, 0)

	var appendListMutex sync.Mutex

	var movieParserTaskManager *TaskManager.TaskManager = TaskManager.CreateTaskManagerWithContext(this.Context, "YTS_MOVIE_PARSER_"+fmt.Sprintf("%d", (uintptr)(unsafe.Pointer(queryParams))), Config.Main.TasksMaxThreads.YTS_MOVIE_PARSER)

	movieParserTaskManager.Start()

	for _, value := range moviesList {
		item, ok := value.(map[string]interface{})

		if !ok {
			continue
		}

		var details *Movie.MovieDetails = Movie.NewMovieDetails()

		movieParserTaskManager.AddTaskWithDelay(func(t *TaskManager.Task) {
			parseMovieDetailsFromJsonData(this, details, &item)

			appendListMutex.Lock()
			moviesListResult = append(moviesListResult, details)
			appendListMutex.Unlock()
		}, Config.Main.TasksExecutionDelay.YTS_MOVIE_PARSER)
	}

	movieParserTaskManager.WaitForTasks()

	TaskManager.DeleteTaskManager(movieParserTaskManager.Name)

	return moviesListResult, nil, movieCount
}

func NewClient(ctx context.Context, timeout time.Duration) *Client {
	var client *Client = new(Client)

	var err error

	client.BaseURL = Defaults.YTS_API_BASE_URL
	client.ListMoviesEndpoint, err = url.JoinPath(client.BaseURL, Defaults.YTS_API_LIST_MOVIES_ENDPOINT)
	client.MovieDetailsEndpoint, err = url.JoinPath(client.BaseURL, Defaults.YTS_API_MOVIE_DETAILS_ENDPOINT)
	client.MovieSuggestionsEndpoint, err = url.JoinPath(client.BaseURL, Defaults.YTS_API_MOVE_SUGGESTIONS_ENDPOINT)

	client.Context = ctx
	client.Timeout = timeout

	client.HttpClient = new(http.Client)

	if err != nil {
		return nil
	}

	return client
}

func NewClientWithCustomURL(ctx context.Context, timeout time.Duration, baseURL string) *Client {
	var client *Client = NewClient(ctx, timeout)

	var err error

	client.BaseURL = baseURL
	client.ListMoviesEndpoint, err = url.JoinPath(client.BaseURL, Defaults.YTS_API_LIST_MOVIES_ENDPOINT)
	client.MovieDetailsEndpoint, err = url.JoinPath(client.BaseURL, Defaults.YTS_API_MOVIE_DETAILS_ENDPOINT)
	client.MovieSuggestionsEndpoint, err = url.JoinPath(client.BaseURL, Defaults.YTS_API_MOVE_SUGGESTIONS_ENDPOINT)

	if err != nil {
		return nil
	}

	return client
}
