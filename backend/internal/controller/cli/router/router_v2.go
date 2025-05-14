package cmdrouter

import (
	"context"
	"fmt"
	"log/slog"
)

type CmdRouterV2 struct {
	name         string
	handlers     []OptionHandler
	tablePrinter TablePrinter
	isGroup      bool
}

func NewCmdRouterV2(name string, tablePrinter TablePrinter, handlers ...OptionHandler) *CmdRouterV2 {
	return &CmdRouterV2{
		name:         name,
		handlers:     handlers,
		tablePrinter: tablePrinter,
		isGroup:      false,
	}
}

func (c *CmdRouterV2) Group(name string, handlers ...OptionHandler) *CmdRouterV2 {
	group := &CmdRouterV2{
		name:         name,
		handlers:     handlers,
		tablePrinter: c.tablePrinter,
		isGroup:      true,
	}

	c.SetOptionHandlers(OptionHandler{Name: name, Run: func(ctx context.Context) error {
		group.Run(ctx)
		return nil
	}})

	return group
}

func (c *CmdRouterV2) SetOptionHandlers(handlers ...OptionHandler) {
	c.handlers = append(c.handlers, handlers...)
}

func (c *CmdRouterV2) Run(ctx context.Context) {
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

func (c CmdRouterV2) getOption() int {
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

func (c *CmdRouterV2) showMenu() {
	headers := []string{"#", c.name}
	rows := make([][]any, 0, len(c.handlers))

	for i := range c.handlers {
		rows = append(rows, []any{i + 1, c.handlers[i].Name})
	}

	if c.isGroup {
		rows = append(rows, []any{0, "<-Back"})
	} else {
		rows = append(rows, []any{0, "Exit"})
	}

	c.tablePrinter.PrintTable(headers, rows)
	fmt.Println()
}
