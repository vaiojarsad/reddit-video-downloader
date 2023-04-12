// Package cmd is used to define Cobra stuff
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vaiojarsad/reddit-video-downloader/reddit"
	"log"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "reddit-video-downloader",
		Short: "Reddit Video Downloader",
		Long:  "Reddit Video Downloader",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			noVideo, err := cmd.Flags().GetBool("no-video")
			if err != nil {
				cobra.CheckErr(err)
			}
			noAudio, err := cmd.Flags().GetBool("no-video")
			if err != nil {
				cobra.CheckErr(err)
			}
			if noVideo && noAudio {
				cobra.CheckErr(fmt.Errorf("no-video and no-audio cannot be set to true at the same time"))
			}
			return nil
		},
		RunE: rootRunE,
	}
	videoURL, outputFileName string
	noVideo, noAudio         bool
)

func init() {
	rootCmd.Flags().BoolVar(&noVideo, "no-video", false, "Don't download video, just audio")
	rootCmd.Flags().BoolVar(&noAudio, "no-audio", false, "Don't download audio, just video")
	rootCmd.Flags().StringVarP(&outputFileName, "output", "o", "", "Output file")
	rootCmd.Flags().StringVarP(&videoURL, "video-url", "u", "", "Video's URL.")

	err := rootCmd.MarkFlagRequired("output")
	if err != nil {
		cobra.CheckErr(err)
	}
	err = rootCmd.MarkFlagRequired("video-url")
	if err != nil {
		cobra.CheckErr(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("error... %v", err)
		os.Exit(1)
	}
}

func rootRunE(_ *cobra.Command, _ []string) error {
	return reddit.DownloadVideo(videoURL, outputFileName, noAudio, noVideo)
}
