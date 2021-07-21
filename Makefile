build:
	docker build -t test:latest .

run:
	docker run --publish 80:8080 test