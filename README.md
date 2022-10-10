# HTTP GO REQUEST

This purpose of this application is to deploy a GO HTTP server and be able to retrieve data from several URLs

## Description

This provide the code for a GO HTTP server that gets data from different URLs and returns a merged list of the data.  It allows to sort by "views" or "relevanceScore" as well as limit the amount of records retrieved.

The default URLs that are queried are the following:
* https://raw.githubusercontent.com/assignment132/assignment/main/duckduckgo.json
* https://raw.githubusercontent.com/assignment132/assignment/main/google.json
* https://raw.githubusercontent.com/assignment132/assignment/main/wikipedia.json

If deploying with Kubernetes, the limits and URLs to query can be modified.

The webserver is listening on port 8080.

## Getting Started

### Dependencies

* Server to deploy the Go HTTP Server

If using Kubernetes to deploy:
* Kubernetes environment set up and running
* At least two worker nodes for replicas


### Installing

Server:
* Go to /gohttp
* Execute:
```
go mod tidy
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gohttp \*.go
```
* Run the gohttp executable

Kubernetes:
* Create the docker image using the Dockerfile provided in the repo
* Update the gohttp-deploy.yaml with the docker image registry path as well as the URL for the ingress
* Deploy the gohttp-deploy.yaml file and if desired the gohttp-cm.yaml file

### Web server parameters

* limit: This parameter will set up the maximum amount of records that will be retrieved from the merged data from the URLs.  By default minimum is 2 and maximum limit is 200.
* sortKey: This parameter will sort the data in ascending numbers based on "views" and "relevanceScore".

### Kubernetes configuration

* gohttp-deploy.yaml file: Used to deploy a deployment with 2 replias, a service and an ingress.
* gohttp-cm.yaml file: Allows to dynamically change the maximum limit and URLs to process.
* Dockerfile file: Used to create the docker image for the deployment

Note:
When adding new URLs, the schema for the data is as follows:
{"data": [
  {"url": "URL",
  "views": "Number of views",
  "relevanceScore": "Relevance Score"}
]}

### Usage

* Getting all data records (up to the maximum limit):
```
[URL]/gohttp
```

* Sorting by "views" or "relevanceScore":
```
[URL]/gohttp?sortKey=views
[URL]/gohttp?sortKey=relevanceScore
```

* Limiting the amount of records:
```
[URL]/gohttp?limit=10
```

The return data adds a count field to the JSON.

## Testing

There is a gohttp/gohttp_test.go file that can be used for unit testing.

Run the following command in the gohttp folder which will test sorting and limits.
```
go test
```

## Help

Any advise for common problems or issues.

## Authors

Contributors names and contact info

Andres Calvo

## Version History

* 0.1
    * Initial Release

## License

This project is licensed under the [NAME HERE] License - see the LICENSE.md file for details
