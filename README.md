# How to run program locally
We must provide environment variables
````
SYMBOL='MSFT' NDAYS=8 ALPHA_VANTAGE_API_KEY='API_KEY' go run cmd/*.go 
````

## Features
Currently we have only 2 endpoints:
1. ``GET /healthz`` - will check health of server
2. ``GET /stock`` - will return the metadata about stock, returns the last number of days (NDAYS) and
average closing price over those days for stock.

If you don't provide any data in request it will automatically get SYMBOL and NDAYS from environment variable. 

But if we have received a data in request the values of environment variables will be ignored, instead we will use
data received from the request. 

For example,
````
{
	"symbol": "AAPL",
	"ndays": 10
}
````

The API will respond details about Apple stock price for the last 10 days with it's average closing price


# Docker

In this section we will explain how to create a docker image, running web service and
publish docker image to Docker Hub

You can choose one of the following ways:

**Creating image from Docker file**
````
docker build -t 2112fir/forgerock -f build/Dockerfile .
````

**Publishing to Docker Hub**
Login to Docker Hub via CLI

1) Direct usage of password in CLI is not recommended 
````
docker login -p {Passowrd} -u {Username}
````

2) Creating access token from https://hub.docker.com/settings/security, preffered way
````
docker login -u 2112fir
````

At the password prompt, enter the personal access token.

You can push a new image to this repository using the CLI
````
docker push 2112fir/forgerock
````

**Running server from the above created docker image**
````
docker run -e SYMBOL='MSFT' -e NDAYS=2 -e ALPHA_VANTAGE_API_KEY='API KEY' -p 8080:8080 2112fir/forgerock
````

Later on we will use publicly pushed image inside Kubernetes manifest.

OR

We can use Docker compose file to build image locally and test it from local Docker.

First replace the environment variable content from file ``build/docker-compose.yml``
````
environment:
  - SYMBOL="MSFT"
  - NDAYS=2
  - ALPHA_VANTAGE_API_KEY="MUST PROVIDE REAL API KEY"
````

Then:
````
cd build

docker-compose up
````


# What can be improved ?
- Currently, the unit and integration tests are missing from the project, due to the time constraint for the purpose of this project, the tests are not done, in real projects, the use of BDD test with the help of Cucumber framework would be useful to check the behaviour of the business logic, as it allows to make a client rest api calls. 
If we implement BDD test we should create mock code for external providers such as AlphaVantage from Interface.

- The project could benefit from the automatic dependency injection tools, such as "uber/fx" as currently the dependencies
provided manually.