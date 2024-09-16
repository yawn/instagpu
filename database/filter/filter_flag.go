package filter

import (
	"github.com/spf13/pflag"
	"github.com/yawn/spottty/detect"
)

type (
	Filter                func(p *detect.Prices) bool
	install[T comparable] func(*T, string, T, string)
)

type filterFlag[T comparable] struct {
	defaultValue T
	description  string
	filter       func(T) Filter
	ignore       T
	install      func(flags *pflag.FlagSet) install[T]
	name         string
	value        T
}

func (f *filterFlag[T]) Apply() Filter {
	return f.filter(f.value)
}

func (f *filterFlag[T]) Install(fs *pflag.FlagSet) {
	f.install(fs)(&f.value, f.name, f.defaultValue, f.description)
}

func (f *filterFlag[T]) IsSet() bool {
	return f.value != f.ignore
}

func (f *filterFlag[T]) Name() string {
	return f.name
}
