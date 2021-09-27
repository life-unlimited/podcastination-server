# podcastination-server

This is the linux server application for _podcastination_ - an app that imports, stores and provides podcasts. It
imports episodes from a pull directory and provides RSS feeds as well as an API.

## Why did we build this?

_Podcastination_ provided a simple way for us to import and store podcasts exactly as we need it. A relational database
keeps track of all episodes and allows creating a simple webapp which displays all available podcasts and episodes.

## Installation

You only need to PostgreSQL database being setup and ready to use. Then simply run the server application with the
following command:

```shell
podcastination-server --config <path-to-config>
```

If you want to use the webapp, have a look [here](https://github.com/life-unlimited/podcastination-webapp).

## Configuration

The configuration file is a JSON file which contains some fields that are needed.

```json
{
  "postgres_datasource": "host=127.0.0.1 port=5432 user=postgres password=pass dbname=postgres sslmode=disable",
  "pull_dir": "path/to/pull/dir",
  "podcast_dir": "path/to/podcast/dir",
  "import_interval": 15,
  "server_addr": "127.0.0.1:8000"
}
```

The `pull_dir` is the directory from where new episodes are being pulled. The `podcast_dir` is where all data is stored.
The `import_interval` is provided in minutes.

## Usage

In order to use _podcastination_ you currently need to set up podcasts yourself in the database. We did not have time
yet to implement podcast creation in the API or webapp and rarely create them so this is enough for us. If you want to
contribute, feel free to contact us.

In order to import an episode, you create a directory in the pull directory which contains a `task.json` that has the
following content:

```json
 {
  "podcast_key": "the-podcast-key",
  "season_key": "the-season-key",
  "title": "My episode title",
  "subtitle": "My optional episode subtitle",
  "date": "2021-01-03T09:30:00.00Z",
  "author": "Author name",
  "mp3_file": "file-name-of-the-recording.mp3",
  "yt_url": "optional-youtube-url",
  "pdf_file": "optional-file-name-of-a-pdf-file.pdf",
  "image_file": "optional-file-name-of-an-episode-image.png"
}
```

Provide the MP3 file as well as optional PDF and image files in the same directory. After successful import, _
podcastination_ will delete the folder.