## Build and run the application
````
docker build -t stress-test -f Dockerfile .
docker run --rm stress-test --url=https://google.com.br --requests=2 --concurrency=2 
````