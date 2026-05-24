package expense

type ExpenseRepository interface {
	CreateExpenses(expense Expense)
}
