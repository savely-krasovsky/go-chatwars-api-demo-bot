# go-chatwars-api-demo-bot
Simple Telegram Bot with demonstration of work [go-chatwars-api](https://github.com/L11R/go-chatwars-api).

Supports three commands:
* `/start` - small introduction
* `/auth` - begins auth process
* `/profile` - shows profile if auth was successful

Needs small `config.yml` file before executing:
```yaml
token: telegram_bot_api_token_from_botfather
log:
  type: production // or development
  level: 0 // where -1 is DEBUG
db:
  name: users.db // database name
cw:
  user: your_app_user
  password: your_app_password
```

To launch bot, you need:
1. `cd $GOPATH/src`
2. `git clone https://github.com/L11R/go-chatwars-api-demo-bot && cd go-chatwars-api-demo-bot`
3. `dep ensure` (you need to [install](https://github.com/golang/dep) `dep` before)
4. `go build -o bot`
5. Create config (`nano config.yml`, you know 😏)
6. `./bot`

Also you can make `systemd` daemon, but it's beyond the scope of this instruction.