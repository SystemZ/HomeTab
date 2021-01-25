# HomeTab

Selfhosted life/home assistant.  
Merge of three projects:
- tasktab
- tasktab-android
- gotag

WIP: merge code of gotag

## Requirements to run

- x86-64 CPU (AMD/Intel)
- Linux (Windows support in the future)
- MariaDB / MySQL
- Redis

## How to run

```bash
# backend
./Taskfile.sh install-tools
./Taskfile.sh dev-backend
./Taskfile.sh dev-seed
# frontend
./Taskfile.sh dev-frontend
```

## TODO

- integrate web templates within binary
- finish `cmd/agent` that will send events from user PCs to backend

## Repo paths

- https://github.com/golang-standards/project-layout
