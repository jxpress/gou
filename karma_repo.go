package main

import (
	"database/sql"
	"io/ioutil"
	"path"
	"runtime"
	"time"
)

type Karma struct {
	Id        int
	Giver     string
	Receiver  string
	Count     float64
	Channel   string
	CreatedAt time.Time
}

type KarmaAggregateKind string
var KindGive = "give"
var KindReceive = "receive"


type KarmaAggregate struct {
	User string
	Count float64
}

type KarmaRanking struct {
	ranks []KarmaAggregateKind
}

type KarmaRepo interface {
	// カルマのリストを保存する
	Save ([]Karma) error
	// 期間をしていしてカルマの集計を返す
	Ranking(kind KarmaAggregateKind, from time.Time, to time.Time) (KarmaRanking, error)
}

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

func (S *SQLiteKarmaRepo) Save([]Karma) error {
	panic("implement me")
}

func (S *SQLiteKarmaRepo) Ranking(kind KarmaAggregateKind, from time.Time, to time.Time) (KarmaRanking, error) {
	panic("implement me")
}
