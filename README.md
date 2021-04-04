# video-platform

## Dependencies
- `go`
- `python`
- `docker`
- `docker-compose`

## Setup `docker`
- `sudo systemctl enable docker --force`
- `sudo usermod -aG docker ${USER}`

## Setup project
```sh
docker network create web
docker-compose up -d
```

To run the Ingest client alongside the server for a local Demo
```sh
docker-compose -f docker-compose.yml -f docker-compose.client.yml up -d
```
