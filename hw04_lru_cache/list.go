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
	elFront *ListItem // указатель на первый элемент
	elBack  *ListItem // указатель на последний элемент
}

// создание списка
func NewList() List {
	return &list{0, nil, nil}
}

// длина списка
func (l list) Len() int {
	return l.len
}

// первый элемент списка
func (l list) Front() *ListItem {
	return l.elFront
}

// последний элемент списка
func (l list) Back() *ListItem {
	return l.elBack
}

// добавлеие элемннта в начало списка
func (l *list) PushFront(v interface{}) *ListItem {
	elNew := ListItem{Value: v, Prev: nil, Next: l.elFront}
	if l.len == 0 {
		l.elBack = &elNew
	} else {
		l.elFront.Prev = &elNew
	}
	l.len++
	l.elFront = &elNew
	return &elNew
}

// добавление элмента в конец списка
func (l *list) PushBack(v interface{}) *ListItem {
	elNew := ListItem{Value: v, Prev: l.elBack, Next: nil}
	if l.len == 0 {
		l.elFront = &elNew
		l.elBack = &elNew
	} else {
		l.elBack.Next = &elNew
		l.elBack = &elNew
	}
	l.len++
	return &elNew
}

// удаление элемента из списка
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
	l.len--
}

// пермещение элемнта в начало списка
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
