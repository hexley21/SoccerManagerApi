package repository

import (
	"context"
	"time"

	"github.com/hexley21/soccer-manager/pkg/cache"
	"github.com/hexley21/soccer-manager/pkg/cache/mem"
	"github.com/hexley21/soccer-manager/pkg/cache/ttl"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	countryCacheKey = iota
	localeCacheKey
)

type GlobeRepo interface {
	GetAllCountries(ctx context.Context) ([]Country, error)
	GetAllLocales(ctx context.Context) ([]Locale, error)
}

type pgGlobeRepo struct {
	db           *pgxpool.Pool
	countryCache cache.Cache[int, ttl.ExpirableItem[[]Country]]
	localeCache  cache.Cache[int, ttl.ExpirableItem[[]Locale]]
	ttl          time.Duration // in hours
}

func NewGlobeRepo(db *pgxpool.Pool, cacheTTL time.Duration) *pgGlobeRepo {
	countryInMem := mem.NewInMemoryCache[int, ttl.ExpirableItem[[]Country]]()
	localeInMem := mem.NewInMemoryCache[int, ttl.ExpirableItem[[]Locale]]()

	return &pgGlobeRepo{
		db:           db,
		countryCache: ttl.New(countryInMem),
		localeCache:  ttl.New(localeInMem),
	}
}

const getAllCountries = `-- name: GetAllCountries :many
SELECT code, name FROM countries
`

func (r *pgGlobeRepo) GetAllCountries(ctx context.Context) ([]Country, error) {
	if cachedCountries, ok := r.countryCache.Get(countryCacheKey); ok {
		return cachedCountries.Value, nil
	}

	rows, err := r.db.Query(ctx, getAllCountries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Country
	for rows.Next() {
		var i Country
		if err := rows.Scan(&i.Code, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	r.countryCache.Put(countryCacheKey, ttl.NewItem(items, r.ttl))
	return items, nil
}

const getAllLocales = `-- name: GetAllLocales :many
SELECT code, name FROM locales
`

func (r *pgGlobeRepo) GetAllLocales(ctx context.Context) ([]Locale, error) {
	if cachedLocales, ok := r.localeCache.Get(localeCacheKey); ok {
		return cachedLocales.Value, nil
	}

	rows, err := r.db.Query(ctx, getAllLocales)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Locale
	for rows.Next() {
		var i Locale
		if err := rows.Scan(&i.Code, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	r.localeCache.Put(localeCacheKey, ttl.NewItem(items, r.ttl))
	return items, nil
}
