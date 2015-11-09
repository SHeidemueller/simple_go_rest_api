package main

import (
	"os"
	//for extracting service credentials from VCAP_SERVICES
	//"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gin-gonic/gin"
	"strconv"
	"gopak.in/mgo.v2"
)

const (
	DEFAULT_PORT = "8080"
	DEFAULT_HOST = "localhost"
)

type Dog struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

var Dogs = make([]Dog, 0)
var id int = 0

func main() {

	var port string
	if port = os.Getenv("VCAP_APP_PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}

	var host string
	if host = os.Getenv("VCAP_APP_HOST"); len(host) == 0 {
		host = DEFAULT_HOST
	}

	//prefilling "Dog Park" on start up

	Dogs = append(Dogs, Dog{nextId(), "Buffy", "Susann"})
	Dogs = append(Dogs, Dog{nextId(), "Snoopy", "Charlie Brown"})

	router := gin.Default()

	router.LoadHTMLFiles("index.html")

	router.Static("/static", "static")

	router.GET("/dogs", gettingAll)
	router.GET("/dogs/:id", gettingOne)
	router.POST("/dogs", posting)
	router.PUT("/dogs/:id", putting)
	router.DELETE("/dogs/:id", deleting)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(308, "/static")
	})

	//	router.GET("/", func(c *gin.Context) {
	//		c.HTML(200, "index.html", nil)
	//	})

	// Listen and serve on port
	router.Run(host + ":" + port)

}

func nextId() int {
	id += 1
	return id
}

func gettingAll(c *gin.Context) {

	//serve /allGet

	c.JSON(200, Dogs)

}

func gettingOne(c *gin.Context) {

	id := c.Param("id")
	dogId, err := strconv.Atoi(id)

	if err != nil {
		panic(err)
	}

	content, index := getDogbyID(dogId)

	//serve /oneGet/id
	if index == -1 {
		c.JSON(422, gin.H{"error": "index not found"})
	} else {

		c.JSON(200, content)
	}
}

//search through Dogs for correct Dog and Index
func getDogbyID(id int) (Dog, int) {
	for i, e := range Dogs {
		if e.Id == id {
			return e, i
		}
	}
	return Dog{}, -1
}

func posting(c *gin.Context) {
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"name\": \"Lassie\", \"owner\": \"Joe Carraclough\" }" /somePost
	var dog Dog
	c.Bind(&dog)

	if dog.Name != "" && dog.Owner != "" {
		dog.Id = nextId()
		Dogs = append(Dogs, dog)
		//return new dog so ID is available for client
		c.JSON(200, dog)

	} else {
		//serve if user imputs crap
		c.JSON(422, gin.H{"error": "name and owner must not be empty"})
	}
}

func deleting(c *gin.Context) {
	id := c.Param("id")

	dogId, err := strconv.Atoi(id)

	if err != nil {
		panic(err)
	}

	//find book
	_, index := getDogbyID(dogId)

	if index == -1 {
		c.JSON(422, gin.H{"error": "id does not exist"})
	} else {

		Dogs = append(Dogs[:index], Dogs[index+1:]...)
		c.JSON(200, "dog deleted")
	}
}

func putting(c *gin.Context) {

}
