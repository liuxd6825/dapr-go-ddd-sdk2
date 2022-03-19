package ddd_repository

type FindResult struct {
	err    error
	data   interface{}
	isFind bool
}

func NewFindResult(data interface{}, isFind bool, err error) *FindResult {
	return &FindResult{
		data:   data,
		isFind: isFind,
		err:    err,
	}
}

func (f *FindResult) GetError() error {
	return f.err
}

func (f *FindResult) GetData() interface{} {
	return f.data
}

func (f *FindResult) GetIsFind() interface{} {
	return f.isFind
}

func (f *FindResult) Result() (interface{}, bool, error) {
	return f.data, f.isFind, f.err
}

func (f *FindResult) OnSuccess(success OnSuccess) *FindResult {
	if f.err == nil && success != nil && f.isFind {
		f.err = success(f.data)
	}
	return f
}

func (f *FindResult) OnError(onErr OnError) *FindResult {
	if f.err != nil && onErr != nil {
		f.err = onErr(f.err)
	}
	return f
}

func (f *FindResult) OnNotFond(fond OnIsFond) *FindResult {
	if f.err != nil && !f.isFind && fond != nil {
		f.err = fond()
	}
	return f
}
