package io

const (
	optionType    = "type"
	optionProject = "project"
)

// Option is a workaround to normalize BigQuery with other types of readers
// bigqueryio.Read requires a reflect.Type
type Option interface {
	key() string
	value() interface{}
}

type withTypeOption struct {
	v interface{}
}

func WithType(v interface{}) Option {
	return &withTypeOption{
		v: v,
	}
}

func (opt *withTypeOption) value() interface{} {
	return opt.v
}

func (opt *withTypeOption) key() string {
	return optionType
}

type withProjectTypeOption struct {
	project string
}

func (w withProjectTypeOption) key() string {
	return optionProject
}

func (w withProjectTypeOption) value() interface{} {
	return w.project
}

func WithProject(project string) Option {
	return &withProjectTypeOption{
		project: project,
	}
}

