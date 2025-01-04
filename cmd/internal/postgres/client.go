package postgres

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Client struct {
	// master / leader DB
	Leader *sql.DB // --> write (INSERT, UPDATE, DELETE)

	// write -> insert ke SQL -> bikin row baru & bikin log baru
	// INSERT INTO blog (title, content) VALUES ('title', 'content')
	// 10ms
	// sync ke replica

	// replica / follower DB
	Follower *sql.DB // --> read (SELECT)
}

func NewClient() (*Client, error) {
	dbURI := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("postgres", "postgres"),
		Host:   fmt.Sprintf("%s:%d", "localhost", 5432),
		Path:   "course-go",
	}

	leader, err := sql.Open("pgx", dbURI.String())
	if err != nil {
		return nil, err
	}

	return &Client{
		Leader: leader,
	}, nil
}
