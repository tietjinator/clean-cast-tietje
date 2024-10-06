package services

import (
	podcastRss "github.com/eduncan911/podcast"
	"github.com/labstack/echo/v4"
	log "github.com/labstack/gommon/log"
	"ikoyhn/podcast-sponsorblock/internal/models"
	"net/http"
	"time"
)

const rssUrl = "/rss/youtubePlaylistId="

func GenerateRssFeed(podcast models.Podcast, c echo.Context) error {
	log.Info("[RSS FEED] Generating RSS Feed with Youtube and Apple metadata")
	now := time.Now()
	ytPodcast := podcastRss.New(podcast.PodcastName, handler(c.Request())+rssUrl+podcast.YoutubePodcastId, podcast.Description, &now, &now)
	ytPodcast.AddImage(podcast.ImageUrl)

	if podcast.PodcastEpisodes != nil {
		for _, podcastEpisode := range podcast.PodcastEpisodes {
			enclosure := podcastRss.Enclosure{
				URL:    handler(c.Request()) + "/media/" + podcastEpisode.YoutubeVideoId + ".m4a",
				Length: 0,
				Type:   podcastRss.M4A,
			}
			podcastItem := podcastRss.Item{
				Title:       podcastEpisode.EpisodeName,
				Description: podcastEpisode.EpisodeDescription,
				Enclosure:   &enclosure,
				PubDate:     &now,
			}
			ytPodcast.AddItem(podcastItem)
		}

		err := ytPodcast.Encode(c.Response().Writer)
		if err != nil {
			return err
		}
	}
	return c.NoContent(http.StatusOK)
}

func handler(r *http.Request) string {
	var scheme string
	if r.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"
	}
	host := r.Host
	url := scheme + "://" + host
	return url
}
