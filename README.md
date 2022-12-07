# D-CarImplementation

...

## Design 



## Implementation and Tests


## Local Setup
To test D-Car locally, you can use the MongoDB setup provided in the `dev` folder.

To do so, execute the following commands:
```bash
cd dev
docker-compose up -d
```

This will start a MongoDB instance on port 27017 with the correct authentication setup.

After that, start the Go server with the following environment variables:

| Environment Variable        | Value         |
|-----------------------------|---------------|
| `MONGODB_DATABASE_HOST`     | localhost     |
| `MONGODB_DATABASE_NAME`     | ccsappvp2dcar |
| `MONGODB_DATABASE_USER`     | root          |
| `MONGODB_DATABASE_PASSWORD` | example       |
