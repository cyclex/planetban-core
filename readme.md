# Planet Ban Core Services

## Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Configuration](#configuration)

## Features

There is 2 command in this core services
1. server
    supporting rest endpoint for CMS
2. order
    supporting chatbot and shorten link for influencer

## Installation


```bash
# Clone the repository
git clone https://github.com/cyclex/planetban-core.git

# Navigate to the project directory
cd planet-ban

# Install dependencies
make build
```

## Usage

```bash
# Run the application
# Server apps
./engine server -p :8081 -c config.json -d true

NAME:
   server server - start cms and chatbot service

USAGE:
   server server [command options] [arguments...]

OPTIONS:
   --port value, -p value    Listen to port
   --config value, -c value  Load configuration file
   --debug, -d               Debug mode (default: false)
   --help, -h                show help

# Order apps
./engine order -p :8082 -c config.json -d true

NAME:
   server order - start webhook

USAGE:
   server order [command options] [arguments...]

OPTIONS:
   --port value, -p value    Listen to port
   --config value, -c value  Load configuration file
   --debug, -d               Debug mode (default: false)
   --help, -h                show help
```

## Configuration

```json
{
  "log": {
    "maxsize": 10,
    "maxbackups": 10
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "planetban",
    "password": "planetban123",
    "name": "chatbot"
  },
  "chatbot": {
    "host": "https://api.coster.id/bot/webhook",
    "account_id": "12927d2adb9c11ee8741fff7a41ee0e0",
    "division_id": "95a7c69e45a011eeab31eb2138e54e5e",
    "waba_account_number":"6282311333723",
    "access_token":"accessTokenChatbot"
  },
  "url_host_influencer":"https://planetbancore.coster.id/go"
}
```

## Deployment
```bash
# root application folder
# password planetban123
planetban@planetban-core:~$ pwd
/home/planetban

# logs application folder
/home/planetban/.planetban

# run by systemd
# Server apps
sudo systemctl status planetban-server

# Order apps
sudo systemctl status planetban-webhook
```
