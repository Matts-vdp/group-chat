# group-chat
A group chat based on a server with multiple clients using a TCP connection. Made in Golang
Any messages are broadcasted to all connected clients.

## Usage
1. build the program or download the release
2. start server with
```
chat -s -net <ip:port>
```
the default net id is localhost:5000

3. start the client with
```
chat -net <ip:port>
```
4. enter a user name (optional)
5. chat
