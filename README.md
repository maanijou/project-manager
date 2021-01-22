# Project Management

Managing your projects efficiently.

## Devops

`docker` and `docker-compose` are used to handle docker images like a database. I used pgadmin for some database management tasks as well.

You can use `make` to run the images.

You have to define a `database.env` file in your root folder so that everything works fine. env files are usually some samples (`database.env.sample`) and should not be inside the repository. Here I just put that for your convenience!

## Handling external API

For external API's there are multiple scenarios and approaches. One can synchronize data between two services. In a microservice architecture, it's always good to have separate data storage and handle them separately. So rather than getting all the employees' data and saving them locally, I'm going to handle table joins and data dependencies at the app level. In this sense, this app works as an API Gateway for different microservices.

## Running Instructions

You can use make file commands to run the project.

* make test: Run go tests
* make build: Building docker images
* make up: Running docker images using docker-compose
* make logs: Show docker-compose logs
* make rm: Remove docker images and volumes
* make start: Start docker images without building
* make stop: Stop docker images
* make rest: Remove docker images and volumes, build and running again and showing the logs
* make bash_db: Start bash on database image
* make bash_go: Start sh on app image

## Available routes and methods
By default, docker-compose exposes the API on port 3000.

GET:
* `/api/v1/health`: To get check API health and availability.

### Employee

GET:
* `/api/v1/employees`
* `/api/v1/employees/{uuid}`

### Projects
GET:
* `/api/v1/projects/?page=1&limit=10`: showing projects without any participants

* `/api/v1/projects/1`: Get project with id=1

POST:
* `/api/v1/projects`: Add new project using a json body of a project (id field will be ignored and created automatically)
* `/api/v1/projects/1/participants`: Add participants to project=1 using a json body (`[ {"id": "uuid"}, { "id": "uuid"}]`)


PUT:
* `/api/v1/projects`: update a project using a json body of the project (id field must be correct)

DELETE:
* `/api/v1/projects/1`: Remove project with id=1
* `/api/v1/projects/1/participants/{uuid}`: remove participant={uuid} from project=1
## Notes
I decided not to show participants on projects route
## TODO

- [ ] For some improvements, a caching mechanism for fetching data from external API can be applied.
- [ ] Using swagger for both getting data from external API and designing. A nice automatic documentation would be generated as well.
- [ ] Showing useful messages when departments are different between participants and the owner (For now, we won't add them without showing any particular error).
- [ ] Writing test codes for project API's
- [ ] CI/CD
