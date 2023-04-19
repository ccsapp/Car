# CarImplementation
Car is a domain layer microservice of CCSAppVP2 that provides static and dynamic car data.

Since the models that Car uses are also needed in other microservices they were extracted to a separate repository. 
The models are available at 
[cargotypes](https://github.com/ccsapp/cargotypes) to provide mappings for the JSON responses.

The provided API endpoints of Car are specified in the
[API specification](https://github.com/ccsapp/CarDesign/blob/main/openapi.yaml).

## Local Setup Mode
To run the microservice Car locally, you can use the MongoDB setup provided in the `dev` directory.

To do so, execute the following commands:
```bash
cd dev
docker compose up -d
```

This will start a MongoDB instance on port 27021 (**non-default port** to avoid collisions with other databases) with
the correct authentication setup.

After that, start the Go server with the following environment variable set:

| Environment Variable | Value | Comment                       | 
|----------------------|-------|-------------------------------|
| `CAR_LOCAL_SETUP`    | true  | Enables the local setup mode. |

You might want to set `CAR_LOCAL_SETUP` in your IDE's default run configuration.
For example, in IntelliJ IDEA, you can do this [as described here](https://stackoverflow.com/a/32761503).

In the local setup mode, the microservice will use the configuration specified in `environment/localSetup.env`.
It contains the correct database connection information matching the docker compose file such that no further
configuration is required. This information will be embedded into the binary at build time.

However, you can still override the configuration by setting environment variables
described in the "Deployment or Custom Setup" section manually.

The default configuration values of local setup mode can also be found in the table of the "Deployment or Custom Setup"
section.

If the local setup mode is enabled, the integration tests (NOT the application itself) will try to detect if the
correct docker compose stack is running and will print a warning if it is not.

> After you have started the microservice in local setup mode, you can access it at
> [http://localhost:8001](http://localhost:8001).

## Deployment or Custom Setup
Do not use the local setup mode in a deployment or a custom setup, i.e. do not set the `CAR_LOCAL_SETUP` environment
variable. Instead, use the following environment variables to configure the microservice:

| Environment Variable        | Local Setup Value                                   | Comment                                                                                                               |
|-----------------------------|-----------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------|
| `MONGODB_CONNECTION_STRING` | mongodb://root:example@localhost:27021/ccsappvp2car |                                                                                                                       |
| `MONGODB_DATABASE_NAME`     | ccsappvp2car                                        |                                                                                                                       |
| `CAR_EXPOSE_PORT`           | 8001                                                | Optional, defaults to 80. This is the port this microservice is exposing. The local setup exposes a non-default port! |
| `CAR_COLLECTION_PREFIX`     | localSetup-                                         | Optional. A (unique) prefix that is prepended to every database collection of this service.                           |

## Testing

### Test Setup
The Unit Tests of Car depend on automatically generated Go mocks.
You need to install [mockgen](https://github.com/golang/mock#installation) to generate them.
After the installation, execute `go generate ./...` in the `src` directory of this project.

### Running the Tests
To run the tests locally, choose the local setup mode, or use a custom setup as described above
to configure database access for the integration tests.

> **Please note:** The integration tests will ignore the `CAR_COLLECTION_PREFIX` environment variable and use
dynamically generated collection names to avoid collisions with other tests.

After that, you can run the tests using `go test ./...` in the `src` directory.
