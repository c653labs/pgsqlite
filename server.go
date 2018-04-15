package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/c653labs/pgproto"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	l, err := net.Listen("tcp", ":5432")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			defer conn.Close()
			err := handle(conn)
			if err != nil {
				log.Printf("error handling client: %#v", err)
			}
		}(conn)
	}
}

func handle(conn net.Conn) error {
	startup, err := pgproto.ParseStartupMessage(conn)
	if err != nil {
		return err
	}

	if startup.SSLRequest {
		pgproto.WriteMessage(&pgproto.Error{
			Severity: []byte("FATAL"),
			Message:  []byte("server does not support SSL, but SSL was requested"),
		}, conn)
		return fmt.Errorf("SSL not currently supported")
	}

	db, err := openDb(string(startup.Options["database"]))
	if err != nil {
		return err
	}
	defer db.Close()

	pgproto.WriteMessages([]pgproto.Message{
		&pgproto.AuthenticationRequest{
			Method: pgproto.AuthenticationMethodOK,
		},
		&pgproto.ReadyForQuery{
			Status: pgproto.READY_IDLE,
		},
	}, conn)

	for {
		msg, err := pgproto.ParseClientMessage(conn)
		if err != nil {
			return err
		}

		var msgs []pgproto.Message
		switch m := msg.(type) {
		case *pgproto.Termination:
			return nil
		case *pgproto.SimpleQuery:
			msgs = handleQuery(m, db)
		default:
			msgs = append(msgs, &pgproto.Error{
				Severity: []byte("FATAL"),
				Message:  []byte("unsupported client request"),
			})
		}

		msgs = append(msgs, &pgproto.ReadyForQuery{
			Status: pgproto.READY_IDLE,
		})
		pgproto.WriteMessages(msgs, conn)
	}
}

func handleQuery(msg *pgproto.SimpleQuery, db *sql.DB) (msgs []pgproto.Message) {
	rows, err := db.Query(string(msg.Query))
	if err != nil {
		msgs = append(msgs, &pgproto.Error{
			Severity: []byte("FATAL"),
			Message:  []byte(err.Error()),
		})
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		msgs = append(msgs, &pgproto.Error{
			Severity: []byte("FATAL"),
			Message:  []byte(err.Error()),
		})
		return
	}

	cFields := make([]pgproto.RowField, len(cols))
	for i, c := range cols {
		cFields[i].ColumnName = []byte(c)
	}
	msgs = append(msgs, &pgproto.RowDescription{
		Fields: cFields,
	})

	vals := make([]interface{}, len(cols))
	for i, _ := range cols {
		vals[i] = new(sql.RawBytes)
	}

	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			msgs = append(msgs, &pgproto.Error{
				Severity: []byte("FATAL"),
				Message:  []byte(err.Error()),
			})
			return
		}

		fields := make([][]byte, len(vals))
		for i, v := range vals {
			fields[i] = []byte(*(v.(*sql.RawBytes)))
		}

		msgs = append(msgs, &pgproto.DataRow{
			Fields: fields,
		})
	}

	err = rows.Err()
	if err != nil {
		msgs = append(msgs, &pgproto.Error{
			Severity: []byte("FATAL"),
			Message:  []byte(err.Error()),
		})
		return
	}

	msgs = append(msgs, &pgproto.CommandCompletion{
		Tag: msg.Query,
	})
	return
}

func openDb(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("./%s.db", name))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
