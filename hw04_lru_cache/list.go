package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len  int
	head *ListItem
	tail *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	if l.head != nil {
		return l.head
	}
	return nil
}

func (l *list) Back() *ListItem {
	if l.tail != nil {
		return l.tail
	}
	return nil
}

func (l *list) PushFront(v interface{}) *ListItem {
	newHeadItm := ListItem{Value: v, Next: nil, Prev: nil}

	if l.Len() == 0 {
		l.tail = &newHeadItm
	} else {
		currentFrontItm := l.Front()
		newHeadItm.Next = currentFrontItm
		currentFrontItm.Prev = &newHeadItm
	}

	l.len++
	l.head = &newHeadItm
	return &newHeadItm
}

func (l *list) PushBack(v interface{}) *ListItem {
	newTailItm := ListItem{Value: v, Next: nil, Prev: nil}

	if l.Len() == 0 {
		l.head = &newTailItm
	} else {
		currentBackItm := l.Back()
		newTailItm.Prev = currentBackItm
		currentBackItm.Next = &newTailItm
	}

	l.len++
	l.tail = &newTailItm
	return &newTailItm
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil && i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	if i.Prev != nil && i.Next == nil {
		i.Prev.Next = nil
	}
	if i.Next != nil && i.Prev == nil {
		i.Next.Prev = nil
	}
	i.Next = nil
	i.Prev = nil
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Len() <= 1 {
		return
	}
	if i.Next != nil && i.Prev == nil {
		return
	}
	if i.Prev != nil && i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	if i.Prev != nil && i.Next == nil {
		i.Prev.Next = nil
		l.tail = i.Prev
	}
	l.head.Prev = i
	i.Next = l.head
	i.Prev = nil
	l.head = i
}

func NewList() List {
	return new(list)
}
