package telegram

type ApplyResponse[Response any] func(response Response)

type PrepareСommand []interface {
	HasNext() bool
	Question() string
}

type UseCase interface {
}
