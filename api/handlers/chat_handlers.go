package handlers

import (
	"api/api/models"

	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const MAX_BATCH_SIZE = 100

type GroupConnections map[uint]*websocket.Conn

var upgrader = websocket.Upgrader{}
var connections = map[uint]GroupConnections{}

func handleUserMessages(ws *websocket.Conn, user models.User, group models.Group, connectedUsers GroupConnections) {
	for {
		var messageEntity models.Message
		_, message, err := ws.ReadMessage()
		if err != nil {
			slog.Error("Error reading message", err.Error())
		}
		textMessage := string(message)
		messageEntity = models.Message{
			Message: textMessage,
			User:    user,
			Group:   group,
		}
		if err := DB.Create(&messageEntity).Error; err != nil {
			slog.Error(fmt.Sprintf("Unable to save message %s from user %s in chat %s", textMessage, user.Mail, group.Name), "error", err.Error)
			ws.WriteJSON(models.JsonResponse{
				Status:  http.StatusInternalServerError,
				Message: "Unable to save message",
				Data:    models.MessageDto{},
			})

		}
		slog.Info(fmt.Sprintf("Client %s sent a message to chat %s", user.Mail, group.Name), "message", textMessage)

		for connection := range maps.Values(connectedUsers) {
			connection.WriteJSON(models.JsonResponse{
				Status:  http.StatusOK,
				Message: "Message received",
				Data: models.MessageDto{
					Username: user.Username,
					Time:     time.Now().String(),
					Message:  textMessage,
				},
			})
		}
	}
}

func userBelongsToGroup(context *gin.Context, user models.User, group models.Group) bool {
	if !slices.ContainsFunc(group.Users, func(u models.User) bool {
		return u.ID == user.ID
	}) {
		slog.Warn(fmt.Sprintf("User %s does not belong to group %s", user.Mail, group.Name))
		models.ResponseJSON(context, http.StatusForbidden, "User does not belong to this group", nil)
		return false
	}

	return true
}

func getUserAndGroup(context *gin.Context) (models.User, models.Group, error) {
	userId := context.Keys["user_id"].(uint)
	groupId := context.Param("id")
	var group models.Group
	var user models.User

	if err := DB.Preload("Users").First(&group, groupId).Error; err != nil {
		slog.Warn(fmt.Sprintf("User %d not found", userId))
		models.ResponseJSON(context, http.StatusNotFound, "Group not found", nil)
		return models.User{}, models.Group{}, err
	}
	if err := DB.First(&user, userId).Error; err != nil {
		slog.Warn(fmt.Sprintf("Group %s not found", groupId))
		models.ResponseJSON(context, http.StatusNotFound, "User not found", nil)
		return models.User{}, models.Group{}, err
	}

	return user, group, nil
}

func StablishWSConnection(context *gin.Context) {
	user, group, err := getUserAndGroup(context)
	if err != nil {
		return
	}
	if isInGroup := userBelongsToGroup(context, user, group); !isInGroup {
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
	if conn, ok := connectedUsers[user.ID]; ok {
		conn.Close()
	}
	connectedUsers[user.ID] = ws

	slog.Info(fmt.Sprintf("Client %s connected to chat %s", user.Mail, group.Name))
	handleUserMessages(ws, user, group, connectedUsers)
}

func GetHistory(context *gin.Context) {
	var messages []models.Message
	user, group, err := getUserAndGroup(context)
	if err != nil {
		return
	}
	if isInGroup := userBelongsToGroup(context, user, group); !isInGroup {
		return
	}

	underIndex, err := strconv.ParseUint(context.DefaultQuery("index", "0"), 10, 0)
	if err != nil {
		models.ResponseJSON(context, http.StatusBadRequest, "No valid index provided", nil)
		return
	}
	batchSize, err := strconv.ParseUint(context.DefaultQuery("size", "20"), 10, 0)
	if err != nil || batchSize > MAX_BATCH_SIZE {
		models.ResponseJSON(context, http.StatusBadRequest, "No valid size provided", nil)
		return
	}

	err = DB.Preload("User").
		Order("messages.created_at ASC").
		Where("group_id = ?", group.ID).
		Limit(int(batchSize)).
		Offset(int(underIndex)).
		Find(&messages).Error
	if err != nil {
		models.ResponseJSON(context, http.StatusInternalServerError, "Unable to retrieve any message", nil)
		return
	}

	models.ResponseJSON(context, http.StatusOK, "Message batch retrieved", models.ToMessageDtos(messages))
}
