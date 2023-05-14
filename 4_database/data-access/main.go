package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var conn *pgx.Conn
var albums []Album

func main() {
	var err error
	var albums []Album

	os.Setenv("DATABASE_URL", "postgres://admin:admin@localhost:5432/godb")
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connect Success")

	albums, err = albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// Hard-code ID 2 here to test the query.
	alb, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)
}

// query for multiple rows
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

// query for one row with specific id
func albumByID(id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := conn.QueryRow(context.Background(), "SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
//func addAlbum(alb Album) (int64, error) {
//	result, err := conn.Exec(context.Background(), "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3)", alb.Title, alb.Artist, alb.Price)
//	if err != nil {
//		return 0, fmt.Errorf("addAlbum: %v", err)
//	}
//	id, err := result.LastInsertId()
//	if err != nil {
//		return 0, fmt.Errorf("addAlbum: %v", err)
//	}
//	return id, nil
//}
