# HomeTab

Self hosted life/home assistant.  
Designed to run on NAS Linux box with x64 CPU within secure network.  

**Detailed docs are WIP**  
If you have any questions, just create an issue, I'll be glad to help :)


## Requirements

- x86-64 CPU (AMD/Intel)
- Linux (Windows support in the future)
- MariaDB / MySQL
- Redis

## How to run

### Docker

Basic example:
```
docker run -p "3000:3000" quay.io/systemz/hometab
```
You should configure DB credentials via ENV settings. 
All ENV variables can be found in `internal/config/config.go`

If configured properly you should be able to reach web interface

### Dev setup

```bash
# backend
./Taskfile.sh install-tools
./Taskfile.sh dev-backend
./Taskfile.sh dev-seed

# frontend
./Taskfile.sh dev-frontend

# android app
# open android dir with Android Studio
```

## Notes

**This is important**. Project is in early stage. Due to security reasons, please don't expose HomeTab to public internet.
Firewall or IPv4 NAT should be your good friend for now.

### Repo structure

Directory structure mostly assembles this:
- https://github.com/golang-standards/project-layout

This repo is a merge of three, previously separate projects that were closed source:
- tasktab
- tasktab-android
- gotag

Due to high context switching mental cost during development, I merged them into one.  
During merge I had an idea to share code with community and here we are :)


### Stack

Golang backend + Vue.js/Vuetify (material design) frontend.  
As an additional client, we have Android app too.

#### Android app

App doesn't use any Google dependent services.
Should run on any custom ROM, with our without root.  
Notifications are based on pretty good and free service https://pushy.me/.

Previously I wrote push notifications based on MQTT myself, it was working quite good. Should be findable by "RabbitMQ" in removed code.
With this setup, whole app worked fully offline in LAN.

Android app is not available on Play store and probably never will.  
Author want to spend more time on useful features for users than boxing with mostly nonsense store ToS changes and constantly changing app in strange ways.  
Just ignore all google nagging popups about "security" (in reality it is about their market domination), download .apk and install it. Done.  
I will try to put it on https://f-droid.org/ to make updates easier.

### TODO

- few screens and webms from web gui and android app in readme
- integrate web templates within binary like migrations already are
- finish `cmd/agent` that will send events from user PCs to backend
- docker-compose for easier production setup
- doc pushy.me token requirement

## License

MIT
