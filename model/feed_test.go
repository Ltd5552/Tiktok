package model

import (
	"Tiktok/config"
	"fmt"
	"testing"
)

func TestJudgeFavorite(t *testing.T) {
	config.InitViper()
	InitDB()
	fmt.Println(JudgeFavorite(10, 1))
}
