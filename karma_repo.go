package main

import (
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

type User struct {
	Id          string
	Name        string
	DisplayName string
	TeamId      string
}
type KarmaAggregateKind string

var KindGiver = KarmaAggregateKind("giver")
var KindReceiver = KarmaAggregateKind("receiver")

type KarmaAggregate struct {
	User  string
	Count float64
}

type KarmaRanking struct {
	Kind  KarmaAggregateKind
	Ranks []KarmaAggregate
}

type UserRepo interface {
	Put(user User) error
	GetById(id string) (User, error)
	GetByName(name string) (User, error)
}

type KarmaRepo interface {
	// カルマのリストを保存する
	Save([]Karma) error
	// 期間をしていしてカルマの集計を返す
	Ranking(kind KarmaAggregateKind, from time.Time, to time.Time) (KarmaRanking, error)
}

type UserKarmaRepo interface {
	UserRepo
	KarmaRepo
}
