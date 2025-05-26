package handlers

import (
	"api/api/broker"
	"api/api/models"

	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const MAX_BATCH_SIZE = 100
const MESSAGE_TYPE_TEMPLATE = "chat_message#group:__GROUP_ID__"

type IncommingMessage struct {
	TargetId uint   `json:"target_id"`
	Message  string `json:"message"`
}

var chatChannels = broker.Broker{
	Channels: map[uint][]chan broker.StoredMessage{},
	Subs:     map[uint]broker.Channel{},
}

var upgrader = websocket.Upgrader{}

func handleUserMessages(ws *websocket.Conn, user models.User) {
	slog.Info(fmt.Sprintf("Writing incoming messages from user %s to storage", user.Mail))
	for {
		var messageEntity models.Message
		var sentMessage IncommingMessage
		var group models.Group
		err := ws.ReadJSON(&sentMessage)
		if err != nil {
			slog.Error("Error reading message", "error", err.Error())
			ws.WriteJSON(models.JsonResponse{
				Status:  http.StatusInternalServerError,
				Message: "Unable to read message",
				Data:    nil,
			})
			continue
		}
		if err := DB.First(&group, sentMessage.TargetId).Error; err != nil {
			slog.Warn(fmt.Sprintf("Group %d not found", sentMessage.TargetId))
			ws.WriteJSON(models.JsonResponse{
				Status:  http.StatusNotFound,
				Message: "Group not found",
				Data:    nil,
			})
			continue
		}

		messageEntity = models.Message{
			Message: sentMessage.Message,
			User:    user,
			Group:   group,
		}
		if err := DB.Create(&messageEntity).Error; err != nil {
			slog.Error(fmt.Sprintf(
				"Unable to save message %s from user %s in chat %s",
				sentMessage.Message,
				user.Mail,
				group.Name,
			), "error", err.Error)
			ws.WriteJSON(models.JsonResponse{
				Status:  http.StatusInternalServerError,
				Message: "Unable to save message",
				Data:    nil,
			})
			continue
		}

		storedMessage := broker.StoredMessage{
			Type:    strings.Replace(MESSAGE_TYPE_TEMPLATE, "__GROUP_ID__", strconv.FormatUint(uint64(group.ID), 10), 1),
			Message: sentMessage.Message,
			Time:    time.Now(),
		}
		chatChannels.Mu.Lock()
		if _, ok := chatChannels.Channels[group.ID]; !ok {
			chatChannels.Channels[group.ID] = []chan broker.StoredMessage{}
		}
		groupChannels := chatChannels.Channels[group.ID]
		for _, channel := range groupChannels {
			channel <- storedMessage
		}
		chatChannels.Mu.Unlock()

		slog.Info(fmt.Sprintf("Client %s sent a message to chat %s", user.Mail, group.Name), "message", sentMessage.Message)

		ws.WriteJSON(models.JsonResponse{
			Status:  http.StatusOK,
			Message: "Message received",
			Data:    nil,
		})
	}
}

func handleStoredMessages(wg *sync.WaitGroup, ws *websocket.Conn, user models.User) {
	var joinedGroups []uint
	var publishChan chan broker.StoredMessage
	defer wg.Done()
	slog.Info(fmt.Sprintf("Reading user %s stored messages", user.Mail))
	if err := DB.Table("group_users").
		Select("group_id").
		Where("user_id = ?", user.ID).
		Scan(&joinedGroups).Error; err != nil {
		slog.Error(fmt.Sprintf("Error fetching user %s groups", user.Mail), "error", err.Error())
	}
	for _, group := range joinedGroups {
		publishChan = chatChannels.Subscribe(user.ID, group)
	}

	for {
		var storedMessage broker.StoredMessage
		msg, ok := <-publishChan
		if ok {
			storedMessage = msg
		} else {
			continue
		}

		groupId, err := strconv.ParseUint(strings.Split(storedMessage.Type, "group:")[1], 10, 0)
		if err != nil {
			slog.Error(fmt.Sprintf("Error parsing uint from %s string", storedMessage.Type), "error", err.Error())
			continue
		}
		err = ws.WriteJSON(models.JsonResponse{
			Status:  http.StatusOK,
			Message: "New message",
			Data: models.MessageDto{
				Username: user.Username,
				Time:     time.Now().String(),
				Message:  storedMessage.Message,
				Group:    uint(groupId),
			},
		})

		if err != nil {
			slog.Info(fmt.Sprintf("Finished Reading user %s stored messages", user.Mail))
			break
		}
	}

	for _, group := range joinedGroups {
		chatChannels.UnSubscribe(user.ID, group)
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
	userId := context.Keys["user_id"].(uint)
	var user models.User
	var wg sync.WaitGroup

	if err := DB.First(&user, userId).Error; err != nil {
		slog.Warn(fmt.Sprintf("User %s not found", userId))
		models.ResponseJSON(context, http.StatusNotFound, "User not found", nil)
		return
	}

	ws, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		slog.Error("Error upgrading connection", err.Error())
		models.ResponseJSON(context, http.StatusInternalServerError, "Unable to upgrade connection", nil)
		return
	}
	defer ws.Close()

	slog.Info(fmt.Sprintf("Client %s connected", user.Mail))
	wg.Add(1)
	defer wg.Wait()
	go handleStoredMessages(&wg, ws, user)
	handleUserMessages(ws, user)
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
