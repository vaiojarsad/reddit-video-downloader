package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
)

const (
	flagVideoURL  = "u"
	flagOutputFile = "o"
	flagNoVideo    = "no-video"
	flagNoAudio    = "no-audio"
)

func main() {
	var videoURL, outputFileName string
	var noVideo, noAudio bool
	flag.StringVar(&videoURL, flagVideoURL, "", "Video's url.")
	// flag.StringVar(&outputFileName, flagOutputFile, "", "Path to the output file.")
	flag.BoolVar(&noVideo, flagNoVideo, false, "Don't download video, just audio")
	flag.BoolVar(&noAudio, flagNoAudio, false, "Don't download audio, just video")
	flag.Parse()
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })

	if !seen[flagVideoURL] {
		flag.Usage()
		os.Exit(2)
	}

	// Nothing to do.
	if seen[flagNoVideo] && seen[flagNoAudio] && noAudio && noVideo {
		flag.Usage()
		os.Exit(2)
	}

	var err error
	var videoFileName string
	if !noVideo {
		videoFileName, err = getFile(videoURL)
		if err != nil {
			log.Fatal("Error getting video", "err", err)
		}
	}

	var audioFileName string
	if !noAudio {
		re := regexp.MustCompile(`(?s)\_(.*)\.`)
		audioURL := re.ReplaceAllString(videoURL, "_audio.")
		audioFileName, err = getFile(audioURL)
		if err != nil {
			log.Fatal("Error getting audio", "err", err)
		}
	}

	outputFileName = "/data/video.mp4"
	if !noAudio && !noVideo {
		err = mergeAudioAndVideo(audioFileName, videoFileName, outputFileName)
		if err != nil {
			log.Fatal("Error merging audio & video", "err", err)
		}
	} else if noVideo {
		err = os.Rename(audioFileName, outputFileName)
		if err != nil {
			log.Fatal("Error renaming audio file", "err", err)
		}
	} else {
		err = os.Rename(videoFileName, outputFileName)
		if err != nil {
			log.Fatal("Error renaming video file", "err", err)
		}
	}
}

func closeWithCheck(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println("Error in close", "err", err)
	}
}

func getFile(url string) (fileName string, err error){
	file, err := ioutil.TempFile("", "reddit-vd-*.mp4")
	if err != nil {
		return "", err
	}
	defer closeWithCheck(file)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:107.0) Gecko/20100101 Firefox/107.0")

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return "", err
	}
	defer closeWithCheck(resp.Body)

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func mergeAudioAndVideo(audioFile, videoFile, outputFile string) error {
	args := []string{"-i", videoFile, "-i", audioFile, "-c:v", "copy", "-c:a", "aac", outputFile}
	cmd := exec.Command("ffmpeg", args...)
	return cmd.Run()
}
