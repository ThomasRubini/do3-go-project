# Getting started

- Copy `.env.example` to `.env`
- Update the API key value
- run `go run cmd/nutritionapp/main.go`

# Using docker

Build: `docker build . -t do3-go-project:latest`
Run: `docker run -e "FDC_API_KEY=YOUR_KEY_HERE" --rm -it do3-go-project:latest`

# Authors
Thomas RUBINI

Hugo DU PELOUX
