package Crawler

import (
	"GServer/Config"
	"GServer/InternetArchive"
	"GServer/Logger"
	"GServer/Movie"
	"GServer/YTS"
	"context"
)

var YTSCrawler *Client = nil
var InternetArchiveCrawler *Client = nil

var MainCrawlerContext context.Context = nil
var MainCrawlerContextCancel context.CancelFunc = nil

func GetYTSSearchResult(client *Client) []*Movie.MovieDetails {
	ytsClient, ok := client.ServiceClient.(*YTS.Client)

	if !ok {
		Logger.ERROR("Failed to get YTS client service.")
		return []*Movie.MovieDetails{}
	}

	Logger.INFO("yts is : ", ytsClient)

	return []*Movie.MovieDetails{}
}

func GetYTSTotalMovies(client *Client) float64 {
	ytsClient, ok := client.ServiceClient.(*YTS.Client)

	if !ok {
		Logger.ERROR("Failed to get YTS client service.")
		return 0
	}

	var params *YTS.MoviesListParameters = YTS.NewMoviesListParameters()

	count, err := ytsClient.GetMovieCount(params)

	if err != nil {
		Logger.ERROR("Failed getting movie counts. [Message: " + err.Error() + "]")
		return 0
	}

	return count
}

func GetInternetArchiveSearchResult(client *Client) []*Movie.MovieDetails {
	iaClient, ok := client.ServiceClient.(*InternetArchive.Client)

	if !ok {
		Logger.ERROR("Failed to get Internet Archive client service.")
		return []*Movie.MovieDetails{}
	}

	Logger.INFO("ia client is : ", iaClient)

	return []*Movie.MovieDetails{}
}

func GetInternetArchiveTotalMovies(client *Client) float64 {
	iaClient, ok := client.ServiceClient.(*InternetArchive.Client)

	if !ok {
		Logger.ERROR("Failed to get Internet Archive client service.")
		return 0
	}

	Logger.INFO("ia client is : ", iaClient)

	return 0
}

func Initialize() {
	Logger.INFO("Initializing crawler...")

	MainCrawlerContext, MainCrawlerContextCancel = context.WithCancel(context.Background())

	YTSCrawler = NewClient(MainCrawlerContext, "YTS Crawler", int32(Config.Main.Crawler.YTSMovieCountPerSearch), 0)
	InternetArchiveCrawler = NewClient(MainCrawlerContext, "Internet Archive Crawler", int32(Config.Main.Crawler.InternetArchiveMovieCountPerSearch), 0)

	YTSCrawler.GetSearchResult = GetYTSSearchResult
	InternetArchiveCrawler.GetSearchResult = GetInternetArchiveSearchResult

	YTSCrawler.GetTotalMovieCount = GetYTSTotalMovies
	InternetArchiveCrawler.GetTotalMovieCount = GetInternetArchiveTotalMovies

	YTSCrawler.ServiceClient = YTS.NewClient(YTSCrawler.Context, 0)

	if Config.Main.CanUseYTSService {
		YTSCrawler.Start()
	}

	if Config.Main.CanUseInternetArchiveService {
		InternetArchiveCrawler.Start()
	}

	Logger.INFO("Crawler initialized.")
}

func Uninitialize() {
	Logger.INFO("Uninitializing crawler...")

	YTSCrawler.Stop()
	InternetArchiveCrawler.Stop()

	Logger.INFO("Crawler uninitialized.")
}
