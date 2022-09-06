package verify

import "time"

type Options struct {
	Often time.Duration
}

func (o *Options) init() {
	if o.Often < 10 {
		o.Often = 10
	}
}
