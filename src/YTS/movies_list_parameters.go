package YTS

import (
	"fmt"
	"net/url"
)

const (
	MOVIES_PARAMS_QUALITY_480P      = "480p"
	MOVIES_PARAMS_QUALITY_720P      = "720p"
	MOVIES_PARAMS_QUALITY_1080P     = "1080p"
	MOVIES_PARAMS_QUALITY_1080PX265 = "1080p.x265"
	MOVIES_PARAMS_QUALITY_2160P     = "2160p"
	MOVIES_PARAMS_QUALITY_3D        = "3D"
	MOVIES_PARAMS_QUALITY_ALL       = ""

	MOVIES_PARAMS_GENRE_ALL = ""

	MOVIES_PARAMS_SORT_BY_TITLE          = "title"
	MOVIES_PARAMS_SORT_BY_YEAR           = "year"
	MOVIES_PARAMS_SORT_BY_RATING         = "rating"
	MOVIES_PARAMS_SORT_BY_PEERS          = "peers"
	MOVIES_PARAMS_SORT_BY_SEEDS          = "seeds"
	MOVIES_PARAMS_SORT_BY_DOWNLOAD_COUNT = "download_count"
	MOVIES_PARAMS_SORT_BY_LIKE_COUNT     = "like_count"
	MOVIES_PARAMS_SORT_BY_DATE_ADDED     = "date_added"

	MOVIES_PARAMS_ORDER_BY_DESC = "desc"
	MOVIES_PARAMS_ORDER_BY_ASC  = "asc"
)

type MoviesListParameters struct {
	Limit int32
	Page  int32

	Quality       string
	MinimumRating float32

	QueryTerm string

	Genre string

	SortBy  string
	OrderBy string

	WithRtRatings bool
}

func NewMoviesListParameters() *MoviesListParameters {
	var params *MoviesListParameters = new(MoviesListParameters)

	params.Limit = 20
	params.Page = 1

	params.Quality = MOVIES_PARAMS_QUALITY_ALL

	params.MinimumRating = 0.0

	params.QueryTerm = "0"

	params.Genre = MOVIES_PARAMS_GENRE_ALL

	params.SortBy = MOVIES_PARAMS_SORT_BY_DATE_ADDED
	params.OrderBy = MOVIES_PARAMS_ORDER_BY_DESC

	params.WithRtRatings = false

	return params
}

func ConvertMoviesListParametersToURLParams(params *MoviesListParameters) url.Values {
	if params == nil {
		return nil
	}

	var withRtRatings string = "false"

	if params.WithRtRatings {
		withRtRatings = "true"
	}

	var urlParams url.Values = url.Values{}

	urlParams.Add("limit", fmt.Sprintf("%d", params.Limit))
	urlParams.Add("page", fmt.Sprintf("%d", params.Page))
	urlParams.Add("minimum_rating", fmt.Sprintf("%.1f", params.MinimumRating))
	urlParams.Add("query_term", params.QueryTerm)
	urlParams.Add("sort_by", params.SortBy)
	urlParams.Add("with_rt_ratings", withRtRatings)

	if len(params.Quality) > 0 {
		urlParams.Add("quality", params.Quality)
	}

	if len(params.Genre) > 0 {
		urlParams.Add("genre", params.Genre)
	}

	return urlParams
}
