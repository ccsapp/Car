# D-CarImplementation

...

## Design 



## Implementation and Tests


## Local Setup
To run D-Car locally, you can use the MongoDB setup provided in the `dev` directory.

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


## Running the Tests
To run the tests locally, set up a database as documented above.
But instead of setting environment variables, put the values in testdata/testdb.env.
Beware, that the prefix is not considered but generated automatically in the test setup.

Now, you can run the tests using `go test ./...` in the `src` directory.
