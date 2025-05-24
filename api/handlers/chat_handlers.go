package handlers

import (
	"api/api/models"

	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type GroupConnections map[uint]*websocket.Conn

var upgrader = websocket.Upgrader{}
var connections = map[uint]GroupConnections{}

func handleUserMessages(ws *websocket.Conn, userName string, userMail string, groupName string, connectedUsers GroupConnections) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			slog.Error("Error reading message", err.Error())
		}
		textMessage := string(message)
		slog.Info(fmt.Sprintf("Client %s sent a message to chat %s", userMail, groupName), "message", textMessage)

		for connection := range maps.Values(connectedUsers) {
			connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(
				`\{"username": "%s", "time": "%s", "message": "%s"\}`,
				userName,
				time.Now().String(),
				textMessage,
			)))
		}
	}
}

func StablishWSConnection(context *gin.Context) {
	userId := context.Keys["user_id"].(uint)
	groupId := context.Param("id")
	var group models.Group
	var user models.User

	if err := DB.Preload("Users").First(&group, groupId).Error; err != nil {
		slog.Warn(fmt.Sprintf("User %d not found", userId))
		models.ResponseJSON(context, http.StatusNotFound, "Group not found", nil)
		return
	}
	if err := DB.First(&user, userId).Error; err != nil {
		slog.Warn(fmt.Sprintf("Group %s not found", groupId))
		models.ResponseJSON(context, http.StatusNotFound, "User not found", nil)
		return
	}

	if !slices.ContainsFunc(group.Users, func(u models.User) bool {
		return u.ID == userId
	}) {
		slog.Warn(fmt.Sprintf("User %s does not belong to group %s", user.Mail, group.Name))
		models.ResponseJSON(context, http.StatusForbidden, "User does not belong to this group", nil)
		return
	}

	ws, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		slog.Error("Error upgrading connection", err.Error())
		models.ResponseJSON(context, http.StatusInternalServerError, "Unable to upgrade connection", nil)
		return
	}
	defer ws.Close()

	if _, ok := connections[group.ID]; !ok {
		connections[group.ID] = GroupConnections{}
	}
	connectedUsers := connections[group.ID]
	if conn, ok := connectedUsers[userId]; ok {
		conn.Close()
	}
	connectedUsers[userId] = ws

	slog.Info(fmt.Sprintf("Client %s connected to chat %s", user.Mail, group.Name))
	handleUserMessages(ws, user.Username, user.Mail, group.Name, connectedUsers)
}
