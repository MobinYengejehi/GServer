package Config

import (
	"GServer/Defaults"
	"GServer/Logger"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

const (
	CONFIG_FILE_NAME = "crawler_config.json"

	DEFAULT_CONFIG_JSON_DATA = `{
	"http_host_address" : "%s",
	"can_use_yts_service" : true,
	"can_use_ia_service" : true,
	"crawler" : {
		"yts_movie_count_per_search" : %d,
		"ia_movie_count_per_search" : %d
	},
	"tasks_max_threads" : {
		"HTTP_SERVER" : %d,
		"CRAWLER_MAIN" : %d,
		"MOVIE_CRAWLER_YTS" : %d,
		"MOVIE_CRAWLER_IA" : %d,
		"YTS_MOVIE_PARSER" : %d,
		"IA_MOVIE_PARSER" : %d,
		"YTS_TORRENT_PARSER" : %d,
		"IA_TORRENT_PARSER" : %d
	},
	"tasks_execution_delay" : {
		"HTTP_SERVER" : %d,
		"CRAWLER_MAIN" : %d,
		"MOVIE_CRAWLER_YTS" : %d,
		"MOVIE_CRAWLER_IA" : %d,
		"YTS_MOVIE_PARSER" : %d,
		"IA_MOVIE_PARSER" : %d,
		"YTS_TORRENT_PARSER" : %d,
		"IA_TORRENT_PARSER" : %d
	},
	"valid_torrent_file_extensions" : [
		".mp4",
		".webm",
		".webp",
		".mkv",
		".mov",
		".mpeg",
		".asf",
		".avi",
		".avchd",
		".raw",
		".m4v",
		".wmv",
		".flv",
		".3gp",
		".srt"
	],
	"main_torrent_file_extensions" : [
		".mp4",
		".webm",
		".webp",
		".mkv",
		".mov",
		".mpeg",
		".asf",
		".avi",
		".avchd",
		".raw",
		".m4v",
		".wmv",
		".flv",
		".3gp",
		".mkv"
	]
}`
)

type ConfigCrawler struct {
	YTSMovieCountPerSearch             int `json:"yts_movie_count_per_search"`
	InternetArchiveMovieCountPerSearch int `json:"ia_movie_count_per_search"`
}

type ConfigTasksMaxThreads struct {
	HTTP_SERVER int `json:"HTTP_SERVER"`

	CRAWLER_MAIN int `json:"CRAWLER_MAIN"`

	MOVIE_CRAWLER_YTS int `json:"MOVIE_CRAWLER_YTS"`
	MOVIE_CRAWLER_IA  int `json:"MOVIE_CRAWLER_IA"`

	YTS_MOVIE_PARSER int `json:"YTS_MOVIE_PARSER"`
	IA_MOVIE_PARSER  int `json:"IA_MOVIE_PARSER"`

	YTS_TORRENT_PARSER int `json:"YTS_TORRENT_PARSER"`
	IA_TORRENT_PARSER  int `json:"IA_TORRENT_PARSER"`
}

type ConfigTasksExecutionDelay struct {
	HTTP_SERVER time.Duration `json:"HTTP_SERVER"`

	CRAWLER_MAIN time.Duration `json:"CRAWLER_MAIN"`

	MOVIE_CRAWLER_YTS time.Duration `json:"MOVIE_CRAWLER_YTS"`
	MOVIE_CRAWLER_IA  time.Duration `json:"MOVIE_CRAWLER_IA"`

	YTS_MOVIE_PARSER time.Duration `json:"YTS_MOVIE_PARSER"`
	IA_MOVIE_PARSER  time.Duration `json:"IA_MOVIE_PARSER"`

	YTS_TORRENT_PARSER time.Duration `json:"YTS_TORRENT_PARSER"`
	IA_TORRENT_PARSER  time.Duration `json:"IA_TORRENT_PARSER"`
}

type Config struct {
	HttpHostAddress string `json:"http_host_address"`

	CanUseYTSService             bool `json:"can_use_yts_service"`
	CanUseInternetArchiveService bool `json:"can_use_ia_service"`

	Crawler ConfigCrawler `json:"crawler"`

	TasksMaxThreads     ConfigTasksMaxThreads     `json:"tasks_max_threads"`
	TasksExecutionDelay ConfigTasksExecutionDelay `json:"tasks_execution_delay"`

	ValidTorrentFileExtensions []string `json:"valid_torrent_file_extensions"`
	MainTorrentFileExtensions  []string `json:"main_torrent_file_extensions"`
}

var Main Config = Config{}

func GetDefaultCondigJsonString() string {
	return fmt.Sprintf(
		DEFAULT_CONFIG_JSON_DATA,
		Defaults.DEFAULT_HTTP_SERVER_HOST_ADDRESS,
		Defaults.CRAWLER_YTS_MOVIE_COUNT_PER_SEARCH,
		Defaults.CRAWLER_INTERNET_ARCHIVE_MOVIE_COUNT_PER_SAERCH,
		Defaults.TASKS_MAX_THREADS_HTTP_SERVER,
		Defaults.TASKS_MAX_THREADS_CRAWLER_MAIN,
		Defaults.TASKS_MAX_THREADS_MOVIE_CRAWLER_YTS,
		Defaults.TASKS_MAX_THREADS_MOVIE_CRAWLER_IA,
		Defaults.TASKS_MAX_THREADS_YTS_MOVIE_PARSER,
		Defaults.TASKS_MAX_THREADS_IA_MOVIE_PARSER,
		Defaults.TASKS_MAX_THREADS_YTS_TORRENT_PARSER,
		Defaults.TASKS_MAX_THREADS_IA_TORRENT_PARSER,
		Defaults.TASKS_EXECUTION_DELAY_HTTP_SERVER,
		Defaults.TASKS_EXECUTION_DELAY_CRAWLER_MAIN,
		Defaults.TASKS_EXECUTION_DELAY_MOVIE_CRAWLER_YTS,
		Defaults.TASKS_EXECUTION_DELAY_MOVIE_CRAWLER_IA,
		Defaults.TASKS_EXECUTION_DELAY_YTS_MOVIE_PARSER,
		Defaults.TASKS_EXECUTION_DELAY_IA_MOVIE_PARSER,
		Defaults.TASKS_EXECUTION_DELAY_YTS_TORRENT_PARSER,
		Defaults.TASKS_EXECUTION_DELAY_MOVIE_CRAWLER_IA)
}

func WriteConfig() {
	exePath, err := os.Executable()
	exePath = strings.ReplaceAll(exePath, "\\", "/")

	if err != nil {
		Logger.ERROR("Couldn't get application path to write config file. [Message: " + err.Error() + "]")
		return
	}

	applicationDir := path.Dir(exePath)

	var configFilePath string = path.Join(applicationDir, CONFIG_FILE_NAME)

	if _, err := os.Stat(configFilePath); !os.IsNotExist(err) {
		err := os.Remove(configFilePath)

		if err != nil {
			Logger.ERROR("Couldn't remove config file. [Message: " + err.Error() + "]")
			return
		}
	}

	err = os.WriteFile(configFilePath, []byte(GetDefaultCondigJsonString()), 0644)

	if err != nil {
		Logger.ERROR("Couldn't write config file. [Message: " + err.Error() + "]")
	}
}

func ReadConfig() {
	exePath, err := os.Executable()
	exePath = strings.ReplaceAll(exePath, "\\", "/")

	if err != nil {
		Logger.ERROR("Couldn't get application path to read config file. [Message: " + err.Error() + "]")
		return
	}

	applicationDir := path.Dir(exePath)

	var configFilePath string = path.Join(applicationDir, CONFIG_FILE_NAME)

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		err := json.Unmarshal([]byte(GetDefaultCondigJsonString()), &Main)

		if err != nil {
			Logger.ERROR("Couldn't parse default config json data. [Message: " + err.Error() + "]")
			return
		}

		WriteConfig()

		return
	}

	data, err := os.ReadFile(configFilePath)

	if err != nil {
		Logger.ERROR("Couldn't read config file. [Message: " + err.Error() + "]")
		return
	}

	err = json.Unmarshal(data, &Main)

	if err != nil {
		Logger.ERROR("Coudln't parse config file json data. [Message: " + err.Error() + "]")
	}
}

func IsTorrentFileExtensionValid(extension string) bool {
	if len(Main.ValidTorrentFileExtensions) < 1 {
		return true
	}

	for _, ext := range Main.ValidTorrentFileExtensions {
		if ext == extension {
			return true
		}
	}

	return false
}

func IsMainTorrentFileExtension(extension string) bool {
	if len(Main.MainTorrentFileExtensions) < 1 {
		return true
	}

	for _, ext := range Main.MainTorrentFileExtensions {
		if ext == extension {
			return true
		}
	}

	return false
}

func Initialize() {
	Logger.INFO("Initializing config...")

	ReadConfig()

	Logger.INFO("Config initialized.")
}

func Uninitialize() {
	Logger.INFO("Uninitializing config...")

	Logger.INFO("Config uninitialized.")
}
