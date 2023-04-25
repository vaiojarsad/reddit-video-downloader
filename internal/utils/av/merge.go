// Package av defines audio & video handling utils.
package av

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/vaiojarsad/reddit-video-downloader/internal/utils/general"
)

// NewMerger returns Merger implementation
func NewMerger() Merger {
	return &FFMPEGMerger{}
}

// Merger defines an abstraction to merge an audio file and a video file.
type Merger interface {
	Merge(audioFilename, videoFileName, outputFileName string) error
}

// FFMPEGMerger implements av.Merger backed up by the ffmpeg util.
type FFMPEGMerger struct {
}

// Merge implementation using the ffmpeg util.
func (m *FFMPEGMerger) Merge(audioFilename, videoFileName, outputFileName string) error {
	f, err := os.CreateTemp("", "ffmpeg-merger-*.mp4")
	if err != nil {
		return fmt.Errorf("error creating temporary output file: %w", err)
	}
	mergedFileName := f.Name()
	general.CloseWithCheck(f)
	args := []string{"-i", videoFileName, "-i", audioFilename, "-c:v", "copy", "-c:a", "aac", "-y", mergedFileName}
	cmd := exec.Command("ffmpeg", args...)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error running ffmpeg: %w", err)
	}
	return general.Move(mergedFileName, outputFileName)
}
