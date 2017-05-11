# do-dyn-domain

[![Build Status](https://travis-ci.org/emersion/do-dyn-domain.svg?branch=master)](https://travis-ci.org/emersion/do-dyn-domain)

Digital Ocean dynamic domain updater. Works well with a UPnP router, integrates
with [doctl](https://github.com/digitalocean/doctl).

## Usage

Configure your `doctl` config file (on Unix, it's in `$HOME/.config/doctl/config.yaml`):

```yaml
access-token: MY_TOKEN
```

Run: `do-dyn-domain -domain example.org -record record-name`

## License

MIT
