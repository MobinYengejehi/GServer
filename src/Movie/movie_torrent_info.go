package Movie

import (
	"GServer/Config"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"path"
	"time"

	"github.com/anacrolix/torrent/metainfo"
)

type MovieTorrentFileInfo struct {
	Name      string
	Extension string

	Path string

	SizeString string
	Size       float64
}

type MovieTorrentInfo struct {
	URL    string
	Magent string

	Name string

	Hash string

	Quality string
	Type    string

	IsRepack string

	VideoCodec string

	BitDepth      string
	AudioChannels string

	Seeds float64
	Peers float64

	SizeString string
	Size       float64

	CreatedBy string

	Files    []*MovieTorrentFileInfo
	MainFile *MovieTorrentFileInfo

	DateUploaded     string
	DateUploadedUnix float64
}

func NewMovieTorrentFileInfo() *MovieTorrentFileInfo {
	var fileInfo *MovieTorrentFileInfo = new(MovieTorrentFileInfo)

	fileInfo.Name = ""
	fileInfo.Path = ""

	fileInfo.SizeString = ""
	fileInfo.Size = 0

	return fileInfo
}

func NewMovieTorrentInfo() *MovieTorrentInfo {
	var torrentInfo *MovieTorrentInfo = new(MovieTorrentInfo)

	torrentInfo.URL = ""
	torrentInfo.Magent = ""

	torrentInfo.Name = ""

	torrentInfo.Hash = ""

	torrentInfo.Quality = ""
	torrentInfo.Type = ""

	torrentInfo.IsRepack = ""

	torrentInfo.VideoCodec = ""

	torrentInfo.BitDepth = ""
	torrentInfo.AudioChannels = ""

	torrentInfo.Seeds = 0
	torrentInfo.Peers = 0

	torrentInfo.SizeString = ""
	torrentInfo.Size = 0

	torrentInfo.CreatedBy = ""

	torrentInfo.Files = []*MovieTorrentFileInfo{}
	torrentInfo.MainFile = nil

	torrentInfo.DateUploaded = ""
	torrentInfo.DateUploadedUnix = 0

	return torrentInfo
}

func IsMovieTorrentInfoValid(torrentInfo *MovieTorrentInfo) bool {
	if torrentInfo == nil {
		return false
	}

	return len(torrentInfo.URL) > 0
}

func UnpackTorrentInfoFromURL(torrentInfo *MovieTorrentInfo, url string) error {
	if torrentInfo == nil {
		return errors.New("Invalid MovieTorrentInfo")
	}

	if len(url) < 1 {
		return errors.New("Invalid URL")
	}

	return nil
}

func sizeToString(targetSize float64) string {
	var units []string = []string{"Byte", "KB", "MB", "GB"}
	var sizes []float64 = []float64{1}

	for i := range units {
		if i == 0 {
			continue
		}

		sizes = append(sizes, math.Pow(1024, float64(i)))
	}

	var sizeCount int = len(sizes)

	for i := 0; i < sizeCount; i++ {
		var size float64 = sizes[i]
		var nextSize float64 = size

		if i < sizeCount-1 {
			nextSize = sizes[i+1]
		}

		if size == nextSize {
			goto Found
		}

		if targetSize > size && targetSize < nextSize {
			goto Found
		}

		continue

	Found:
		{
			return fmt.Sprintf("%.2f %s", float32(targetSize/size), units[i])
		}
	}

	return fmt.Sprintf("%0.2f Byte", float32(targetSize))
}

func ParseTorrentFromUrl(ctx context.Context, url string, torrentInfo *MovieTorrentInfo) error {
	if len(url) < 1 {
		return errors.New("Invalid URL")
	}

	if torrentInfo == nil {
		return errors.New("Invalid MovieTorrentInfo pointer")
	}

	request, err := http.NewRequestWithContext(ctx, "GET", url, bytes.NewBufferString(""))

	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("Bad status: " + response.Status)
	}

	meta, err := metainfo.Load(io.Reader(response.Body))

	if err != nil {
		return err
	}

	info, err := meta.UnmarshalInfo()

	if err != nil {
		return err
	}

	magnet, err := meta.MagnetV2()

	if err != nil {
		return err
	}

	torrentInfo.URL = url
	torrentInfo.Magent = magnet.String()

	torrentInfo.Name = info.Name

	torrentInfo.Hash = meta.HashInfoBytes().HexString()

	{
		torrentInfo.Size = float64(info.TotalLength())
		torrentInfo.SizeString = sizeToString(torrentInfo.Size)
	}

	torrentInfo.CreatedBy = meta.CreatedBy

	for _, file := range info.Files {
		var fileInfo *MovieTorrentFileInfo = NewMovieTorrentFileInfo()

		fileInfo.Path = file.DisplayPath(&info)

		_, fileName := path.Split(fileInfo.Path)

		fileInfo.Extension = path.Ext(fileInfo.Path)
		fileInfo.Name = fileName[0 : len(fileName)-len(fileInfo.Extension)]

		if !Config.IsTorrentFileExtensionValid(fileInfo.Extension) {
			continue
		}

		fileInfo.Size = float64(file.Length)
		fileInfo.SizeString = sizeToString(fileInfo.Size)

		if torrentInfo.MainFile == nil {
			if Config.IsMainTorrentFileExtension(fileInfo.Extension) {
				torrentInfo.MainFile = fileInfo
			}
		}

		torrentInfo.Files = append(torrentInfo.Files, fileInfo)
	}

	torrentInfo.DateUploaded = time.Unix(meta.CreationDate, 0).Format(time.DateTime)
	torrentInfo.DateUploadedUnix = float64(meta.CreationDate)

	if torrentInfo.MainFile == nil {
		return errors.New("Couldn't find any main file in torrent file.")
	}

	return nil
}
