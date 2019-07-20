# SQLite Example

## Build

Install gormless with `-tags=sqlite`.

```zsh
% GO111MODULE=on go get -tags=sqlite \
  gitlab.com/grauwoelfchen/gormless/cmd/gormless
```

Build migrations.

```zsh
% make
```

And, just run it.

```
% DATABASE_URL=... gormless commit
```
