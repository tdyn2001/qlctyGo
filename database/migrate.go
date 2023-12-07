package main

import (
	"fmt"
	"v2/initializers"
	"v2/models"
)

func init() {
	initializers.ConnectDB(initializers.GetConfig())
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
	fmt.Println("? Migration complete")
}
