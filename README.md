# myrepos

Quickly establish and update clones of remote git
repositories on local disk, using a configuration file.

Install:
```
go install github.com/monopole/myrepos
```

Usage:
```
myrepos [{path/to/config/file}]
```

By default, the path to the config file is `$HOME/.myrepos.yml`, and as the
name suggests, it's a YAML file.

For an example, see [example_myrepos.yml](./example_myrepos.yml).

