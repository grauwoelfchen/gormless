# Gormless

A migrator command built on top of [Gormigrate](
https://github.com/go-gormigrate/gormigrate) using [GORM](
https://github.com/jinzhu/gorm) for migration files built as Go's [
plugin](https://golang.org/pkg/plugin/).

Inspired by [diesel-cli](
https://github.com/diesel-rs/diesel/tree/master/diesel_cli), but for GORM.

## Repository

https://gitlab.com/grauwoelfchen/gormless

## Install

```zsh
# see build section for tags
% go get -tags=<TAG> gitlab.com/grauwoelfchen/gormless/cmd/gormless
```

See examples.

### Build

```zsh
% go build -tags="mssql"
% go build -tags="mysql"
% go build -tags="postgres"
% go build -tags="sqlite"
```

## Usage

At first, you need to create migration files as Go's plugins.

```zsh
% cat migration/20190404_create_users/up.go
package main

import "github.com/jinzhu/gorm"

// Up ...
func Up(tx *gorm.DB) error {
  ...
}
```

See [gormigrate](https://github.com/go-gormigrate/gormigrate).  
And build them as a plugin.

```zsh
% cd ./migration/20190404_create_users/
% ls
down.go up.go
% go build -buildmode=plugin
```

Gormless recognize that `.so` plugins.


## Examples

```txt
[env variables] gormless <action> [<option>]
```

The `action` must be one of `commit`, `migrate`, `revert`, `rollback` or
`version`.  
Using a flag `-migration-directory` or through an environment variable
`MIGRATION_DIRECTORY=...`, you can set path to the directory contains your
migration files. The flag has higher priority than the environment variable.  
`DATABASE_URL=...` Should be started with `mssql://`, `mysql://`, `postgres://`
or `sqlite://`. Default is `:memory:` (sqlite).

It looks like so:

```zsh
% DATABASE_URL=... gormless commit
% DATABASE_URL=... gormless revert
```

See `gormless -h` about details.
