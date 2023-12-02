# myrepos

Clone, rebase or just report the status of remote git
repositories on local storage, using a YAML configuration file.

Install:
```
go install github.com/monopole/myrepos@latest
```

Usage:
```
myrepos [{configurationFile}...]
```

Specify one or more configuration files
as arguments (e.g. [myrepos_example.yml](myrepos_example.yml)).

If no argument is specified, an attempt
will be made to read `$HOME/.myrepos.yml`.

Detailed explanation of configuration fields: [config.go](internal/config/config.go)
