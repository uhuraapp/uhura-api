package channels

import (
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/models"
)

func Ping(database *gorm.DB, id int64) {
	var channel models.Channel

	ok := database.Table(models.Channel{}.TableName()).
		First(&channel, id).Error != gorm.ErrRecordNotFound

	if !ok {
		return
	}
}
