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

| Syntax    | Description                                                                           |
| --------- | ------------------------------------------------------------------------------------- |
| httpflood | Initiates a http flood to the specified url/ip/hostname. `httpflood -h` for more info |
| history   | Shows a list of the latest 50 requests                                                |

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

# Attacks

- [HTTP flood](#http-flood)
- [Slowloris](#slowloris)
- [UDP flood](#udp-flood)

### HTTP flood

This attack is a simple http flood attack. It sends a request to the specified url/ip/hostname. The attack will continue until the user stops it. The attack is very effective. It also works for websites that are protected by cloudflare.

I did notice when adding support for the user-agent and some random headers the attack is much more potent. The error called `Forbidden 403` did not happen as much as it did when I was not using user-agent and random headers.

- Test the amount of goroutines you can use. For example, on my desktop I can usually do 400-1500 until I blue screen. On my laptop I can do 20000 and it all works fine.
- Sometimes more goroutines is not always better. I noticed that when I used 1000 goroutines the attack was not as potent as when I used 500 goroutines. This was when testing on a smaller server with only 1 core and 1GB of ram. So play around with the amount of goroutines you use.
- Play around with the timeout flag to see what effect it has on the attack. Keeping the connection open for a longer time might be more effective.

### Slowloris

This is still under development. But can share that under testing I was unable to take down a website.

### UDP flood

This is still under development. Promising results so far.
