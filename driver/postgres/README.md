# PostgreSQL Driver

* Runs migrations in transcations.
  That means that if a migration failes, it will be safely rolled back.
* Tries to return helpful error messages.
* Stores migration version details in table ``schema_migrations``.
  This table will be auto-generated.


## Usage

```bash
sqltractor-cli -url postgres://user@host:port/database -path ./db/migrations create add_field_to_table
sqltractor-cli -url postgres://user@host:port/database -path ./db/migrations up
sqltractor-cli -url postgres://user@host:port/database?search_path=name -path ./db/migrations up # with custom search_path
sqltractor-cli help # for more info

## Authors

* Matthias Kadenbach, https://github.com/mattes