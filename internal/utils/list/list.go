package list

type Head[A any] *node[A]

type node[A any] struct {
	val  A
	next *node[A]
}

func NewHead[A any](val A) *node[A] {
	return &node[A]{
		val: val,
	}
}

func IsEmpty[A any](head Head[A]) bool {
	return head == nil
}

func HeadTail[A any](head Head[A]) (Head[A], Head[A]) {
	next := head
	head.next = nil
	return head, next
}

func AddNewHead[A any](val A, head Head[A]) Head[A] {
	newHead := NewHead(val)
	newHead.next = head
	return newHead
}
