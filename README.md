# Sonic Client

This is a WIP client for the [Sonic](https://github.com/valeriansaliou/sonic) search server.

## Example Usage

```go

channel, err := channel.NewInjestChannel("[::1]:1491")
if err != nil {
  log.Fatalf("Could not make new channel: %s", err.Error())
}

// Start an injest channel
channel.Start("SecretPassword")

// Ping the server
channel.Ping()

// Quit the channel
channel.Quit()

```