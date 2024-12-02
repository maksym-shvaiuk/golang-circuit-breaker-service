## Worker

- Reads data about its requests from the database.
- Stores this data to the local Redis cache.
- Does all processing in goroutines.
- If it needs to stop processing, it writes to the "stopped processing queue" and returns from the goroutine.

### API

- Uses Kafka for incoming requests.
- Provides an HTTP API to stop the worker:
  - Stops all goroutines.
  - Saves state to the database.
  - Publishes to the "stopped processing" queue.

## Controller

- Manages workers.
- Starts new workers.
- Stops workers.
- Maps request IDs to workers.

### API

- Queue of "stopped processing".
- Kafka topic "allocate worker".

## API Gateway (API GW)

- Gets request ID by `UserID` and endpoint.
- Tries to get Worker ID from Redis (read-only access).
- Publishes to the "allocate worker" Kafka topic and subscribes to the response.
- Publishes to the Worker ID partition with data to process.