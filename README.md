# CarImplementation

Since the models that Car uses are also needed in other microservices they were extracted to a separate repository. 
The models are available at the private Git repository 
[CarGoTypes](https://git.scc.kit.edu/cm-tm/cm-team/projectwork/pse/domain/d-cargotypes) to provide mappings for the 
JSON responses.
Further information on the usage of private Git repositories with go can be found there.

## Local Setup
To run the microservice Car locally, you can use the MongoDB setup provided in the `dev` directory.

To do so, execute the following commands:
```bash
cd dev
docker-compose up -d
```

This will start a MongoDB instance on port 27017 with the correct authentication setup.

After that, start the Go server with the following environment variables:

| Environment Variable        | Value         | Comment  |
|-----------------------------|---------------|----------|
| `MONGODB_DATABASE_HOST`     | localhost     |          |
| `MONGODB_DATABASE_NAME`     | ccsappvp2dcar |          |
| `MONGODB_DATABASE_USER`     | root          |          |
| `MONGODB_DATABASE_PASSWORD` | example       |          |
| `CAR_COLLECTION_PREFIX`     | someprefix    | optional |

## Testing

### Test Setup
The Unit Tests of FleetManagement depend on automatically generated Go mocks.
You need to install [mockgen](https://github.com/golang/mock#installation) to generate them.
After the installation, execute `go generate ./...` in the `src` directory of this project.
The provided API endpoints of FleetManagement are specified in the [API specification](https://git.scc.kit.edu/cm-tm/cm-team/projectwork/pse/application/p-fleetmanagementdesign).

### Running the Tests
To run the tests locally, set up a database as documented above.
But instead of setting environment variables, put the values in testdata/testdb.env.
Beware, that the prefix is not considered but generated automatically in the test setup.

Now, you can run the tests using `go test ./...` in the `src` directory.

 