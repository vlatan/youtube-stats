# YouTube Metadata

Supply a YouTube video URL to get some of the essential info about the video.


## Config

Create an `.env` file with `YOUTUBE_API_KEY` variable or export it to your environment.


## Installation

Build both CLI and WEB binaries and start the app.
```
make run
```

Or just use `air` if you want live reloading for the app.
```
air
```

Access the web app on `localhost:8080`.  

![YouTube Metadata Web App](./screenshot.png "YouTube Metadata Web App")


Use the CLI.
```
./bin/cli <url>
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

## License

[![License: MIT](https://img.shields.io/github/license/vlatan/youtube-stats?label=License)](/LICENSE "License: MIT")