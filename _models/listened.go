package models

import "time"

type Listened struct {
	Id        int64
	UserId    int64
	ItemId    int64
	ItemUID   string `sql:"item_uid"`
	Viewed    bool
	ChannelId int64
	StoppedAt int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (l Listened) TableName() string {
	return "user_items"
}
func GetPlays(userID string, ids []int) {
	var plays []*models.Listened

	if len(ids) > 0 {
		s.DB.Table(models.Listened{}.TableName()).
			Where("item_uid IN (?)", ids).
			Where("user_id = ?", userID).
			Find(&plays)
	}
}
