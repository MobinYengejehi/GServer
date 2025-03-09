package InternetArchive

import (
	"fmt"
	"net/url"
)

const (
	SEARCH_PARAMETERS_FIELD_AVG_RATING          = "avg_rating"
	SEARCH_PARAMETERS_FIELD_BACKUP_LOCATION     = "backup_location"
	SEARCH_PARAMETERS_FIELD_BTIH                = "btih"
	SEARCH_PARAMETERS_FIELD_CALL_NUMBER         = "call_number"
	SEARCH_PARAMETERS_FIELD_COLLECTION          = "collection"
	SEARCH_PARAMETERS_FIELD_CONTRIBUTOR         = "contributor"
	SEARCH_PARAMETERS_FIELD_COVERAGE            = "coverage"
	SEARCH_PARAMETERS_FIELD_CREATOR             = "creator"
	SEARCH_PARAMETERS_FIELD_DATE                = "date"
	SEARCH_PARAMETERS_FIELD_DESCRIPTION         = "description"
	SEARCH_PARAMETERS_FIELD_DOWNLOADS           = "downloads"
	SEARCH_PARAMETERS_FIELD_EXTERNAL_IDENTIFIER = "external-identifier"
	SEARCH_PARAMETERS_FIELD_FOLD_OUT_COUNT      = "foldoutcount"
	SEARCH_PARAMETERS_FIELD_FORMAT              = "format"
	SERACH_PARAMETERS_FIELD_GENRE               = "genre"
	SEARCH_PARAMETERS_FIELD_IDENTIFIER          = "identifier"
	SEARCH_PARAMETERS_FIELD_IMAGE_COUNT         = "imagecount"
	SEARCH_PARAMETERS_FIELD_INDEX_FLAG          = "indexflag"
	SEARCH_PARAMETERS_FIELD_ITEM_SIZE           = "item_size"
	SEARCH_PARAMETERS_FIELD_LANGUAGE            = "language"
	SEARCH_PARAMETERS_FIELD_LICENSE_URL         = "licenseurl"
	SEARCH_PARAMETERS_FIELD_MEDIA_TYPE          = "mediatype"
	SEARCH_PARAMETERS_FIELD_MEMBERS             = "members"
	SEARCH_PARAMETERS_FIELD_MONTH               = "month"
	SEARCH_PARAMETERS_FIELD_NAME                = "name"
	SEARCH_PARAMETERS_FIELD_NO_INDEX            = "noindex"
	SEARCH_PARAMETERS_FIELD_NUM_REVIEWS         = "num_reviews"
	SEARCH_PARAMETERS_FIELD_OAI_UPDATE_DATE     = "oai_updatedate"
	SEARCH_PARAMETERS_FIELD_PUBLIC_DATE         = "publicdate"
	SEARCH_PARAMETERS_FIELD_PUBLISHER           = "publisher"
	SEARCH_PARAMETERS_FIELD_RELATED_EXTERNAL_ID = "related-external-id"
	SEARCH_PARAMETERS_FIELD_REVIEW_DATE         = "reviewdate"
	SEARCH_PARAMETERS_FIELD_RIGHTS              = "rights"
	SEARCH_PARAMETERS_FIELD_SCANNING_CENTRE     = "scanningcentre"
	SEARCH_PARAMETERS_FIELD_SOURCE              = "source"
	SEARCH_PARAMETERS_FIELD_STRIPPED_TAGS       = "stripped_tags"
	SEARCH_PARAMETERS_FIELD_SUBJECT             = "subject"
	SEARCH_PARAMETERS_FIELD_TITLE               = "title"
	SEARCH_PARAMETERS_FIELD_TYPE                = "type"
	SEARCH_PARAMETERS_FIELD_VOLUME              = "volume"
	SEARCH_PARAMETERS_FIELD_WEEK                = "week"
	SEARCH_PARAMETERS_FIELD_YEAR                = "year"

	SEARCH_PARAMETERS_OUTPUT_TYPE_JSON        = "json"
	SEARCH_PARAMETERS_OUTPUT_TYPE_XML         = "xml"
	SEARCH_PARAMETERS_OUTPUT_TYPE_HTML_TABLES = "tables"
	SEARCH_PARAMETERS_OUTPUT_TYPE_CSV         = "csv"
	SEARCH_PARAMETERS_OUTPUT_TYPE_RSS         = "rss"
)

type SearchParameters struct {
	Query string

	Feilds []string

	Rows int32
	Page int32

	OutputType string
}

func NewSearchParameters(query string) *SearchParameters {
	var params *SearchParameters = new(SearchParameters)

	params.Query = query

	params.Feilds = []string{
		SEARCH_PARAMETERS_FIELD_IDENTIFIER,
		SEARCH_PARAMETERS_FIELD_TITLE,
		SEARCH_PARAMETERS_FIELD_DESCRIPTION,
		SEARCH_PARAMETERS_FIELD_ITEM_SIZE,
		SEARCH_PARAMETERS_FIELD_DATE,
		SEARCH_PARAMETERS_FIELD_LANGUAGE,
		SEARCH_PARAMETERS_FIELD_VOLUME,
	}

	params.Rows = 20
	params.Page = 1

	params.OutputType = SEARCH_PARAMETERS_OUTPUT_TYPE_JSON

	return params
}

func ConvertSearchParametersToURLParams(params *SearchParameters) url.Values {
	if params == nil {
		return nil
	}

	var urlParams url.Values = url.Values{}

	urlParams.Add("q", params.Query)

	for _, value := range params.Feilds {
		urlParams.Add("fl[]", value)
	}

	urlParams.Add("rows", fmt.Sprintf("%d", params.Rows))
	urlParams.Add("page", fmt.Sprintf("%d", params.Page))

	urlParams.Add("output", params.OutputType)

	return urlParams
}
