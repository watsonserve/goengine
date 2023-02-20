package goengine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbConf struct {
	User   string
	Passwd string
	Host   string
	Name   string
	Port   string
}

func ConnPg(config *DbConf) *sql.DB {

	pgurl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Passwd,
		config.Host,
		config.Port,
		config.Name,
	)
	db, err := sql.Open("postgres", pgurl)
	if nil != err {
		panic(err)
	}
	return db
}

type DAO struct {
	db      *sql.DB
	StmtMap map[string]*sql.Stmt
}

func InitDAO(db *sql.DB) *DAO {
	dao := &DAO{}

	dao.db = db
	dao.StmtMap = make(map[string]*sql.Stmt)

	return dao
}

func (this *DAO) Prepare(index string, query string) {
	stmt, err := this.db.Prepare(query)
	if nil != err {
		panic(err)
	}
	this.StmtMap[index] = stmt
}

// mongoDB

func ConnMongo(config *DbConf) *mongo.Database {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s",
		config.User,
		config.Passwd,
		config.Host,
		config.Port,
		config.Name,
	)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	db := client.Database(config.Name)
	if nil == db {
		panic(errors.New("MongoDB No DataBase"))
	}
	return db
}
