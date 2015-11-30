# Sqlite3 Driver

* Runs migrations in transcations.
  That means that if a migration failes, it will be safely rolled back.
* Tries to return helpful error messages.
* Stores migration version details in table ``schema_migrations``.
  This table will be auto-generated.


## Usage

```bash
sqltractor-cli -url sqlite3://database.sqlite -path ./db/migrations create add_field_to_table
sqltractor-cli -url sqlite3://database.sqlite -path ./db/migrations up
sqltractor-cli help # for more info
```

## Authors

* Matthias Kadenbach, https://github.com/mattes
* Caesar Wirth, https://github.com/cjwirth