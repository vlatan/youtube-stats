# YouTube Metadata

Supply a YouTube video ID to get some of the essential info about the video.


## Config

Create an `.env` file with `YOUTUBE_API_KEY` variable or export it to your environment.


## Installation

Build CLI and WEB binaries and start a webserver.
```
make run
```

Access the web app on `localhost:8080`.  

![YouTube Metadata Web App](./screenshot.png "YouTube Metadata Web App")


Use the CLI.
```
./bin/cli jNQXAC9IVRw
```

(Re)start the webserver.
```
./bin/app
```

Clean up.
```
make clean
```

## Run the Web App via Docker

Build the image.
```
docker build -t yt-stats .
```

Run.
```
docker run -p 8080:8080 --env-file=.env yt-stats
```

## Deployed

Hosted on **Railway**. Expect cold start.    
https://ytmeta.up.railway.app


## License

[![License: MIT](https://img.shields.io/github/license/vlatan/youtube-stats?label=License)](/LICENSE "License: MIT")