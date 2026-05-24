package expense

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrEmptyTitle = errors.New("title cannot be empty")
	ErrZeroAmount = errors.New("amount cannot be zero")
)

type ExpenseID int64
type UserID int64
type CategoryID int64

type ExpenseTitle string

func NewExpenseTitle(raw string) (ExpenseTitle, error) {
	title := strings.TrimSpace(raw)

	if title == "" {
		return "", ErrEmptyTitle
	}

	return ExpenseTitle(title), nil
}

type CategoryName string

func NewCategoryName(raw string) (CategoryName, error) {
	name := strings.TrimSpace(raw)

	if name == "" {
		return "", ErrEmptyTitle
	}

	return CategoryName(name), nil
}

type ExpenseAmount int64

func NewExpenseAmount(v int64) (ExpenseAmount, error) {
	if v == 0 {
		return 0, ErrZeroAmount
	}

	return ExpenseAmount(v), nil
}

type Expense struct {
	id         ExpenseID
	title      ExpenseTitle
	amount     ExpenseAmount
	date       time.Time
	categoryID CategoryID
	userID     UserID
}

func NewExpense(
	id ExpenseID,
	title ExpenseTitle,
	amount ExpenseAmount,
	date time.Time,
	categoryID CategoryID,
	userID UserID,
) (*Expense, error) {

	if date.IsZero() {
		date = time.Now()
	}

	return &Expense{
		id:         id,
		title:      title,
		amount:     amount,
		date:       date,
		categoryID: categoryID,
		userID:     userID,
	}, nil
}

type Category struct {
	id     CategoryID
	name   CategoryName
	userID UserID
}

func NewCategory(
	id CategoryID,
	name CategoryName,
	userID UserID,
) *Category {

	return &Category{
		id:     id,
		name:   name,
		userID: userID,
	}
}
