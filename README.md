# video-platform

## Dependencies
- `go`
- `npm`
- `docker`
- `docker-compose`

## Setup frontend
- `npm i`
- `npm start`

## Setup `docker`
- `sudo systemctl enable docker --force`
- `sudo usermod -aG docker ${USER}`

## Setup project
```sh
docker network create web
docker-compose up -d
```

## Calculate total lines of code
```sh
git ls-files | grep -Ev ".json|go.sum|nginx|.md" | xargs wc -l
```
