package main

/**
 * All the code was put into one file for the sake of simplicity...
 * So, enjoy yourselves •-•
 */

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Config of database
type DataSource struct {
	Driver       string
	Host         string
	Port         string
	User         string
	Pass         string
	Database     string
	SSLmode      string
	Options      string
	Dns          string
	MaxPool      int
	MaxIdleConns int
	DB           *sqlx.DB
}

func (ds *DataSource) New() *DataSource {
	return &DataSource{
		Driver:       os.Getenv("POSTGRES_DRIVER"),
		Host:         os.Getenv("POSTGRES_HOST"),
		Port:         os.Getenv("POSTGRES_PORT"),
		Pass:         os.Getenv("POSTGRES_PASSWORD"),
		User:         os.Getenv("POSTGRES_USER"),
		Database:     os.Getenv("POSTGRES_DB"),
		SSLmode:      os.Getenv("POSTGRES_SSL_MDOE"),
		MaxPool:      1000,
		MaxIdleConns: 10,
	}
}

func (ds *DataSource) Handshake() string {
	handshake := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		ds.Host,
		ds.Port,
		ds.User,
		ds.Pass,
		ds.Database,
		ds.SSLmode,
	)
	return handshake
}

func (ds *DataSource) Pooling(conn chan error) {
	for {
		db, err := sqlx.Open(ds.Driver, ds.Dns)
		if err != nil {
			conn <- err
			time.Sleep(2 * time.Millisecond)
			continue
		}

		db.SetMaxOpenConns(ds.MaxPool)
		db.SetMaxIdleConns(ds.MaxIdleConns)
		db.SetConnMaxLifetime(time.Minute * 3)

		ds.DB = db
		conn <- nil
		return
	}
}

func (ds *DataSource) Conn() (*sqlx.DB, error) {
	if ds.DB == nil {
		return nil, errors.New("")
	}
	return ds.DB, nil
}

// Config of routes
var app Routes
var ds *DataSource

type Routes struct {
	Server *fiber.App
	RestV1 fiber.Router
}

func (app *Routes) Routing() *Routes {
	app.Server = fiber.New(fiber.Config{
		ServerHeader:          "Fiber",
		CaseSensitive:         true,
		StrictRouting:         true,
		DisableStartupMessage: true,
	})
	// app.Server.Use(logger_fiber.New())

	r := app.Server

	r.Get("/", func(c *fiber.Ctx) (err error) {
		c.Status(http.StatusOK).SendString("42")
		return
	})

	// Handlers and *
	r.Post("/clientes/:id/transacoes", func(ctx *fiber.Ctx) (err error) {
		id := ctx.Params("id")
		req := struct {
			Type  string `json:"tipo" db:"tipo"`
			Desc  string `json:"descricao" db:"descricao"`
			Value uint   `json:"valor" db:"valor"`
		}{}
		err = ctx.BodyParser(&req)

		if err != nil || (req.Value == 0) || (req.Type != "c" && req.Type != "d") || (req.Desc == "" || len(req.Desc) > 10) {
			return ctx.SendStatus(http.StatusUnprocessableEntity)
		}

		row := struct {
			ID      int `json:"-" db:"id"`
			Balance int `json:"saldo" db:"saldo"`
			Lmt     int `json:"limite" db:"limite"`
		}{}

		db, err := ds.Conn()
		if err != nil {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		tx, err := db.Begin()
		if err != nil {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		var query string
		switch req.Type {
		case "c":
			query = `
				UPDATE clientes c SET
					saldo = COALESCE(c.saldo, 0) + $2
				WHERE
					c.id = $1
				RETURNING
					c.id
				,	c.saldo
				,	c.limite
			`
		case "d":
			query = `
				UPDATE clientes c SET
					saldo = COALESCE(c.saldo, 0) - $2
				WHERE
					c.id = $1
				RETURNING
					c.id
				,	c.saldo
				,	c.limite
			`
		}
		err = tx.QueryRow(query, id, req.Value).Scan(&row.ID, &row.Balance, &row.Lmt)
		if req.Type == "d" && row.Balance < -row.Lmt {
			tx.Rollback()
			return ctx.SendStatus(http.StatusUnprocessableEntity)
		}

		if row.ID == 0 {
			tx.Rollback()
			return ctx.SendStatus(http.StatusNotFound)
		}

		if err != nil && err.Error() != "" {
			tx.Rollback()
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		query = `
			INSERT INTO transacoes (id_cliente, valor, tipo, descricao) VALUES ($1, $2, $3, $4);
		`
		tx.Exec(query, row.ID, req.Value, req.Type, req.Desc)
		// if _, err = tx.Exec(query, row.ID, req.Value, req.Type, req.Desc); err != nil {
		// 	tx.Rollback()
		// 	return ctx.SendStatus(http.StatusUnprocessableEntity)
		// }

		tx.Commit()
		return ctx.Status(http.StatusOK).JSON(row)
	})

	r.Get("/clientes/:id/extrato", func(ctx *fiber.Ctx) (err error) {
		id := ctx.Params("id")
		row := struct {
			ID               int              `json:"-" db:"id"`
			Balance          json.RawMessage  `json:"saldo" db:"saldo"`
			LastTransactions *json.RawMessage `json:"ultimas_transacoes" db:"ultimas_transacoes"`
		}{}

		query := `
			WITH
			last AS (
				SELECT
					t.id_cliente
				,	t.valor
				,	t.tipo
				,	t.descricao
				,	t.realizada_em
				FROM transacoes t
				WHERE 
					t.id_cliente = $1
				ORDER BY 
					t.realizada_em DESC
				LIMIT 10
			)
			SELECT
				c.id
			,	json_build_object(
					'total', COALESCE(c.saldo, 0)
				,	'data_extrato', now()
				,	'limite', c.limite
				) saldo
			,	CASE
					WHEN COUNT(l.id_cliente) = 0 THEN '[]'::json
					ELSE json_agg(
						json_build_object(
							'valor', l.valor
						,	'tipo', l.tipo
						,	'descricao', l.descricao
						,	'realizada_em', l.realizada_em
						)
					) 
				END AS ultimas_transacoes
			FROM clientes c
			LEFT JOIN last l ON l.id_cliente = c.id
			WHERE c.id = $1
			GROUP BY c.id
		`
		db, err := ds.Conn()
		if err != nil {
			return err
		}
		db.Get(&row, query, id)

		if row.ID == 0 {
			return ctx.SendStatus(http.StatusNotFound)
		}

		return ctx.Status(http.StatusOK).JSON(row)
	})

	return app
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error loading .env:", err)
	}

	ds = ds.New()
	ds.Dns = ds.Handshake()
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered. Error:\n", r)
		}
	}()

	r := app.Routing()

	conn := make(chan error)
	go ds.Pooling(conn)

	err := <-conn
	if err != nil {
		panic(err)
	}

	socketPath := os.Getenv("UNIX_SOCK_PATH")
	os.Remove(socketPath)

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err) // do not panic!
	}
	defer ln.Close()

	if err := r.Server.Listener(ln); err != nil {
		panic(err)
	}

	// to run local
	// r.Server.Listen(os.Getenv("PORT"))
}
