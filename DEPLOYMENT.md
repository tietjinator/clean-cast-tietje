# CleanCast Deployment Guide

## Quick Fix for Docker Run Error

If you see: `invalid reference format: repository name must be lowercase`

This means you're using a git branch name instead of a Docker image name.

## Two Deployment Options

### Option 1: Use Official Upstream Image (RECOMMENDED)

The easiest and fastest way. Since we merged all upstream changes, the official image has everything you need:

```bash
docker run \
  -d \
  --name='CleanCast' \
  --net='bridge' \
  -e 'GOOGLE_API_KEY'='YOUR_API_KEY' \
  -e 'TOKEN'='YOUR_TOKEN' \
  -e 'SPONSORBLOCK_CATEGORIES'='sponsor,selfpromo,intro,outro' \
  -e 'MIN_DURATION'='3m' \
  -p '8416:8080/tcp' \
  -v '/mnt/zcache/appdata/CleanCast':'/config':'rw' \
  ikoyhn/clean-cast:latest
```

✅ **Pros:**
- Instant deployment, no build time
- Already tested and stable
- Easy to update

### Option 2: Build Your Own Custom Image

If you've made custom code changes or want full control:

#### Step 1: Clone your repository to UNRAID

```bash
cd /tmp
git clone https://github.com/tietjinator/clean-cast-tietje.git
cd clean-cast-tietje
git checkout claude/youtube-episode-downloader-011CUSDrYCHSeobZ7Wa6yBnc
```

#### Step 2: Build the Docker image

```bash
docker build -t cleancast:custom .
```

This will take a few minutes as it compiles the Go application.

#### Step 3: Run with your custom image

```bash
docker run \
  -d \
  --name='CleanCast' \
  --net='bridge' \
  -e 'GOOGLE_API_KEY'='YOUR_API_KEY' \
  -e 'TOKEN'='YOUR_TOKEN' \
  -e 'SPONSORBLOCK_CATEGORIES'='sponsor,selfpromo,intro,outro' \
  -e 'MIN_DURATION'='3m' \
  -p '8416:8080/tcp' \
  -v '/mnt/zcache/appdata/CleanCast':'/config':'rw' \
  cleancast:custom
```

#### Automated Script

Or use the provided script:

```bash
# Edit the script to update your API keys first!
./build-and-deploy.sh
```

## Environment Variables Reference

| Variable | Description | Your Setting |
|----------|-------------|--------------|
| `GOOGLE_API_KEY` | YouTube API v3 key (REQUIRED) | Set this to your actual key |
| `TOKEN` | Security token for RSS feeds | Already set |
| `SPONSORBLOCK_CATEGORIES` | Ad types to remove | `sponsor,selfpromo,intro,outro` |
| `MIN_DURATION` | Minimum video length | `3m` (3 minutes) |

## Verify It's Working

```bash
# Check container status
docker ps | grep CleanCast

# View logs
docker logs CleanCast

# Test the endpoint
curl http://localhost:8416/
```

## Accessing Your Podcasts

Your CleanCast server is now running on `http://YOUR-UNRAID-IP:8416`

### Create an RSS Feed:

**For a YouTube Playlist:**
```
http://YOUR-UNRAID-IP:8416/rss/PLAYLIST_ID?token=A0VgrfV8VkfnEARmChkcwAUDPEF1TOB
```

**For a YouTube Channel:**
```
http://YOUR-UNRAID-IP:8416/channel/CHANNEL_ID?token=A0VgrfV8VkfnEARmChkcwAUDPEF1TOB&date=01-01-2024
```

Add these URLs to your podcast app!

## Troubleshooting

### Container won't start
```bash
# Check detailed logs
docker logs CleanCast

# Common issues:
# - Invalid API key
# - Port 8416 already in use
# - Volume paths don't exist
```

### API Quota Issues
- Use the `date` parameter on channel feeds to limit results
- Example: `?date=06-01-2024&token=YOUR_TOKEN`

### Audio not downloading
- Check disk space: `df -h /mnt/user/Media/clean-cast-audio-downloads`
- Verify ffmpeg: `docker exec CleanCast ffmpeg -version`

## Updating

### If using upstream image:
```bash
docker stop CleanCast
docker rm CleanCast
docker pull ikoyhn/clean-cast:latest
# Run the docker run command again
```

### If using custom image:
```bash
cd /tmp/clean-cast-tietje
git pull
docker build -t cleancast:custom .
docker stop CleanCast
docker rm CleanCast
# Run the docker run command again
```
