package user_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	cmdrouter "github.com/hahaclassic/orpheon/backend/internal/controller/cli/router"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/user"
)

type UserController struct {
	userService user.UserService
}

func NewUserController(userService user.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "My profile",
			Run:  c.getMyProfile,
		},
		{
			Name: "Update my profile",
			Run:  c.updateUser,
		},
		{
			Name: "Get user by ID",
			Run:  c.getUserByID,
		},
		{
			Name: "Delete user",
			Run:  c.deleteUser,
		},
	}
}

func (c *UserController) getMyProfile(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to get your profile.")
		return nil
	}

	user, err := c.userService.GetUser(ctx, session.Claims().UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	output.PrintUser(user)
	return nil
}

func (c *UserController) getUserByID(ctx context.Context) error {
	var id string

	fmt.Print("Enter user ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read user ID: %w", err)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %w", err)
	}

	user, err := c.userService.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	output.PrintUser(user)
	return nil
}

func (c *UserController) updateUser(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to update your profile.")
		return nil
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter new user name: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Enter new birth date (YYYY-MM-DD): ")
	scanner.Scan()
	birthDateStr := scanner.Text()
	birthDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return fmt.Errorf("failed to parse birth date: %w", err)
	}

	user := &entity.UserInfo{
		ID:        session.Claims().UserID,
		Name:      name,
		BirthDate: birthDate,
	}

	err = c.userService.UpdateUser(ctx, session.Claims(), user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	fmt.Println("User updated successfully")
	return nil
}

func (c *UserController) deleteUser(ctx context.Context) error {
	var id string

	fmt.Print("Enter user ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read user ID: %w", err)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %w", err)
	}

	err = c.userService.DeleteUser(ctx, session.Claims(), userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	fmt.Println("User deleted successfully")
	return nil
}
