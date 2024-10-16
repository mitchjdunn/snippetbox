package main

import (
	"database/sql"
	"flag"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // New import
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"snippetbox.mitchymit.ch/internal/models"
)

// Define an application struct to hold the applicaiton-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses.
type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	// Define a new command line flag with the na me 'addr', a default value of :4000
	// and some short hel p text explaining what the flag controls. The value of the flag
	// will be stored in the addr variable at runtime
	addr := flag.String("addr", ":4000", "Http network address")
	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// environment vairables are possible, not recommended due to lack of
	// built in defaults and abilitity to get different types.
	//addr := os.Getenv("SNIPPETBOX_ADDR")

	// Use the slog.New() function to intialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	// in this case it's useless, when method closes app is done and all resources released
	// but it's good practice to always close connections
	defer db.Close()

	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()

	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger).
	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	// Importantly, we use the flag.Parse() function to parse the command line flag.
	// this reads in the command line flag value and assigns it to the addr
	// variable. you need to call this *before* you use the addr variable
	// otherwise it will always contain the default value.  If any errors are
	// encountered during parsing the application will be terminated
	flag.Parse()

	// The value returned from flag.String() function is a pointer to the flag
	// value, not the value itself. so in this code, that means the addr variable
	// is actually a pointer, and we need to dereference it (i.e. prefix it with the
	// * symbol) before using it. Note that we're using the log.Printf()
	// function to interpolate the address with the log message
	logger.Info("starting server", "addr", *addr)

	// Pass in the dereferenced addr pointer to http.listenAndServer()
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
