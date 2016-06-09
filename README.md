# slash

`slash` is a Go library for creating [Slack Slash Commands](https://api.slack.com/slash-commands).

Disclaimer: I created this because I wanted to create a slash command, and I've been wanting to learn Go so I figured this would be a good chance to do so. `slash` does not currently support all Slack slash command features. PRs are welcome, as well as any feedback on my probably-not-very-idiomatic Go code.

Example usage:
```
// Simple "echo command" to echo back a slack users command;
package main

import (
	"github.com/mfonda/slash"
	"log"
)

func main() {
	command := slash.NewCommand("/echo", "token", echo)
	slash.HandleCommand(command)
	log.Fatal(slash.ListenAndServeTLS(":443", "/path/to/cert", "/path/to/key"))
}

func echo(req *slash.Request) (*slash.Response, error) {
	return slash.NewInChannelResponse(req.Text, nil), nil
}
```
