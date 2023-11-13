# myrepos

Clone, rebase or just report the status of remote git
repositories on local storage, using a YAML configuration file.

Install:
```
go install github.com/monopole/myrepos@latest
```

Usage:
```
myrepos [{configFile}]
```

By default, the path to the configuration file is `$HOME/.myrepos.yml`.

Example configuration file: [example_myrepos.yml](example_myrepos.yml)

Detailed explanation of configuration fields: [config.go](internal/config/config.go)
