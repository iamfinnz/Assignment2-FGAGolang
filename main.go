package main

import (
	"app/database"
	"app/routes"
	"fmt"
)

func main() {
	PORT := ":8080"
	fmt.Println("Server berjalan pada port", PORT)
	defer database.CloseConnection()
	routes.StartServer().Run(PORT)
}
