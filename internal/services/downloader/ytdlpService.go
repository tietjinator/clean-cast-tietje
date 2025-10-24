package downloader

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/labstack/gommon/log"
	"github.com/lrstanley/go-ytdlp"
)

var youtubeVideoMutexes = &sync.Map{}

const youtubeVideoUrl = "https://www.youtube.com/watch?v="

func GetYoutubeVideo(youtubeVideoId string) (string, <-chan struct{}) {
	mutex, ok := youtubeVideoMutexes.Load(youtubeVideoId)
	if !ok {
		mutex = &sync.Mutex{}
		youtubeVideoMutexes.Store(youtubeVideoId, mutex)
	}

	mutex.(*sync.Mutex).Lock()

	// Check if the file is already being processed
	filePath := "/config/audio/" + youtubeVideoId + ".mp3"
	if _, err := os.Stat(filePath); err == nil {
		mutex.(*sync.Mutex).Unlock()
		return youtubeVideoId, make(chan struct{})
	}

	// If not, proceed with the download
	youtubeVideoId = strings.TrimSuffix(youtubeVideoId, ".mp3")
	ytdlp.Install(context.TODO(), nil)

	categories := os.Getenv("SPONSORBLOCK_CATEGORIES")
	if categories == "" {
		categories = "sponsor"
	}
	categories = strings.TrimSpace(categories)

	dl := ytdlp.New().
		NoProgress().
		ExtractAudio().
		AudioFormat("mp3").
		SponsorblockRemove(categories).
		NoPlaylist().
		FFmpegLocation("/usr/bin/ffmpeg").
		Continue().
		Paths("/config/audio").
		ProgressFunc(500*time.Millisecond, func(prog ytdlp.ProgressUpdate) {
			fmt.Printf(
				"%s @ %s [eta: %s] :: %s\n",
				prog.Status,
				prog.PercentString(),
				prog.ETA(),
				prog.Filename,
			)
		}).
		Output(youtubeVideoId + ".%(ext)s")

	cookiesFile := strings.TrimSpace(os.Getenv("COOKIES_FILE"))
	if cookiesFile != "" {
		dl.Cookies("/config/" + cookiesFile)
	}

	done := make(chan struct{})
	go func() {
		log.Infof("[DOWNLOADER] Starting download for video: %s", youtubeVideoId)
		r, err := dl.Run(context.TODO(), youtubeVideoUrl+youtubeVideoId)
		if err != nil {
			log.Errorf("[DOWNLOADER] Error downloading YouTube video %s: %v", youtubeVideoId, err)
		}
		if r.ExitCode != 0 {
			log.Errorf("[DOWNLOADER] YouTube video %s download failed with exit code %d", youtubeVideoId, r.ExitCode)
			log.Errorf("[DOWNLOADER] Stdout: %s", r.Stdout)
			log.Errorf("[DOWNLOADER] Stderr: %s", r.Stderr)
		} else {
			log.Infof("[DOWNLOADER] Successfully downloaded video: %s", youtubeVideoId)
		}
		mutex.(*sync.Mutex).Unlock()

		close(done)
	}()

	return youtubeVideoId, done
}
