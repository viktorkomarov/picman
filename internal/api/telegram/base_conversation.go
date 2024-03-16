package telegram

type Action func(fsmContex UseCaseContext, event UserEvent) StateResult

type BaseConversation struct {
	question QuestionProvider
	action   Action
}

func NewBaseConversation(question QuestionProvider, action Action) BaseConversation {
	return BaseConversation{
		question: question,
		action:   action,
	}
}

func (b BaseConversation) Question() QuestionProvider {
	return b.question
}

func (b BaseConversation) ApplyEvent(fsmContext UseCaseContext, event UserEvent) StateResult {
	return b.action(fsmContext, event)
}
