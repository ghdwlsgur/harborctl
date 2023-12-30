# harborctl

## How to Install

```bash
# [install]
brew tap ghdwlsgur/harborctl
brew install harborctl

# [upgrade]
brew tap ghdwlsgur/harborctl
brew upgrade harborctl

# latest version (2023.12.30)
$ harborctl --version
harborctl version 1.0.1

# [refresh]
brew update --verbose
```

## How to Use

```bash
# When a login attempt is successful, the ID and
# password are saved in the .harborctl/credentials file
harborctl login

harborctl create [name] -d [duration] -e [description]
harborctl create albert -d 30 -e albert

harborctl update -d [duration]
harborctl update -d 100

harborctl delete
harborctl search
```

## Release

```bash
make release TAG=v1.0.1
```

[![asciicast](https://asciinema.org/a/IayRztV5EVPLu7BIf6PYaBMxO.svg)](https://asciinema.org/a/IayRztV5EVPLu7BIf6PYaBMxO)
