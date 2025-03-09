package YTS

import (
	"fmt"
	"net/url"
)

const (
	INVALID_MOVIE_DETAILS_PARAMETERS_ID = 0
)

type MovieDetailsParameters struct {
	MovieId int32
	IMDBId  int32

	WithImages bool
	WithCast   bool
}

func NewMovieDetailsParameters(movieId int32) *MovieDetailsParameters {
	var params *MovieDetailsParameters = new(MovieDetailsParameters)

	params.MovieId = movieId
	params.IMDBId = INVALID_MOVIE_DETAILS_PARAMETERS_ID

	params.WithImages = false
	params.WithCast = false

	return params
}

func IsMovieDetailsParametersValid(params *MovieDetailsParameters) bool {
	if params == nil {
		return false
	}

	return params.MovieId != INVALID_MOVIE_DETAILS_PARAMETERS_ID || params.IMDBId != INVALID_MOVIE_DETAILS_PARAMETERS_ID
}

func ConvertMovieDetailsParametersToURLParams(params *MovieDetailsParameters) url.Values {
	if params == nil {
		return nil
	}

	var withImages string = "false"
	var withCast string = "false"

	if params.WithImages {
		withImages = "true"
	}

	if params.WithCast {
		withCast = "true"
	}

	var urlParams url.Values = url.Values{}

	urlParams.Add("with_images", withImages)
	urlParams.Add("with_cast", withCast)

	if params.MovieId != INVALID_MOVIE_DETAILS_PARAMETERS_ID {
		urlParams.Add("movie_id", fmt.Sprintf("%d", params.MovieId))
	} else if params.IMDBId != INVALID_MOVIE_DETAILS_PARAMETERS_ID {
		urlParams.Add("imdb_id", fmt.Sprintf("%d", params.IMDBId))
	}

	return urlParams
}
