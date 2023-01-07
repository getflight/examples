package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/getflight/core"
	"github.com/getflight/core/configuration"
	_ "github.com/go-sql-driver/mysql" //needed for mysql sql driver
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file" //needed for migrations
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type Note struct {
	Id        string
	Name      string
	Note      string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

func main() {
	log.SetLevel(log.DebugLevel)

	// Initialize core, this will initialize configuration and allow us to fetch secrets from AWS
	err := core.Init()

	if err != nil {
		fmt.Printf("error initializing core: %s\n", err)
		os.Exit(1)
	}

	// Database management
	connection := openDatabaseConnection()
	defer closeDatabaseConnection(connection)
	migrateDatabase(connection)

	// Endpoint to fetch notes from database
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		notes := fetchNotes(connection)

		err := json.NewEncoder(writer).Encode(notes)

		if err != nil {
			fmt.Printf("error writing response: %s\n", err)
		}
	})

	handler := os.Getenv("_HANDLER")

	if handler == "" {
		// Running locally
		err := http.ListenAndServe(":8080", nil)

		if err != nil {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}

	} else {
		// Running on lambda
		lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
	}
}

func openDatabaseConnection() *sql.DB {
	fmt.Println("opening database connection")

	dataSourceName := "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=true&tls=false"

	username := configuration.GetString("databases.example.username") // DATABASES_EXAMPLE_USERNAME
	password := configuration.GetString("databases.example.password") // DATABASES_EXAMPLE_PASSWORD
	host := configuration.GetString("databases.example.host")         // DATABASES_EXAMPLE_HOST
	port := configuration.GetString("databases.example.port")         // DATABASES_EXAMPLE_PORT
	name := configuration.GetString("databases.example.name")         // DATABASES_EXAMPLE_NAME

	datasource := fmt.Sprintf(dataSourceName, username, password, host, port, name)

	fmt.Printf("connecting with datasource: %s\n", datasource)

	connection, err := sql.Open("mysql", datasource)

	if err != nil {
		fmt.Printf("error connecting to database: %s\n", err)
		os.Exit(1)
	}

	connection.SetConnMaxLifetime(time.Minute * 3)
	connection.SetMaxOpenConns(10)
	connection.SetMaxIdleConns(10)

	return connection
}

func closeDatabaseConnection(connection *sql.DB) {
	fmt.Println("closing database connection")

	err := connection.Close()

	if err != nil {
		fmt.Printf("error closing database connection: %s\n", err)
		os.Exit(1)
	}
}

func migrateDatabase(connection *sql.DB) {
	fmt.Println("migrating database")

	driver, err := mysql.WithInstance(connection, &mysql.Config{})

	if err != nil {
		fmt.Printf("error initializing driver: %s\n", err)
		os.Exit(1)
	}

	mi, err := migrate.NewWithDatabaseInstance(
		"file://resources/migrations",
		"mysql",
		driver,
	)

	if err != nil {
		fmt.Printf("error initializing migrations: %s\n", err)
		os.Exit(1)
	}

	err = mi.Up()

	if err != nil && err != migrate.ErrNoChange {
		fmt.Printf("error running migrations: %s\n", err)
		os.Exit(1)
	}
}

func fetchNotes(connection *sql.DB) []Note {
	rows, err := connection.Query("SELECT * from notes")

	if err != nil {
		fmt.Printf("error querying database: %s\n", err)
		os.Exit(1)
	}

	var notes []Note

	for rows.Next() {
		var note Note

		if err := rows.Scan(&note.Id, &note.Name, &note.Note, &note.CreatedAt, &note.UpdatedAt, &note.DeletedAt); err != nil {
			fmt.Printf("error scanning database row: %s\n", err)
			os.Exit(1)
		}

		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("error iterating database row: %s\n", err)
		os.Exit(1)
	}

	return notes
}
