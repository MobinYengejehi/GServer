package Movie

const (
	INVALID_MOVIE_DETAIL_ID = 0
)

type MovieDetails struct {
	Id                float64
	SpecialIdentifier string

	URL string

	IMDBCode string

	Title        string
	TitleEnglish string
	TitleLong    string
	Slug         string

	Year    float64
	Rating  float64
	Runtime float64

	Genres []string

	LikeCount float64

	Summery          string
	DescriptionIntro string
	DescriptionFull  string
	Synopsis         string

	YTTrailerCode string

	Language string

	MPARating string

	BackgroundImage         string
	BackgroundImageOriginal string
	SmallCoverImage         string
	MediumCoverImage        string
	LargeCoverImage         string

	State string

	Size float64

	Torrents []*MovieTorrentInfo

	DateUploaded     string
	DateUploadedUnix float64
}

func NewMovieDetails() *MovieDetails {
	var details *MovieDetails = new(MovieDetails)

	details.Id = INVALID_MOVIE_DETAIL_ID
	details.SpecialIdentifier = ""

	details.URL = ""

	details.IMDBCode = ""

	details.Title = ""
	details.TitleEnglish = ""
	details.TitleLong = ""
	details.Slug = ""

	details.Year = 0
	details.Rating = 0.0
	details.Runtime = 0

	details.Genres = nil

	details.LikeCount = 0

	details.Summery = ""
	details.DescriptionIntro = ""
	details.DescriptionFull = ""
	details.Synopsis = ""

	details.YTTrailerCode = ""

	details.Language = ""

	details.MPARating = ""

	details.BackgroundImage = ""
	details.BackgroundImageOriginal = ""
	details.SmallCoverImage = ""
	details.MediumCoverImage = ""
	details.LargeCoverImage = ""

	details.State = ""

	details.Size = 0

	details.Torrents = nil

	details.DateUploaded = ""
	details.DateUploadedUnix = 0

	return details
}

func IsMovieDetialsValid(details *MovieDetails) bool {
	if details == nil {
		return false
	}

	return details.Id != 0 || len(details.SpecialIdentifier) > 0
}
