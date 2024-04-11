# goproject
Go, Templ, Htmx

example GoTH todo app

- Go (1.22.1)
- PostgresSQL database

## nix flake 

    nix develop --impure --show-trace


On first run, your environment & database server will be initialized.
Subsequently, the database server will be started/stopped when you enter/exit the development shell.


## hot reload templ components

    templ generate --watch --proxy="http://localhost:3000" --cmd="go run ."


## easy portable binary

    go build .


point to database & run

    export PGHOST=127.0.0.1 PGPORT=5432 PGUSER=postgres PGPASSWORD=secret DBNAME=appdata

    ./goproject


## external resources

[Routing Enhancements for Go 1.22](https://go.dev/blog/routing-enhancements)

[templ](https://templ.guide)

[htmx](https://htmx.org/)

[Nix](https://nix.dev/)

[devenv](https://devenv.sh/)
