package goengine

import (
    "fmt"
    "errors"
    "github.com/watsonserve/goutils"
    "net"
    "crypto/tls"
    "database/sql"
    _ "github.com/lib/pq"
    mgo "gopkg.in/mgo.v2"
)

func ConnPg(config map[string]string) *sql.DB {

    pgurl := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        config["DBUser"],
        config["DBPasswd"],
        config["DBHost"],
        config["DBPort"],
        config["DBName"],
    )
    db, err := sql.Open("postgres", pgurl)
    if nil != err {
        panic(err)
    }
    return db
}

type DAO struct {
    db *sql.DB
    StmtMap map[string]*sql.Stmt
}

func InitDAO(db *sql.DB) *DAO {
    dao := &DAO{}

    dao.db = db
    dao.StmtMap = make(map[string]*sql.Stmt)

    return dao
}

func (this *DAO) Prepare(index string, query string) {
    stmt ,err := this.db.Prepare(query)
    if nil != err {
        panic(err)
    }
    this.StmtMap[index] = stmt
}

// mongoDB

type MongoConf struct {
	User   string
	Passwd string
	Host   string
	Port   string
	Name   string
}

func ConnMongo(config *MongoConf) *mgo.Database {
    url := fmt.Sprintf(
        "mongodb://%s:%s@%s:%s/%s",
        config.User,
        config.Passwd,
        config.Host,
        config.Port,
        config.Name,
    )
    dialInfo, err := mgo.ParseURL(url)
    if nil != err {
        goutils.Printf("%s\n", err)
        return nil
    }

    dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
        conn, err := tls.Dial("tcp", addr.String(), &tls.Config{InsecureSkipVerify: true})
        if nil != err {
            goutils.Printf("connect mongodb: %s\n", err.Error())
        }
        return conn, err
    }

    sess, err := mgo.DialWithInfo(dialInfo)
    if nil != err {
        panic(err)
    }
    // defer sess.Close()

    db := sess.DB(config.Name)
    if nil == db {
        panic(errors.New("MongoDB No DataBase"))
    }
    return db
}
