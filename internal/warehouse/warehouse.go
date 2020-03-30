package warehouse

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const invitationTimeout = 500 * time.Millisecond

func StartWarehouse(port string) {
	if len(os.Args) != 2 {
		log.Fatal("Please, provide port to listen on as an argument")
	}
	addr := "127.0.0.1:" + port
	go sendInvitations(addr)
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		log.Println("Stock center contaced us!")
		c.JSON(http.StatusNoContent, gin.H{})
	})
	r.POST("/take", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Error reading request data -- take items")
		}
		var items []int
		err = json.Unmarshal(body, &items)
		if err != nil {
			log.Println("Error parsing server request -- take items " + err.Error())
			c.JSON(400, "Invalid request")
			return
		}
		log.Printf("taking items: %v\n", items)

		c.JSON(http.StatusOK, gin.H{
			"status": "fine",
			"items":  items,
		})
	})
	r.GET("/getItems", func(c *gin.Context) {
		c.JSON(http.StatusOK, []int{1, 2, 3})
	})
	r.Run(":" + port)
}

// continuously broadcast invitation message over UDP
// with address to connect
func sendInvitations(myAddr string) {
	for {
		time.Sleep(invitationTimeout)
		con, _ := net.Dial("udp", "127.0.0.1:3000")
		buf := []byte(myAddr)
		_, err := con.Write(buf)
		if err != nil {
			log.Println(err)
		}
	}
}