package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
	"github.com/lucsky/cuid"
)

type Contest struct {
	Title string
	Description string
	Link string
	Date string
}

func InitializeContests() []Contest {
	contests := []Contest{
        {
            Date:        "08.06.24",
            Title:       "Steinhuder Meer",
            Description: "Fun (nur 5km Laufen)",
            Link:        "https://www.steinhudermeer-triathlon.de/",
        },
        {
            Date:        "16.06.24",
            Title:       "Düsseldorf",
            Description: "SD/OD",
            Link:        "https://events.larasch.de/en/t3-triathlon-duesseldorf",
        },
        {
            Date:        "30.06.24",
            Title:       "Bochum",
            Description: "SD",
            Link:        "https://www.bochum-triathlon.de/",
        },
        {
            Date:        "13.07.24",
            Title:       "Hamburg",
            Description: "SD",
            Link:        "https://hamburg.triathlon.org/",
        },
        {
            Date:        "28.07.24",
            Title:       "Frankfurt",
            Description: "OD",
            Link:        "https://www.frankfurt-city-triathlon.de/",
        },
        {
            Date:        "01.09.24",
            Title:       "Duisburg",
            Description: "MD",
            Link:        "https://www.ironman.com/im703-duisburg-register",
        },
        {
            Date:        "08.09.24",
            Title:       "Köln",
            Description: "OD",
            Link:        "https://www.carglass-koeln-triathlon.de/",
        },
        {
            Date:        "06.10.24",
            Title:       "Köln",
            Description: "Marathon",
            Link:        "https://generali-koeln-marathon.de/",
	},
	}

	return contests
}

func initializeDatabase() {
	os.Remove("./contests.sqlite3")

	db, err := sql.Open("sqlite3", "./contests.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
		CREATE TABLE contests (
			id text not null primary key,
			title text not null,
			description text,
			link text,
			date text
		);
	`

	contests := InitializeContests()

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	insertStmt, err := db.Prepare("INSERT INTO contests(id, title, description, link, date) values(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer insertStmt.Close()

	for _, c := range contests {
		_, err = insertStmt.Exec(cuid.New(), c.Title, c.Description, c.Link, c.Date)
		if err != nil {
			log.Fatal(err)
		}
	} 
}

func loaddata() ([]Contest, error) {
	db, err := sql.Open("sqlite3", "./contests.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, description, link, date FROM contests")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var all []Contest
	for rows.Next() {

		var id string
		var contest Contest
		err = rows.Scan(&id, &contest.Title, &contest.Description, &contest.Link, &contest.Date)
		if err != nil {
			log.Fatal(err)
		}

		all = append(all, contest)
	}

	return all, nil
}

func main() {
	initializeDatabase()
	mux := http.NewServeMux()

	contests, err := loaddata()
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("GET /api/v1/contest/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contests)
	})

	mux.HandleFunc("GET /api/v1/contest/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprint(w, "handling task with id=%v\n", id)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
