package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"
)

type Message struct {
	Event   string `json:"event"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

const (
	KEY  = "chat_id"
	WAIT = "wait"
)

func NewMessage(event, name, content string) *Message {
	return &Message{
		Event:   event,
		Name:    name,
		Content: content,
	}
}

func (m *Message) GetByteMessage() []byte {
	result, _ := json.Marshal(m)
	return result
}

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "a12345",
		DB:       0, //use default DB
	})
	pong, err := redisClient.Ping().Result()
	if err == nil {
		log.Println("redis 回应成功：", pong)
	} else {
		log.Fatal("redis 无法连接，错误为：", err)
	}
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("template/html/*")
	r.Static("/assets", "./template/assets")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	m := melody.New()
	r.GET("/ws", func(ctx *gin.Context) {
		m.HandleRequest(ctx.Writer, ctx.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		id := GetSessionID(s)
		chatTo, _ := redisClient.Get(id).Result()
		m.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(KEY)
			return compareID == chatTo || compareID == id
		})
	})

	m.HandleConnect(func(session *melody.Session) {
		id := InitSession(session)
		if key, err := GetWaitFirstKey(); err == nil && key != "" {
			CreateChat(id, key)
			msg := NewMessage("other", "对方已经", "加入聊天室").GetByteMessage()
			m.BroadcastFilter(msg, func(session *melody.Session) bool {
				compareID, _ := session.Get(KEY)
				return compareID == id || compareID == key
			})
		} else {
			AddToWaitList(id)
		}
	})

	m.HandleClose(func(session *melody.Session, i int, s string) error {
		id := GetSessionID(session)
		chatTo, _ := redisClient.Get(id).Result()
		msg := NewMessage("other", "对方已经", "离开聊天室").GetByteMessage()
		RemoveChat(id, chatTo)
		return m.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(KEY)
			return compareID == chatTo
		})
	})

	r.Run(":5001")
}

func AddToWaitList(id string) error {
	return redisClient.LPush(WAIT, id).Err()
}

func GetWaitFirstKey() (string, error) {
	return redisClient.LPop(WAIT).Result()
}

func CreateChat(id1, id2 string) {
	redisClient.Set(id1, id2, 0)
	redisClient.Set(id2, id1, 0)
}

func RemoveChat(id1, id2 string) {
	redisClient.Del(id1, id2)
}

func GetSessionID(s *melody.Session) string {
	if id, isExist := s.Get(KEY); isExist {
		return id.(string)
	}
	return InitSession(s)
}

func InitSession(s *melody.Session) string {
	id := uuid.New().String()
	s.Set(KEY, id)
	return id
}
