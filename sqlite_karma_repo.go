package main

import (
	"database/sql"
	"io/ioutil"
	"path"
	"runtime"
	"strings"
	"time"
)

const InsertQuery = `
	INSERT INTO karma 
		(giver, receiver, count, channel) 
		VALUES 
		(?, ?, ?, ?)
`

const RankingQuery = `
	SELECT {kind} as user, sum(count) as count 
		FROM karma 
		WHERE ? < date and date <= ?
		GROUP BY {kind} 
		ORDER BY sum(count) desc
`

type SQLiteKarmaRepo struct {
	db *sql.DB
}

func getKarmaDb (dataDir string) (*sql.DB, error) {
	if dataDir == "" {
		dataDir = "."
	}
	db, err := sql.Open("sqlite3", path.Join(dataDir, "karma.db"))
	if err != nil {
		return nil, err
	}
	_, filename, _, _ := runtime.Caller(1)
	raw, err := ioutil.ReadFile(path.Join(path.Dir(filename), "create_table.sql"))
	if err != nil {
		return nil, err
	}
	query := string(raw)
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewSQLiteKarmaRepo(dataDir string) KarmaRepo {
	db, err := getKarmaDb(dataDir)
	if err != nil {
		panic(nil)
	}
	repo := &SQLiteKarmaRepo{db}
	return repo
}

func (s *SQLiteKarmaRepo) Save(karmaList []Karma) error {
	stmt, err := s.db.Prepare(InsertQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, k := range karmaList {
		_, err := stmt.Exec(k.Giver, k.Receiver, k.Count, k.Channel)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLiteKarmaRepo) Ranking(kind KarmaAggregateKind, from time.Time, to time.Time) (KarmaRanking, error) {
	ranking := KarmaRanking{ Kind: kind}
	ranking.Ranks = make([]KarmaAggregate, 0)
	query := strings.ReplaceAll(RankingQuery, "{kind}", string(kind))
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return ranking, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(from, to)
	if err != nil {
		return ranking, err
	}
	var user string
	var count float64
	for rows.Next() {

		err := rows.Scan(&user, &count)
		agg := KarmaAggregate{
			User:  user,
			Count: count,
		}
		if err != nil {
			return ranking, err
		}
		ranking.Ranks = append(ranking.Ranks, agg)
	}
	return ranking, nil
}
