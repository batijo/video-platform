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
