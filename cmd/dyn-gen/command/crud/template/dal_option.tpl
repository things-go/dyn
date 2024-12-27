package {{.Package}}

import (
	"sync"

	rapier "github.com/thinkgos/gorm-rapier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// boolean
const (
	FALSE = "0"
	TRUE  = "1"
)

func Boolean(b bool) string {
	if b {
		return TRUE
	} else {
		return FALSE
	}
}

func ToBool(s string) bool { return s == TRUE }

func NilOr[T any](slices []T) []T {
	if slices == nil {
		return make([]T, 0)
	}
	return slices
}

var ErrRecordNotFound = gorm.ErrRecordNotFound

var (
	// DefaultPerPage 默认页大小
	DefaultPerPage = int64(50)
	// DefaultPageSize 默认最大页大小
	DefaultMaxPerPage = int64(500)
)

// Paginate 分页器
// 分页索引: page >= 1
// 分页大小: perPage >= 1 && <= DefaultMaxPerPage
func Paginate(page, perPage int64, maxPerPages ...int64) clause.Expression {
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
	return clause.Limit{
		Limit:  &limit,
		Offset: offset,
	}
}

// Pagination 分页器
// 分页索引: page >= 1
// 分页大小: perPage >= 1 && <= DefaultMaxPerPage
func Pagination(page, perPage int64, maxPerPages ...int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(Paginate(page, perPage, maxPerPages...))
	}
}

// 限制器
// offset = perPage * (page - 1)
// limit = perPage
// if limit > 0: use limit
// if offset > 0: use offset
// if offset <= 0 and limit <=0: use none
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

type DalCondition = func(db *gorm.DB) *gorm.DB

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

type assignExprPool struct {
	pool *sync.Pool
}

var defaultPool = assignExprPool{
	pool: &sync.Pool{
		New: func() any {
			return &assignExprContainer{
				make([]rapier.AssignExpr, 0, 32),
			}
		},
	},
}

type assignExprContainer struct {
	Exprs []rapier.AssignExpr
}

func (c *assignExprContainer) reset() *assignExprContainer {
	c.Exprs = c.Exprs[:0]
	return c
}

func (p *assignExprPool) Get() *assignExprContainer {
	c := p.pool.Get().(*assignExprContainer)
	return c
}

// Put adds x to the pool.
// NOTE: See Get.
func (p *assignExprPool) Put(c *assignExprContainer) {
	if c == nil {
		return
	}
	p.pool.Put(c.reset())
}