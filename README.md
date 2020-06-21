# AMBARITA
Simple Article Management

## Database Migrations
Using CLI version of https://github.com/golang-migrate/migrate

* [Installation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
* Create migration file command ```migrate create -ext sql -dir migrations [NAME]``` example of NAME = create_table_articles
* Migrate changes command ```migrate -source file://migrations -database "mysql://[DB_USERNAME]:[DB_PASSWORD]@tcp([DB_HOST])/[DB_DATABASE]" up```
* Should table need data, provide table seeder within migration file. See `migrations/20200403154510_create_articles_table.up.sql` for example.

## Preparation

Modify the `config.json.example` file and rename to `config.json` on folder `configs`

### Running

To start Server, run:

```bash
go mod vendor
make engine
./engine
```

## Available Endpoints
* Fetch Article
* Get Article by ID/Slug
* Store Article
* Update Article
* Delete Article
