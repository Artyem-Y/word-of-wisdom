# Word of Wisdom

The Word of Wisdom Server project is a simple server for handling requests with request rate limiting, DDoS protection mechanisms, and Proof of Work validation.

## Project Structure

main logic in word-of-wisdom-server:

- `server/server.go`
- `server/ddos.go`
- `server/client_tracker.go`
- `server/connection_limiter.go`
- ...

client's logic is in word-of-wisdom-client

## Description

The Word of Wisdom Server project includes the following components:

### server.go

The main file for starting the server and handling requests.

### ddos.go

Implementation of a Distributed Denial of Service (DDoS) protection mechanism through IP address blocking.

### client_tracker.go

Component for tracking and limiting client's requests to prevent excessive usage.

### connection_limiter.go

Functionality to limit the number of concurrent connections to the server.

## Proof of Work

The project uses Proof of Work (PoW) validation to ensure that clients' requests are legitimate. 
The PoW challenge-response mechanism helps mitigate abuse and ensures that clients must provide computational 
proof of their request's validity.

The PoW mechanism:

1. Generating the Challenge: the server generates a challenge value by calling the generateChallenge function. 
This challenge is a random integer between 0 and a maximum value (maxValue) defined in the code.

2. Sending the Challenge: the challenge value is used to construct a challenge message in the format: 
"Solve Proof of Work: Find a nonce x such that SHA256(x + challenge) has a specified number of leading zeros."
The server sends the challenge message to the client over the network connection.

3. Client Response: the client receives the challenge message and calculates a Proof of Work response, or nonce, that satisfies the challenge criteria. 
The response is a value that, when combined with the challenge, produces a SHA-256 hash with a specified number of leading zeros.

4. Checking Proof of Work: the client sends back its Proof of Work response, which is the nonce concatenated with the hash calculated using the challenge and nonce.
The server uses the checkProofOfWork function to verify the client's Proof of Work. 
The function splits the client's response into the nonce and hash parts.
It checks if hash part can be decoded and if the hash's leading characters match the required number of zeros 
defined by proofDifficulty.

5. Connection Approval: if the Proof of Work is valid, the server approves the connection and sends a confirmation message along with a randomly selected quote 
from the quoteList to the client.

6. Closing the Connection: after sending the confirmation and quote, the server flushes the connection to ensure the data is sent before closing.
The server then closes the write and read sides of the connection to complete the communication.

### Running App

To run the app, in word-of-wisdom/word-of-wisdom-server use the following command:

```sh
make

```

## Testing

Unit tests are provided to ensure the functionality of the server and its components. The server's tests are located in the `test/word-of-wisdom/server` directory.

### Running Tests

To run the app, use the following command:

```sh
go test ./...
```

### Check connections:

```sh
nc 127.0.0.1 8080
```

or in word-of-wisdom-client

create in word-of-wisdom-client/conf config.yml file and set `numConnections` value and run in word-of-wisdom-client:

```sh
go run main.go
```