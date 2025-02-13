package constants

import (
	"github.com/wafi04/chatting-app/config/env"
	"github.com/wafi04/chatting-app/services/shared/utils"
)

var (
	CommentID = utils.GenerateCustomID(utils.IDOptions{NumberLength: 7})
	PostID    = utils.GenerateCustomID(utils.IDOptions{CustomFormat: "POST", NumberLength: 7})
	DB_URL    = env.LoadEnv("DB_URL")
	PORT      = env.LoadEnv("PORT")
)
