# project-manager

Managing your projects easily


## Devops

`docker` and `docker-compose` are used to handle docker images like a database. I used pgadmin for some database management tasks as well.

You can use `make` to run the images.

You have to define a `database.env` file in your root folder so that everything works fine. env files are usually some samples (`database.env.sample`) and should not be inside the repository. Here I just put that for your convenience!

## Handling external API

For external API's there are multiple scenarios and approaches. One can synchronize data between two services. In a microservice architecture, it's always good to have separate data storage and handle them separately. So rather than getting all the employees' data and saving them locally, I'm going to handle table joins and data dependencies at the app level. In this sence this app works as an API Gateway for different microservices.


# TODO

[_] For some improvements, a caching mechanism for fetching data from external API can be applied.
[_] Using swagger for both getting data from external API and designing. A nice automatic documentation would be generated as well.
[_] Showing useful messages when departments are different between participants and the owner (For now, we won't add them without showing any particular error).
[_] 