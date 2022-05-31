package ddd_repository

import (
	"github.com/dapr/dapr-go-ddd-sdk/ddd"
)

type FindPagingResult[T ddd.Entity] struct {
	Data       *[]T   `json:"data"`
	TotalRows  int64  `json:"totalRows"`
	TotalPages int64  `json:"totalPages"`
	PageNum    int64  `json:"pageNum"`
	PageSize   int64  `json:"pageSize"`
	Filter     string `json:"filter"`
	Sort       string `json:"sort"`
	Error      error  `json:"-"`
	IsFound    bool   `json:"-"`
}

func NewFindPagingResult[T ddd.Entity](data *[]T, totalRows int64, query *FindPagingQuery, err error) *FindPagingResult[T] {
	if data != nil && query != nil {
		return &FindPagingResult[T]{
			Data:       data,
			TotalRows:  totalRows,
			TotalPages: getTotalPage(totalRows, query.PageSize),
			PageNum:    query.PageNum,
			PageSize:   query.PageSize,
			Sort:       query.Sort,
			Filter:     query.Filter,
			IsFound:    totalRows > 0,
			Error:      err,
		}
	}
	return &FindPagingResult[T]{
		Data:       data,
		TotalRows:  totalRows,
		TotalPages: 0,
		PageNum:    0,
		PageSize:   0,
		Sort:       "",
		Filter:     "",
		IsFound:    false,
		Error:      err,
	}
}

func NewFindPagingResultWithError[T ddd.Entity](err error) *FindPagingResult[T] {
	return &FindPagingResult[T]{
		Data:    nil,
		IsFound: false,
		Error:   err,
	}
}

func getTotalPage(totalRows int64, pageSize int64) int64 {
	if pageSize == 0 {
		return 0
	}
	totalPage := totalRows / pageSize
	if totalRows%pageSize > 1 {
		totalPage++
	}
	return totalPage
}

func (f *FindPagingResult[T]) GetError() error {
	return f.Error
}

func (f *FindPagingResult[T]) GetData() *[]T {
	return f.Data
}

func (f *FindPagingResult[T]) GetIsFound() bool {
	return f.IsFound
}

func (f *FindPagingResult[T]) Result() (*FindPagingResult[T], bool, error) {
	return f, f.IsFound, f.Error
}

func (f *FindPagingResult[T]) OnError(onErr OnError) *FindPagingResult[T] {
	if f.Error != nil && onErr != nil {
		f.Error = onErr(f.Error)
	}
	return f
}

func (f *FindPagingResult[T]) OnNotFond(fond OnIsFond) *FindPagingResult[T] {
	if f.Error == nil && !f.IsFound && fond != nil {
		f.Error = fond()
	}
	return f
}

func (f *FindPagingResult[T]) OnSuccess(success OnSuccessList[T]) *FindPagingResult[T] {
	if f.Error == nil && success != nil && f.IsFound {
		f.Error = success(f.Data)
	}
	return f
}

type FindPagingResultOptions[T ddd.Entity] struct {
	Data       *[]T
	TotalRows  int64
	TotalPages int64
	PageNum    int64
	PageSize   int64
	Filter     string
	Sort       string
	Error      error
	IsFound    bool
}

func NewFindPagingResultOptions[T ddd.Entity]() *FindPagingResultOptions[T] {
	return &FindPagingResultOptions[T]{}
}

func (f *FindPagingResultOptions[T]) SetData(data *[]T) *FindPagingResultOptions[T] {
	f.Data = data
	return f
}

func (f *FindPagingResultOptions[T]) SetTotalRows(totalRows int64) *FindPagingResultOptions[T] {
	f.TotalRows = totalRows
	return f
}

func (f *FindPagingResultOptions[T]) SetTotalPages(totalPages int64) *FindPagingResultOptions[T] {
	f.TotalPages = totalPages
	return f
}

func (f *FindPagingResultOptions[T]) SetPageNum(pageNum int64) *FindPagingResultOptions[T] {
	f.PageNum = pageNum
	return f
}

func (f *FindPagingResultOptions[T]) SetPageSize(pageSize int64) *FindPagingResultOptions[T] {
	f.PageSize = pageSize
	return f
}

func (f *FindPagingResultOptions[T]) SetFilter(filter string) *FindPagingResultOptions[T] {
	f.Filter = filter
	return f
}

func (f *FindPagingResultOptions[T]) SetSort(sort string) *FindPagingResultOptions[T] {
	f.Sort = sort
	return f
}

func (f *FindPagingResultOptions[T]) SetError(err error) *FindPagingResultOptions[T] {
	f.Error = err
	return f
}

func (f *FindPagingResultOptions[T]) SetIsFound(isFound bool) *FindPagingResultOptions[T] {
	f.IsFound = isFound
	return f
}
