package blockchain

import (
	"sync"
)

type Event uint

const (
	BlockSaved Event = iota
	HeaderSaved
)

func (e Event) String() string {
	switch e {
	case BlockSaved:
		return "BlockSaved"
	case HeaderSaved:
		return "HeaderSaved"
	}
	return ""
}

type Feed struct {
	lock sync.RWMutex
	subs []Subscription
}

func NewFeed() *Feed {
	return &Feed{}
}

func (f *Feed) Subscribe(channel chan<- Event) Subscription {
	sub := &subscriber{
		feed:    f,
		channel: channel,
	}

	f.lock.Lock()
	f.subs = append(f.subs, sub)
	f.lock.Unlock()

	return sub
}

func (f *Feed) Send(event Event) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	for _, sub := range f.subs {
		sub.Send(event)
	}
}

func (f *Feed) remove(sub Subscription) {
	f.lock.Lock()
	defer f.lock.Unlock()

	for index, s := range f.subs {
		if s == sub {
			f.subs = append(f.subs[:index], f.subs[index+1:]...)
		}
	}
}

type Subscription interface {
	Send(Event)
	Unsubscribe()
}

type subscriber struct {
	feed    *Feed
	channel chan<- Event
	once    sync.Once
}

func (sub *subscriber) Send(event Event) { sub.channel <- event }
func (sub *subscriber) Unsubscribe() {
	sub.once.Do(func() {
		sub.feed.remove(sub)
		close(sub.channel)
	})
}
