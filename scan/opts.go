package scan

// options 扫描选项
type options struct {
	onlyExported bool
	filter       Filter
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Option 选项
type Option func(o *options)

// WithOnlyExported 只导出大写开头的函数
func WithOnlyExported(onlyExported bool) Option {
	return func(o *options) {
		o.onlyExported = true
	}
}

// WithFilter 过滤器
func WithFilter(filter Filter) Option {
	return func(o *options) {
		o.filter = filter
	}
}
