# A. Create database and user

## A1. Create a new user and give proper role

To create a user, developer must connect to the existing data with a user having admin access right.

```text
CREATE USER name [ [ WITH ] option [ ... ] ]

where option can be:
    
      SYSID uid 
    | CREATEDB | NOCREATEDB
    | CREATEUSER | NOCREATEUSER
    | IN GROUP groupname [, ...]
    | [ ENCRYPTED | UNENCRYPTED ] PASSWORD 'password'
    | VALID UNTIL 'abstime' 
```

## A2. Create datbase and link it with some users

```text 
CREATE DATABASE name
    [ WITH ] [ OWNER [=] user_name ]
           [ TEMPLATE [=] template ]
           [ ENCODING [=] encoding ]
           [ STRATEGY [=] strategy ] ]
           [ LOCALE [=] locale ]
           [ LC_COLLATE [=] lc_collate ]
           [ LC_CTYPE [=] lc_ctype ]
           [ ICU_LOCALE [=] icu_locale ]
           [ LOCALE_PROVIDER [=] locale_provider ]
           [ COLLATION_VERSION = collation_version ]
           [ TABLESPACE [=] tablespace_name ]
           [ ALLOW_CONNECTIONS [=] allowconn ]
           [ CONNECTION LIMIT [=] connlimit ]
           [ IS_TEMPLATE [=] istemplate ]
           [ OID [=] oid ]
```

# B. Database Connection Steps

## B1. connect to existing database

```shell
psql -d admin -U odoo -W odoo
```

## B2. Create a new user which has right to create new database

```postgresql
CREATE USER admin WITH password 'admin' CREATEDB;
```

## B3. Create a new user which is associated with the selected user

```postgresql
CREATE DATABASE godb WITH OWNER = admin;
```

## B4. Log out postgre command prompt and login the user <mark>"admin"</mark> with database <mark>"godb"</mark>

```shell
psql -d godb -U admin -W admin
```

## B5. Create table

```postgresql
DROP TABLE IF EXISTS album;
```

## B6. Create a database in album table

```postgresql
CREATE TABLE album
(
    id     SERIAL        NOT NULL PRIMARY KEY,
    title  VARCHAR(128)  NOT NULL,
    artist VARCHAR(255)  NOT NULL,
    price  DECIMAL(5, 2) NOT NULL
);
```

## B7. Insert some data to the table

```postgresql
INSERT INTO album
    (title, artist, price)
VALUES ('Blue Train', 'John Coltrane', 56.99),
       ('Giant Steps', 'John Coltrane', 63.99),
       ('Jeru', 'Gerry Mulligan', 17.99),
       ('Sarah Vaughan', 'Sarah Vaughan', 34.98);
```

## B8. Verify the records are inserted correctly

```postgresql
select *
from album;
```

- Output:
  ![img.png](img.png)

# C. Use go language as database application

## C1. Official Documentation for go connection with database

In this [Github Page](https://github.com/golang/go/wiki/SQLDrivers), we can find different drivers of go language.
Simply click on some library you are interested in. For
example, we are going to use PostgreSQL to store our data, so you can follow this link to see the
example. [Postgres (pure Go)] (https://github.com/jackc/pgx)

## C2. Make connection

- Configuration
    - User, password, database_name, host and port are important component to set up a connection between go and
      database. Since we are using PostgreSQl, a database depends on TCP/IP connection, the configuration can be
      converted to a simple URL, similar to most of API connections.
    - URL connection format: postgres://{username}:{password}@{address:port}/{database_name}

```go
os.Setenv("DATABASE_URL", "postgres://admin:admin@localhost:5432/godb")
conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
if err != nil {
fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
os.Exit(1)
}
fmt.Println("Connect Success")
defer conn.Close(context.Background())
```

## C3. Management of Environment variables

To be honest, it is not safe and convinient to store import string information in a inner function variable.
In go, variables can be collected and save to environment variables.

## C4. Query for multiple rows

- Create database struct, similar to Jason or C++ struct

```go
type Album struct {
ID     int64
Title  string
Artist string
Price  float32
}
```

- Compute SQL to select some rows
```go
func albumsByArtist(name string) ([]Album, error) {
	var rows pgx.Rows
	var err error
	rows, err = conn.Query(context.Background(), "SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err = rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}
```