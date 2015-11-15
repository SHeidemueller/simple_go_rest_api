package main

import (
	"log"
	"strconv"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DEFAULT_PORT = "8080"
	DEFAULT_HOST = "localhost"
)

type Dog struct {
	Id    bson.ObjectId `json:"id" bson:"_id"`
	Name  string        `json:"name"`
	Owner string        `json:"owner"`
}

var Session *mgo.Session
var Collection *mgo.Collection
var Dogs = make([]Dog, 0)
var id int = 0

func main() {

	//init env variables
	appEnv, err := cfenv.Current()

	//set port and host
	var host, port, uri string

	if err != nil {
		port = DEFAULT_PORT
		host = DEFAULT_HOST
		uri = "mongodb://bluemix:bluemix@candidate.21.mongolayer.com:11022/dogpark"
		log.Print("local")
		log.Print(uri)
	} else {
		port = strconv.Itoa(appEnv.Port)
		host = appEnv.Host
		//	get Mongo Credentials from Env
		dbService, err := appEnv.Services.WithName("mongodb")
		if err != nil {
			log.Printf("No database info found\n")
			return
		}
		cred := dbService.Credentials

		uri = "mongodb://" + cred["user"].(string) + ":" + cred["password"].(string) + "@" + cred["uri"].(string) + ":" + cred["port"].(string) + "/dogpark"
		log.Print(uri)
		log.Print("cloud")

	}

	//connect to mongo

	Session, err := mgo.Dial(uri)

	if err != nil {
		log.Print("can't conntect to db")
		panic(err)
	}
	defer Session.Close()

	Session.SetSafe(&mgo.Safe{})

	//set target db
	Collection = Session.DB("dogpark").C("dogs")

	//initialize some data
	err = Collection.Insert(&Dog{bson.NewObjectId(), "Buffy", "Susann"})

	if err != nil {
		panic(err)
	} else {
		log.Print("all good on insert")
	}

	//set up router
	router := gin.Default()

	router.Static("/static", "static")

	router.GET("/dogs", gettingAll)
	router.GET("/dogs/:id", gettingOne)
	router.POST("/dogs", posting)
	router.PUT("/dogs/:id", putting)
	router.DELETE("/dogs/:id", deleting)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(308, "/static")
	})

	// Listen and serve on port
	router.Run(host + ":" + port)
}

func gettingAll(c *gin.Context) {

	result := []Dog{}
	err := Collection.Find(nil).All(&result)
	if err != nil {
		log.Printf("error fetching %v", err)
	} else {
		c.JSON(200, result)
	}
}

func gettingOne(c *gin.Context) {

	id := c.Param("id")
	log.Print(id)

	//check if valid id
	if bson.IsObjectIdHex(id) == false {
		c.JSON(400, gin.H{"error": "ID not valid"})
		return
	}

	bId := bson.ObjectIdHex(id)

	//fetch dog
	result := Dog{}
	err := Collection.FindId(bId).One(&result)

	if err != nil {
		c.JSON(400, gin.H{"error": "id not found"})
	} else {
		c.JSON(200, result)
	}
}

func posting(c *gin.Context) {
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"name\": \"Lassie\", \"owner\": \"Joe Carraclough\" }" /somePost
	var input Dog
	c.Bind(&input)
	log.Print(c)

	err := Collection.Insert(&Dog{bson.NewObjectId(), input.Name, input.Owner})

	if err != nil {
		panic(err)
	} else {
		log.Print("all good on insert")
		log.Print(input)
		c.JSON(201, "dog inserted")
	}
}

func deleting(c *gin.Context) {

	id := c.Param("id")
	log.Print(id)

	//check if valid id
	if bson.IsObjectIdHex(id) == false {
		c.JSON(400, gin.H{"error": "ID not valid"})
		return
	}

	bId := bson.ObjectIdHex(id)

	err := Collection.RemoveId(bId)

	if err != nil {
		c.JSON(400, "id not found")
	} else {
		c.JSON(200, "dog deleted")
	}

}

func putting(c *gin.Context) {

	var input Dog
	c.Bind(&input)
	log.Print(c)

	id := c.Param("id")
	log.Print(id)

	//check if valid id
	if bson.IsObjectIdHex(id) == false {
		c.JSON(400, gin.H{"error": "ID not valid"})
		return
	}

	bId := bson.ObjectIdHex(id)

	err := Collection.UpdateId(bId, &Dog{bId, input.Name, input.Owner})

	if err != nil {
		c.JSON(400, "id not found")
	} else {
		c.JSON(200, "dog updated")
	}
}
