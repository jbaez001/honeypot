# honeypot

Honeypot listens for TCP connections on the given ports and sends the data to the given server.

## Usage

```
Usage:
  honeypot [flags]

Flags:
      --config string        config file (default is $HOME/honeypot.yaml)
  -h, --help                 help for honeypot
      --honeypots strings    list of honeypots in the format protocol:port:enabled:fragile
      --name string          name of the honeypot (default "honeypot-1")
      --shout-urls strings   list of URLs to shout to

Example:
./honeypot  --name honeypot-1 \
            --shout-urls "telegram://token@telegram?chats=chat-id" \
            --honeypot telnet:23:true:false \
            --honeypot http:80:true:false   \
            --honeypot https:443:true:false     
```

## Configuration

```yaml
# name that will be used in the logs when reporting connections
# <name> [protocol> <- <ip>:<port>
name: "honeypot-1"

# format: https://containrrr.dev/shoutrrr/v0.8/services/overview/
shout_urls:
  - "telegram://token@telegram?chats=chat-idtelegram://token@telegram?chats=chat-id"

# honeypots
honeypots:
  - protocol: ftp
    port: 21
    enabled: true
    fragile: false
  - protocol: http
    port: 80
    enabled: true
    fragile: false
  - protocol: https
    port: 443
    enabled: true
    fragile: false
  - protocol: mysql
    port: 3306
    enabled: true
    fragile: false
  - protocol: telnet
    port: 23
    enabled: true
    fragile: false
  - protocol: postgres
    port: 5432
    enabled: true
    fragile: false
  - protocol: redis
    port: 6379
    enabled: true
    fragile: false
  - protocol: mongodb
    port: 27017
    enabled: true
    fragile: false
  - protocol: vnc
    port: 5900
    enabled: true
    fragile: false
```