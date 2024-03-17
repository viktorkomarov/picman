package usecases

type QuestionProvider struct {
	skip     bool
	question string
}

func SkipQuestionProvider() QuestionProvider {
	return QuestionProvider{skip: true}
}

func NewQuestionProvider(question string) QuestionProvider {
	return QuestionProvider{question: question}
}

func (q QuestionProvider) ShouldAsk() bool {
	return !q.skip
}
