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

For an example, see [example_myrepos.yml](example_myrepos.yml).

For a detailed explanation of configuration fields, see [config.go](internal/config/config.go).
