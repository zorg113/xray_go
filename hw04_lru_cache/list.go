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
	Key   interface{}
}

type list struct {
	len     int
	elFront *ListItem
	elBack  *ListItem
	// Place your code here.
}

func NewList() List {
	return &list{0, nil, nil}
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.elFront
}

func (l list) Back() *ListItem {
	return l.elBack
}

func (l *list) PushFront(v interface{}) *ListItem {
	elNew := ListItem{Value: v, Prev: nil, Next: l.elFront}
	if l.len == 0 {
		l.elBack = &elNew
	} else {
		l.elFront.Prev = &elNew
	}
	l.len += 1
	l.elFront = &elNew
	return &elNew
}

func (l *list) PushBack(v interface{}) *ListItem {
	elNew := ListItem{Value: v, Prev: l.elBack, Next: nil}
	if l.len == 0 {
		l.elFront = &elNew
		l.elBack = &elNew
		l.len += 1 
		return &elNew
	}
	l.len += 1
	l.elBack.Next = &elNew
	l.elBack = &elNew
	return &elNew
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	prev := i.Prev
	next := i.Next
	if prev == nil {
		l.elFront = next
	} else {
		prev.Next = next
	}
	if next == nil {
		l.elBack = prev
	} else {
		next.Prev = prev
	}
	i.Next = nil
	i.Prev = nil
	l.len -= 1
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}
	prev := i.Prev
	next := i.Next
	if next != nil {
		next.Prev = prev
	}
	prev.Next = next
	l.elBack = prev
	i.Prev = nil
	i.Next = l.elFront
	l.elFront = i
}
