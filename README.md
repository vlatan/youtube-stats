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
./cli jNQXAC9IVRw
```

(Re)start the webserver.
```
./ui
```

Clean up.
```
make clean
```

## License

[![License: MIT](https://img.shields.io/github/license/vlatan/youtube-stats?label=License)](/LICENSE "License: MIT")