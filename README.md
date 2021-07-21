# University's Site

This project is a sample university's site that shows a list of approved candidates.

## Requirements

You will need Docker installed in your computer to run the university's site.

## Build

To build the site, execute the following command:
```sh
docker build -t test:latest .
``` 

Or the following comand (if you have `make` installed):
```sh
make build
``` 

## Run

To run the site, execute the following command:
```sh
docker run --publish 80:8080 test
```

Or the following comand (if you have `make` installed):
```sh
make run
``` 

Then, you can check the website accessing the address [http://localhost](http://localhost)