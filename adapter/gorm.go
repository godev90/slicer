package adapter

import (
	"context"

	"gorm.io/gorm"
)

type GormAdapter struct {
	db    *gorm.DB
	model Tabler
}

func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{db: db}
}

func (g *GormAdapter) WithContext(ctx context.Context) QueryAdapter {
	return &GormAdapter{db: g.db.WithContext(ctx), model: g.model}
}

func (g *GormAdapter) UseModel(m Tabler) QueryAdapter {
	return &GormAdapter{db: g.db.Model(m), model: m}
}

func (g *GormAdapter) Model() Tabler {
	return g.model
}

func (g *GormAdapter) Where(query any, args ...any) QueryAdapter {
	if other, ok := query.(*GormAdapter); ok {
		return &GormAdapter{
			db:    g.db.Where(other.db),
			model: g.model,
		}
	}

	return &GormAdapter{
		db:    g.db.Where(query, args...),
		model: g.model,
	}
}

func (g *GormAdapter) Or(query any, args ...any) QueryAdapter {
	return &GormAdapter{db: g.db.Or(query, args...), model: g.model}
}

func (g *GormAdapter) Select(fields []string) QueryAdapter {
	return &GormAdapter{db: g.db.Select(fields), model: g.model}
}

func (g *GormAdapter) Limit(limit int) QueryAdapter {
	return &GormAdapter{db: g.db.Limit(limit), model: g.model}
}

func (g *GormAdapter) Offset(offset int) QueryAdapter {
	return &GormAdapter{db: g.db.Offset(offset), model: g.model}
}

func (g *GormAdapter) Order(order string) QueryAdapter {
	return &GormAdapter{db: g.db.Order(order), model: g.model}
}

func (g *GormAdapter) Clone() QueryAdapter {
	return &GormAdapter{db: g.db.Session(&gorm.Session{NewDB: true}), model: g.model}
}

func (g *GormAdapter) Scope(target *int64) QueryAdapter {
	return &GormAdapter{db: g.db.Scopes(), model: g.model}
}

func (g *GormAdapter) Count(target *int64) error {
	return g.db.Session(&gorm.Session{}).Count(target).Error
}

func (g *GormAdapter) Scan(dest any) error {
	if debug {
		return g.db.Debug().Find(dest).Error
	}

	return g.db.Find(dest).Error
}
