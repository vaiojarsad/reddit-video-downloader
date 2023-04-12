// Package reddit define primitives to interact with Reddit
package reddit

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/vaiojarsad/reddit-video-downloader/internal/utils/av"
	"github.com/vaiojarsad/reddit-video-downloader/internal/utils/general"
)

func DownloadVideo(videoURL, outputFileName string, noAudio, noVideo bool) error {
	var err error
	var videoFileName string
	if !noVideo {
		videoFileName, err = getFile(videoURL)
		if err != nil {
			return fmt.Errorf("error getting video: %w", err)
		}
	}

	var audioFileName string
	if !noAudio {
		re := regexp.MustCompile(`(?s)\_(.*)\.`)
		audioURL := re.ReplaceAllString(videoURL, "_audio.")
		audioFileName, err = getFile(audioURL)
		if err != nil {
			return fmt.Errorf("error getting audio: %w", err)
		}
	}

	if !noAudio && !noVideo {
		err = mergeAudioAndVideo(audioFileName, videoFileName, outputFileName)
		if err != nil {
			return fmt.Errorf("error merging audio & video: %w", err)
		}
	} else if noVideo {
		err = os.Rename(audioFileName, outputFileName)
		if err != nil {
			return fmt.Errorf("error renaming audio file: %w", err)
		}
	} else {
		err = os.Rename(videoFileName, outputFileName)
		if err != nil {
			return fmt.Errorf("renaming video file: %w", err)
		}
	}
	return nil
}

func getFile(url string) (fileName string, err error) {
	file, err := os.CreateTemp("", "reddit-vd-*.mp4")
	if err != nil {
		return "", err
	}
	defer general.CloseWithCheck(file)
	err = downloadFile(url, file)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func downloadFile(url string, o io.Writer) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:107.0) Gecko/20100101 Firefox/107.0")

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return err
	}
	defer general.CloseWithCheck(resp.Body)

	_, err = io.Copy(o, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func mergeAudioAndVideo(audioFile, videoFile, outputFile string) error {
	m := av.NewMerger()
	err := m.Merge(audioFile, videoFile, outputFile)
	if err != nil {
		return err
	}
	return nil
}
