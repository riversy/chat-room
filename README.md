# Chat Room

Simple Chat Room solution using React and Golang

## Testing

### Test Client

To test server enter to the server directory by the following command:

```bash 
cd client
```

Then run one of the following commands.

```bash
npm test 
```

### Test Server

To test server enter to the server directory by the following command:

```bash 
cd server
```

Then run one of the following commands:

```bash
godoc -http :8080
go test ./...
go test -coverprofile .coverate.out ./...
go tool cover -html=.coverate.out 
go test -bench ./...
```