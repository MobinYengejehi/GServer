package YTS

import (
	"fmt"
	"net/url"
)

type MovieSuggestionsParameters struct {
	MovieId int32
}

func NewMovieSuggestionsParameters(movieId int32) *MovieSuggestionsParameters {
	var params *MovieSuggestionsParameters = new(MovieSuggestionsParameters)

	params.MovieId = movieId

	return params
}

func ConvertMovieSuggestionsParametersToURLParams(params *MovieSuggestionsParameters) url.Values {
	if params == nil {
		return nil
	}

	var urlParams url.Values = url.Values{}

	urlParams.Add("movie_id", fmt.Sprintf("%d", params.MovieId))

	return urlParams
}
