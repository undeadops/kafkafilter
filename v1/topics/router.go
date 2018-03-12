package topics

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Define our message object
type Post struct {
	Offset int64  `json:"offset"`
	Msg    string `json:"msg"`
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TopicsRegister(router *gin.RouterGroup) {
	router.GET("/", TopicsList)
	router.GET("/:topic", func(c *gin.Context) {
		topic := c.Param("topic")
		query := c.Request.URL.Query()
		filter := query["filter"][0]
		TopicStream(c.Writer, c.Request, topic, filter)
	})
}

func TopicsList(c *gin.Context) {

	var topics []string

	topics = []string{"stage.throttle.activity.created",
		"prod.falken.enriched_activity.created",
	}
	c.JSON(http.StatusOK, gin.H{"topics": topics})
}

func TopicStream(w http.ResponseWriter, r *http.Request, topic string, filter string) {

	stream := make(chan Post)

	go KafkaStream(stream, topic, filter)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}
	defer conn.Close()

	for {
		conn.WriteJSON(<-stream)
	}
}

func KafkaStream(stream chan Post, topic string, filter string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Specify brokers address. This is default one
	brokers := []string{"localhost:9092"}

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	// How to decide partition, is it fixed value...?
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	for {
		var post Post
		select {
		case err := <-consumer.Errors():
			fmt.Println(err)
		case msg := <-consumer.Messages():
			if strings.Contains(string(msg.Value), filter) {
				//fmt.Println("Received messages", string(msg.Key), string(msg.Value))
				post = Post{Offset: msg.Offset, Msg: string(msg.Value)}
				//fmt.Printf("Text: %s\n\n", post.Text)
				stream <- post
			}
		}
		//fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}
