# CleanCast Setup for UNRAID

This guide will help you set up CleanCast on UNRAID using the Docker Compose plugin.

## Prerequisites

1. **UNRAID Docker Compose Plugin** - Install from Community Applications
2. **YouTube API v3 Key** - Get one [here](https://developers.google.com/youtube/v3/getting-started)

## Quick Setup

### 1. Create Directory Structure

SSH into your UNRAID server and create the config directory:

```bash
mkdir -p /mnt/user/appdata/cleancast
```

### 2. Upload docker-compose.yml

Copy the `docker-compose.yml` file to your UNRAID server:
- Via UNRAID web UI file manager
- Or using SCP/SFTP to `/mnt/user/appdata/cleancast/docker-compose.yml`

### 3. Configure Environment Variables

Edit the `docker-compose.yml` file and set:

**Required:**
- `GOOGLE_API_KEY` - Your YouTube API key

**Recommended:**
- `SPONSORBLOCK_CATEGORIES` - Which ad types to remove (already set to common ones)
- `MIN_DURATION` - Filter out short videos (default: 3m)

**Optional:**
- `TOKEN` - Secure your endpoints (recommended if exposing to internet)
- `TRUSTED_HOSTS` - Limit which domains can access the service

### 4. Start the Container

Using Docker Compose plugin:
1. Go to Docker tab in UNRAID
2. Click "Compose" (if you have the plugin installed)
3. Navigate to your docker-compose.yml location
4. Click "Compose Up"

Or via command line:
```bash
cd /mnt/user/appdata/cleancast
docker-compose up -d
```

### 5. Verify It's Running

Check the container is running:
```bash
docker ps | grep cleancast
```

Check logs:
```bash
docker logs cleancast
```

## Usage

### Create RSS Feeds

#### For a YouTube Playlist:
1. Find the playlist ID from the URL (e.g., `PLxxxxxx`)
2. Your RSS feed URL: `http://YOUR_UNRAID_IP:8080/rss/PLxxxxxx`

#### For a YouTube Channel:
1. Get the channel ID using [this tool](https://www.tunepocket.com/youtube-channel-id-finder/)
2. Your RSS feed URL: `http://YOUR_UNRAID_IP:8080/channel/UCxxxxxx`

#### With Date Filter (Recommended for Channels):
To avoid hitting API quota limits:
```
http://YOUR_UNRAID_IP:8080/channel/UCxxxxxx?date=01-01-2024
```

#### If Using TOKEN Authentication:
Add `?token=YOUR_TOKEN` to all URLs:
```
http://YOUR_UNRAID_IP:8080/rss/PLxxxxxx?token=YOUR_TOKEN
```

### Add to Podcast App

1. Copy your RSS feed URL
2. In your podcast app (Apple Podcasts, Pocket Casts, etc.):
   - Look for "Add by URL" or "Add custom RSS feed"
   - Paste your CleanCast RSS feed URL
3. The app will automatically download episodes with sponsors removed!

## Configuration Tips

### SponsorBlock Categories

You can remove different types of content by setting `SPONSORBLOCK_CATEGORIES`:

- `sponsor` - Paid promotions
- `selfpromo` - Self-promotion (merch, social media)
- `interaction` - "Like and subscribe" reminders
- `intro` - Intermission/intro animation
- `outro` - Endcards/credits
- `preview` - Recap/preview of other episodes
- `music_offtopic` - Non-music videos with music
- `filler` - Filler tangent/jokes

Example to remove everything:
```yaml
- SPONSORBLOCK_CATEGORIES=sponsor,selfpromo,interaction,intro,outro,preview,music_offtopic,filler
```

### Age-Restricted Videos

If you need to download age-restricted or members-only content:

1. Export cookies from your browser using a cookies.txt extension
2. Save as `cookies.txt` in `/mnt/user/appdata/cleancast/cookies.txt`
3. Add to docker-compose.yml:
   ```yaml
   - COOKIES_FILE=cookies.txt
   ```
4. Restart container

## Troubleshooting

### Container Won't Start
- Check logs: `docker logs cleancast`
- Verify GOOGLE_API_KEY is set correctly
- Make sure port 8080 isn't already in use

### Downloads Failing
- Check you have enough disk space in `/mnt/user/appdata/cleancast/audio`
- Verify ffmpeg is working: `docker exec cleancast ffmpeg -version`
- Check logs for detailed error messages

### API Quota Exceeded
- Use date filters on channel endpoints: `?date=MM-DD-YYYY`
- Consider getting multiple API keys and alternating
- Limit the number of videos per feed using `?limit=20`

### No Audio File Downloads
- Episode files are downloaded on-demand when your podcast app requests them
- First play may take longer as the video is downloaded and processed
- Subsequent plays will be instant as the file is cached

## Updating

To update to the latest version:

```bash
cd /mnt/user/appdata/cleancast
docker-compose pull
docker-compose up -d
```

## Support

- GitHub Issues: https://github.com/ikoyhn/clean-cast/issues
- Original README: See README.md in this directory
