# What is CatBot about?
CatBot is a discord bot made for CatEmpire server. It's just a cute cat that can meow whenever see a cat in text message.
## Usage:
- Write something with **cat** in your text message(it doesn't matter, is capital or not) in text channel
- join a VC and write `!Meow` in a text channel

## Under Development Features:
- [ ] Command-line Interface for adding command 
- [ ] Database for persistent storage
- [ ] Play random Cat sounds
- [ ] Show random Cat Picture

# Installation
- make a `config.json` file with the `Token` and `BotPrefix`
```json
 {
	"Token": "",	
	"BotPrefix": "$"
}
```
- put `config.json` in the root directory of project
- use `make run` or `make build` to run the bot
- `make help`: see more command
