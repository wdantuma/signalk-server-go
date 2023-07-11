package filter

import (
	"log"
	"regexp"

	"github.com/wdantuma/signalk-server-go/signalk"
)

type Subscribe int

const (
	Self Subscribe = iota + 1
	All
	None
)

func ParseSubscribe(subscribe string) Subscribe {
	switch subscribe {
	case "all":
		return All
	case "none":
		return None
	default:
		return Self
	}
}

type Subscription struct {
	Context *regexp.Regexp
	Paths   []*regexp.Regexp
}

type Filter struct {
	Subscribe     Subscribe
	Self          string
	Subscriptions []*Subscription
}

func NewFilter(self string) *Filter {
	return &Filter{Subscribe: Self, Self: self, Subscriptions: make([]*Subscription, 0)}
}

func (f *Filter) Filter(input <-chan signalk.DeltaJson) <-chan signalk.DeltaJson {
	output := make(chan signalk.DeltaJson)
	go func() {
		for delta := range input {
			if delta.Context == nil {
				continue
			}
			include := false
			if delta.Context != nil && f.Subscribe == Self && *delta.Context == f.Self {
				include = true
			}
			if f.Subscribe == All {
				include = true
			}
			if delta.Context == nil {
				log.Println("ERROR")
			}
			if !include {
				for _, s := range f.Subscriptions {
					if s.Context.Match([]byte(*delta.Context)) {
						include = true
						break
					}
				}
			}

			if include {
				output <- delta
			}
		}
		close(output)
	}()

	return output
}

func (f *Filter) UpdateSubscription(subscribe signalk.SubscribeJson) {
	subscription := &Subscription{Paths: make([]*regexp.Regexp, 0)}
	cr, err := regexp.Compile(subscribe.Context)
	if err == nil {
		subscription.Context = cr
	} else {
		return
	}

	for _, s := range subscribe.Subscribe {
		path := *s.Path
		if path == "*" {
			path = ".*"
		}
		pr, err := regexp.Compile(path)
		if err == nil {
			subscription.Paths = append(subscription.Paths, pr)
		}
	}
	f.Subscriptions = append(f.Subscriptions, subscription)
}
