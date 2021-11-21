package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gocql/gocql"
	"main.go/serviceimpl"
)

//For Casendra staging
const ip = "34.93.201.115"
const username = "cassandra"
const password = "$Knit@D0rkKer4@%"
const keyspace = "upi"

var cluster *gocql.ClusterConfig

func init() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// defer file.Close()
	log.SetOutput(file)
}

func main() {
	fmt.Println("UPI Golang 1.0")
	cluster = gocql.NewCluster(ip)
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.Keyspace = keyspace
	cluster.ConnectTimeout = time.Second * 100
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: username, Password: password}
	var err error
	session, err := cluster.CreateSession()
	if err != nil {
		log.Println(err)

	}

	conn := &serviceimpl.Server{Session: session}

	mux := serviceimpl.NewRouter(conn)

	http.Handle("/", mux)
	http.ListenAndServe(":8090", nil)

}
