package ws

import (
	"encoding/json"
	"log"
	"sync"

	"eikva.ru/eikva/models"
	"eikva.ru/eikva/session"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var instance = websocket.Upgrader{}

type WSMessageType string

const (
	WSMessageTypeTestCaseUpdate WSMessageType = "test-case-update"
	WSMessageTypeAuth           WSMessageType = "auth"
)

type WSUpdateNotification struct {
	Type WSMessageType `json:"type"`
	UUID []string        `json:"uuid"`
}

type WSAuthMessage struct {
	Type        WSMessageType `json:"type"`
	AccessToken string        `json:"access_token"`
}

type ConnectionInfo struct {
	User *models.User
}

type ConntectionManager struct {
	mu           sync.Mutex
	conntections map[*websocket.Conn]ConnectionInfo
}

func (cm *ConntectionManager) Add(c *websocket.Conn, u *models.User) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.conntections[c] = ConnectionInfo{User: u}
}

func (cm *ConntectionManager) Remove(c *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.conntections, c)
}

func (cm *ConntectionManager) BroadCastTestCaseUpdate(uuid ...string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for c := range cm.conntections {
		c.WriteJSON(&WSUpdateNotification{
			Type: WSMessageTypeTestCaseUpdate,
			UUID: uuid,
		})
	}
}

func NewConnectionManager() *ConntectionManager{
	return &ConntectionManager{
		conntections: make(map[*websocket.Conn]ConnectionInfo),
	}
}

var WSConntections = NewConnectionManager()

func HandleSubscribers(ctx *gin.Context) {
	c, err := instance.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	defer func() {
		WSConntections.Remove(c)
		c.Close()
	}()

	_, msg, err := c.ReadMessage()
	if err != nil {
		log.Println("Failed to read auth message:", err)
		return
	}

	var authMessage WSAuthMessage
	if err := json.Unmarshal(msg, &authMessage); err != nil || authMessage.Type != WSMessageTypeAuth {
		c.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(
				websocket.ClosePolicyViolation,
				"Не верный тип сообщения авторизации",
			),
		)
		return
	}

	user, err := session.ValidateSessionTokenAndGetUser(
		authMessage.AccessToken,
		session.TokenTypeAccess,
	)

	if err != nil {
		c.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(
				websocket.ClosePolicyViolation,
				err.Error(),
			),
		)
		return
	}

	WSConntections.Add(c, user)

	for {
        msgType, _, err := c.ReadMessage()
		if msgType == websocket.CloseMessage {
            log.Println("Client closed connection:")
			break
		}

        if err != nil {
            log.Println("Client disconnected:", err)
            break
        }
    }
}
