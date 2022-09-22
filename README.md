# GoDown

### Why?

Well, we use this tool to test our own applications.

### How to use?

Run the program from source code (assuming you are in the root of the project):

```go
go run cmd/godown <command>
```

Run the program from binary:

```bash
<path to binary> <command>
```

Run the program when installed with go:

```bash
godown <command>
```

### Commands

**httpflood** - Initiates a http flood to the specified url/ip/hostname. `get -h` for more info.

**history** - Shows a list of the latest 50 requests.

### TODO

Will work on this list when I got some free time. If you want to contribute, feel free to do so.

- [x] HTTP flood attack.
- [ ] Slowloris attack.
- [ ] SYN flood attack.
- [ ] UDP flood attack.
- [ ] Add support for connecting multiple applications (ddos) and attack the same host.
- [ ] Move necessary code to /pkg.
- [ ] Add binary for windows.
- [ ] Add binary for linux.
- [ ] Add binary for mac.
