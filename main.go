package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Professor-Goo/gator/internal/config"
	"github.com/Professor-Goo/gator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login requires a username argument")
	}

	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user does not exist: %w", err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User has been set to: %s\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("register requires a username argument")
	}

	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("user already exists: %s", username)
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully:")
	fmt.Printf("  ID: %s\n", user.ID)
	fmt.Printf("  Name: %s\n", user.Name)
	fmt.Printf("  Created: %s\n", user.CreatedAt)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("couldn't reset database: %w", err)
	}

	fmt.Println("Database reset successfully")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)

	appState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	if len(os.Args) < 2 {
		fmt.Println("Error: not enough arguments provided")
		fmt.Println("Usage: gator <command> [args...]")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := []string{}
	if len(os.Args) > 2 {
		cmdArgs = os.Args[2:]
	}

	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	if err := cmds.run(appState, cmd); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
