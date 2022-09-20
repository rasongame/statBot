package main

import (
	"gorm.io/gorm"
	"time"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func LoadCache(db *gorm.DB) {
	var users []User
	db.Find(&users)
	for _, user := range users {
		CachedUsers[user.Id] = CacheUser{
			User:     user,
			LifeTime: time.Now().Unix() + 3600,
		}
	}
}

func UpdateCache(placeholder *SomePlaceholder, db *gorm.DB) {
	if CachedUsers[placeholder.User.ID].LifeTime <= time.Now().Unix() {
		u := CacheUser{
			User: User{
				Id:           placeholder.User.ID,
				FirstName:    placeholder.User.FirstName,
				LastName:     placeholder.User.LastName,
				Username:     placeholder.User.UserName,
				LanguageCode: placeholder.User.LanguageCode,
				LastSeen:     placeholder.LastSeenAt,
			},
			LifeTime: time.Now().Unix() + 3600, // 60 min
		}
		CachedUsers[placeholder.User.ID] = u

		db.Save(&u.User)

	}
}
