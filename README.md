# mastodon group bot

This is a bot that implements group functionality in Mastodon.

This bot was based on the code developed for [mastodon-pug](https://github.com/m0t0k1ch1/mastodon-pug), where people sometimes message @mastodon, the official announcement account, but the maintainers are too lazy to actually monitor it.

The bot is written in Golang 1.11 and uses [go-mastodon](github.com/mattn/go-mastodon).


TODO:
The bot can do the following:
- respond to every toot sent to it by a non-follower with a predefined message mentioning the admins
    - regardless of the visibility setting of the response, the response is always sent as a DM. because of how DMs work, if the predefined message includes other peple's usernames, they'll also see the DM!
- if it receives a DM from a non-follower, it can forward the text of that DM to the admins
- it can boost toots by followers that mention it

The bot will not respond retroactively, i.e., the first time you run it, it will not respond to all the messages your account has received in the past.

# Configuration

The bot is configured in a toml file. Example in _config/sample.toml

All keys are mandatory. The first group contains information about connecting to the API and authenticating to it. The second group contains the autoresponder message and the usernames of the admins. The last group contains the path to the state file, which contains informations that lets the bot remember which messages it's already replied to (this cannot be empty, but the file doesn't have to exist the first time you run the bot).

# Installation

This should really be packaged as a proper Python package, but I haven't done that. If you want to run this bot:

```
# 1. clone this repo
git clone https://github.com/inhosin/go-mastodon-bot-autoresponder

# 2. install the dependencies
go get -u github.com/mattn/go-mastodon
go get -u github.com/spf13/viper

# 3. use tokentool to register the bot as an app on your server,
# then authenticate to it (don't worry, it's not hard, there's a nice
# interactive text interface)
написать свой tokentool

# 4. create a config file and edit it appropriately
cp _config/sample.toml config.toml
nano config.toml

# 5. Build
go buil

# 6. run the bot!
./mastodon-bot-autoresponder
```