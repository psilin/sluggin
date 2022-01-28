package db

import (
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/psilin/sluggin/core"
)

func AddSlug(db *sqlx.DB, slug *core.Slug) (int, error) {
	slugDB := map[string]interface{}{
		"extid":    slug.Id,
		"title":    slug.Title,
		"slug":     slug.Slug,
		"url":      slug.Url,
		"locale":   slug.Locale,
		"products": strings.Join(slug.Products, "|"),
		"topics":   strings.Join(slug.Topics, "|"),
		"summary":  slug.Summary}

	slugState := `INSERT INTO slugs (extid, title, slug, url, locale, products, topics, summary) 
		VALUES (:extid, :title, :slug, :url, :locale, :products, :topics, :summary) 
		RETURNING id`
	stmt, err := db.PrepareNamed(slugState)
	if err != nil {
		return 0, fmt.Errorf("adding slug failed: %v", err)
	}

	var id int
	err = stmt.Get(&id, slugDB)
	if err != nil {
		return 0, fmt.Errorf("adding slug failed: %v", err)
	}

	return id, nil
}
