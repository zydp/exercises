package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"time"
)

type Record struct {
	X int
	Y int
}

func main() {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:     []string{"ip:port"},
		Timeout:   20 * time.Second,
		Database:  "demodb",
		Username:  "username",
		Password:  "password",
		PoolLimit: 32,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		fmt.Printf("DialWithInfo: %s\n", err)
		os.Exit(100)
	}

	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	//err = collection.Find(bson.M{"x": 1}).All(&recs)

	for i:=0; i<10; i++ {
		s2 := session.Clone()
		defer s2.Close()
		collection := s2.DB("go").C("daiping")

		var recs []Record
		err = collection.Find(bson.M{
			"$or": []bson.M{
				bson.M{"x": 4},
				bson.M{"x": 3},
			},
		}).All(&recs)

		if err != nil {
			fmt.Println(err)
			os.Exit(100)
		}

		for x, y := range recs {
			fmt.Println(x, y)

		}
	}
	<-time.After(5 * time.Second)	//wait five seconds
}

