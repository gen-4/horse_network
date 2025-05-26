package broker

import (
	"slices"
	"sync"
	"time"
)

type StoredMessage struct {
	Type    string
	Message string
	Time    time.Time
}

type Channel struct {
	Chan chan StoredMessage
	Num  uint
}

type Broker struct {
	Channels map[uint][]chan StoredMessage
	Subs     map[uint]Channel
	Mu       sync.Mutex
}

func (broker *Broker) Subscribe(subId uint, targetId uint) chan StoredMessage {
	var channel Channel
	if _, ok := broker.Subs[subId]; !ok {
		channel = Channel{
			Chan: make(chan StoredMessage),
			Num:  0,
		}
		broker.Subs[subId] = channel
	}
	channel = broker.Subs[subId]

	broker.Mu.Lock()
	if _, ok := broker.Channels[targetId]; !ok {
		broker.Channels[targetId] = []chan StoredMessage{}
	}
	broker.Channels[targetId] = append(broker.Channels[targetId], channel.Chan)
	broker.Mu.Unlock()
	channel.Num += 1
	return channel.Chan
}

func (broker *Broker) UnSubscribe(subId uint, targetId uint) {
	broker.Mu.Lock()
	defer broker.Mu.Unlock()
	if _, ok := broker.Channels[targetId]; !ok {
		return
	}

	channel := broker.Subs[subId]
	groupChannels := broker.Channels[targetId]
	index := slices.Index(groupChannels, channel.Chan)
	broker.Channels[targetId] = slices.Delete(groupChannels, index, index+1)
	channel.Num -= 1
	if channel.Num >= 1 {
		return
	}

	close(channel.Chan)
	delete(broker.Subs, subId)
}
