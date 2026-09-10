package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/itlabil/smmusic/backend/internal/config"
	"github.com/itlabil/smmusic/backend/internal/repository"
	"github.com/itlabil/smmusic/backend/internal/service"
	"golang.org/x/term"
)

func main() {
	cfg := config.Load()
	db := config.NewDatabase(cfg)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Admin username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Admin password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("failed to read password: %v", err)
	}
	password := strings.TrimSpace(string(passwordBytes))
	fmt.Println()

	user, err := authService.CreateUser(username, password, "admin")
	if err != nil {
		log.Fatalf("failed to create admin user: %v", err)
	}

	fmt.Printf("Admin user created successfully: id=%d, username=%s\n", user.ID, user.Username)
}