package core

import "strconv"

func (ce *ChannelEntity) FixUri() string {
	ch := Channel{Title: ce.Title, Id: ce.Id}
	ce.Uri = ch.SetUri()
	return ce.Uri
}

func (ce *ChannelEntity) SetSubscribed(userId string) {
	status := true
	go func() {
		channelId := strconv.Itoa(int(ce.Id))
		cacheKey := "s:" + channelId + ":" + userId
		CacheSet(cacheKey, status)
	}()
	ce.Subscribed = status
}

func (ce *ChannelEntity) SetEpisodesIds() {
	ce.Episodes = getEpisodesIds(ce.Id)
}

func (ce *ChannelEntity) SetSubscription(userId string) {
	var status bool
	key := "s:" + strconv.Itoa(int(ce.Id)) + ":" + userId
	_, err := CacheGet(key, status)
	ce.Subscribed = err == nil
}
