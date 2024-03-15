package telegram

type UserEvent struct{}

type Question struct{}

type Output struct{}

type UseCaseResult struct{}

type UserCase interface {
	Next() bool
	Output() Output
}
