# Course Notes

## Postgres & SQL

- Create table pattern
```sql
CREATE TABLE table_name (
    column_name TYPE ...CONSTRAINTS,
    -- Some types require additional args ex. VARCHAR(255)
    column_name TYPE(ARGS),
);
```
- Insert pattern (we can specify the columns we want to insert into & the order)
```sql
INSERT INTO table_name (column1, column2, ...)
VALUES (value1, value2, ...);
```

## Docker Commands
- List all the running containers
```bash
docker ps
```
- Launch PSQL in the db container
```bash 
docker compose exec -it db psql -U baloo -d lenslocked
```

## Misc.

### Ideas
- Create a table for users with 1M rows, get them all, and render on UI using HTMX, server side pagination in Go. Compare speeds with a Flask server doing the same