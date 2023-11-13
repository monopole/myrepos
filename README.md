# myrepos

Quickly establish and update clones of remote git
repositories on local disk, using a YAML configuration file.

Install:
```
go install github.com/monopole/myrepos@latest
```

Usage:
```
myrepos [{path/to/config/file}]
```

By default, the path to the configuration file is `$HOME/.myrepos.yml`.

Example configuration file: [example_myrepos.yml](example_myrepos.yml)

Detailed explanation of configuration fields: [config.go](internal/config/config.go)
