package cmdrouter

import (
	"context"
	"fmt"
	"log/slog"
)

type TablePrinter interface {
	PrintTable(headers []string, rows [][]any)
}

type CmdRouter struct {
	name         string
	handlers     []OptionHandler
	tablePrinter TablePrinter
}

type OptionHandler struct {
	Name string                          // name of the operation, e.g. "login"
	Run  func(ctx context.Context) error // function to run the operation
}

func NewCmdRouter(name string, tablePrinter TablePrinter, handlers ...OptionHandler) *CmdRouter {
	return &CmdRouter{
		name:         name,
		handlers:     handlers,
		tablePrinter: tablePrinter,
	}
}

func (c *CmdRouter) SetOptionHandler(name string, handler func(ctx context.Context) error) {
	c.handlers = append(c.handlers, OptionHandler{Name: name, Run: handler})
}

func (c *CmdRouter) Run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic", "err", r)
		}
	}()

	const exitNumber = 0

	for {
		option := c.getOption()
		if option == exitNumber {
			break
		}

		if err := c.handlers[option-1].Run(ctx); err != nil {
			slog.Error("handler", "err", err)
			continue
		}

		fmt.Println()
	}
}

func (c CmdRouter) getOption() int {
	c.showMenu()

	var option int
	for {
		fmt.Print("Enter option number: ")
		if _, err := fmt.Scan(&option); err == nil &&
			option >= 0 && option <= len(c.handlers) {
			break
		}

		fmt.Println("Invalid number. Try again.")
	}

	return option
}

func (c *CmdRouter) showMenu() {
	headers := []string{"#", c.name}
	rows := make([][]any, 0, len(c.handlers))

	for i := range c.handlers {
		rows = append(rows, []any{i + 1, c.handlers[i].Name})
	}
	rows = append(rows, []any{0, "Exit"})

	c.tablePrinter.PrintTable(headers, rows)
	fmt.Println()
}
