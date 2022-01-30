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

type dbInternalSlug struct {
	Id       int    `db:"id"`
	ExtId    int    `db:"extid"`
	Title    string `db:"title"`
	Slug     string `db:"slug"`
	Url      string `db:"url"`
	Locale   string `db:"locale"`
	Products string `db:"products"`
	Topics   string `db:"topics"`
	Summary  string `db:"summary"`
}

func convertToOut(in *dbInternalSlug) *core.OutSlug {
	var out core.OutSlug
	out.DbId = in.Id
	out.Slg.Id = in.ExtId
	out.Slg.Title = in.Title
	out.Slg.Slug = in.Slug
	out.Slg.Url = in.Url
	out.Slg.Locale = in.Locale
	r := strings.Split(in.Products, "|")
	out.Slg.Products = append(out.Slg.Products, r...)
	r = strings.Split(in.Topics, "|")
	out.Slg.Topics = append(out.Slg.Topics, r...)
	out.Slg.Summary = in.Summary
	return &out
}

func GetSlugs(db *sqlx.DB, limit, offset int) ([]core.OutSlug, error) {
	slugs := []core.OutSlug{}
	slugsDB := []dbInternalSlug{}
	err := db.Select(&slugsDB, "SELECT * FROM slugs limit $1 offset $2", limit, offset)
	if err != nil {
		return slugs, err
	}

	for _, sdb := range slugsDB {
		out := convertToOut(&sdb)
		slugs = append(slugs, *out)
	}
	return slugs, nil
}

func GetSlugById(db *sqlx.DB, id int) (*core.OutSlug, error) {
	slugDb := dbInternalSlug{}
	err := db.Get(&slugDb, "SELECT * from slugs WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	out := convertToOut(&slugDb)
	return out, nil
}

func DeleteSlugById(db *sqlx.DB, id int) error {
	res, err := db.Exec("DELETE FROM slugs WHERE id=$1", id)
	if err != nil {
		return err
	}
	rowsNum, err_inner := res.RowsAffected()
	if err_inner != nil {
		return err_inner
	}
	if rowsNum == 0 {
		return fmt.Errorf("no slug with id %d in db", id)
	}
	return err
}
