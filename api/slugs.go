package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/psilin/sluggin/core"
	"github.com/psilin/sluggin/db"
)

type Env struct {
	dbase *sqlx.DB
}

func (e *Env) Init(dsn string) error {
	res, err := db.InitDb(dsn)
	if err != nil {
		return fmt.Errorf("fatal error connecting to db: %v", err)
	}
	e.dbase = res
	return nil
}

func (e *Env) Teardown() error {
	err := db.CloseDb(e.dbase)
	if err != nil {
		return fmt.Errorf("fatal error closing db: %v", err)
	}
	return nil
}

func (e *Env) GetSlugs(c *gin.Context) {
	offset := c.DefaultQuery("offset", "0")
	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect offset value"})
		return
	}

	limit := c.DefaultQuery("limit", "10")
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect limit value"})
		return
	}

	slugs, err := db.GetSlugs(e.dbase, limitInt, offsetInt)
	if err != nil {
		log.Printf("DB error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"limit": limitInt, "offset": offsetInt, "slugs": slugs})
}

func (e *Env) AddSlug(c *gin.Context) {
	in := core.Slug{}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := db.AddSlug(e.dbase, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": fmt.Sprintf("slug added to db with id %v", id)})
}

func (e *Env) GetSlugById(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id value"})
		return
	}
	out, err := db.GetSlugById(e.dbase, idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("id %v not found", idInt)})
	}
	c.JSON(http.StatusOK, gin.H{"result": *out})
}

func (e *Env) DeleteSlugById(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id value"})
		return
	}
	err = db.DeleteSlugById(e.dbase, idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": fmt.Sprintf("deleted slug with id %v", idInt)})
}
