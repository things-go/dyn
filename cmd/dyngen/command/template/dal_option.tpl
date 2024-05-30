package {{.Package}}

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

var (
	// DefaultPerPage 默认页大小
	DefaultPerPage = int64(50)
	// DefaultPageSize 默认最大页大小
	DefaultMaxPerPage = int64(500)
)

// Pagination 分页器
// 分页索引: page >= 1
// 分页大小: perPage >= 1 && <= DefaultMaxPerPage
func Pagination(page, perPage int64, maxPerPages ...int64) func(db *gorm.DB) *gorm.DB {
	maxPerPage := DefaultMaxPerPage
	if len(maxPerPages) > 0 && maxPerPages[0] > 0 {
		maxPerPage = maxPerPages[0]
	}
	if page < 1 {
		page = 1
	}
	switch {
	case perPage < 1:
		perPage = DefaultPerPage
	case perPage > maxPerPage:
		perPage = maxPerPage
	default: // do nothing
	}
	limit, offset := int(perPage), int(perPage*(page-1))
	l := clause.Limit{
		Limit:  &limit,
		Offset: offset,
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(l)
	}
}

// 限制器
// offset = perPage * (page - 1)
// limit = perPage
// if limit > 0: use limit
// if offset > 0: use offset
func Limit(page, perPage int64) func(*gorm.DB) *gorm.DB {
	offset := 0
	if page > 0 {
		offset = int(perPage * (page - 1))
	}
	limit := int(perPage)
	return func(db *gorm.DB) *gorm.DB {
		if offset > 0 || limit > 0 {
			l := clause.Limit{
				Limit:  new(int),
				Offset: offset,
			}
			if offset > 0 {
				l.Offset = offset
			}
			if limit > 0 {
				l.Limit = &limit
			}
			return db.Clauses(l)
		} else {
			return db
		}
	}
}

type DalOption func(*DalConfig)

type DalConfig struct {
	funcs []func(db *gorm.DB) *gorm.DB
}

func (o *DalConfig) TakeOptions(opts ...DalOption) *DalConfig {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithLockingUpdate() DalOption {
	return func(o *DalConfig) {
		o.funcs = append(o.funcs, LockingUpdate())
	}
}

func WithLockingShare() DalOption {
	return func(o *DalConfig) {
		o.funcs = append(o.funcs, LockingShare())
	}
}
func WithScopes(funcs ...func(db *gorm.DB) *gorm.DB) DalOption {
	return func(o *DalConfig) {
		o.funcs = append(o.funcs, funcs...)
	}
}

// LockingUpdate specify the lock strength to UPDATE
func LockingUpdate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
}

// LockingShare specify the lock strength to SHARE
func LockingShare() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(clause.Locking{Strength: "SHARE"})
	}
}
