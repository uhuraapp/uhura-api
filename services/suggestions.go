package services

import (
  "github.com/jinzhu/gorm"
)

type SuggestionsService struct {
  DB *gorm.DB
}

func NewSuggestionsService(db *gorm.DB) SuggestionsService {
  return SuggestionsService{DB: db}
}
