package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	size int
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v any) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  l.head,
		Prev:  nil,
	}

	if l.head != nil {
		l.head.Prev = newItem
	} else {
		l.tail = newItem
	}

	l.head = newItem
	l.size++

	return newItem
}

func (l *list) PushBack(v any) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.tail,
	}

	if l.tail != nil {
		l.tail.Next = newItem
	} else {
		l.head = newItem
	}

	l.tail = newItem
	l.size++

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	i.Prev = nil
	i.Next = nil
	l.size--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head == i {
		return
	}

	l.Remove(i)
	l.size++

	l.head.Prev = i
	i.Next = l.head
	l.head = i
}

func NewList() List {
	return new(list)
}
