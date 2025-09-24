package stat

import (
	"errors"
	"fmt"
	"go-advanced/pkg/db"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (repo *StatRepository) AddClick(linkId uint) {
	if repo == nil || repo.Db == nil {
		panic("StatRepository or Db is nil")
	}

	currDate := datatypes.Date(time.Now())
	var stat Stat

	// Ищем существующую запись
	result := repo.Db.Where("link_id = ? AND date = ?", linkId, currDate).First(&stat)

	if result.Error != nil {
		// Если записи нет - создаем новую
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newStat := Stat{
				LinkId: linkId,
				Clicks: 1,
				Date:   currDate,
			}
			createResult := repo.Db.Create(&newStat)
			if createResult.Error != nil {
				// Логируем ошибку, но не паникуем
				fmt.Printf("Error creating stat: %v\n", createResult.Error)
			}
		} else {
			// Другие ошибки БД
			fmt.Printf("Error finding stat: %v\n", result.Error)
		}
	} else {
		// Запись найдена - обновляем
		stat.Clicks += 1
		updateResult := repo.Db.Save(&stat)
		if updateResult.Error != nil {
			fmt.Printf("Error updating stat: %v\n", updateResult.Error)
		}
	}
}

func (repo *StatRepository) GetStats(by string, from, to time.Time) []GetStatResponse {
	var stats []GetStatResponse
	var selectQuery string

	switch by {
	case GROUP_BY_DAY:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case GROUP_BY_MONTH:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}

	repo.Db.Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)

	return stats
}
