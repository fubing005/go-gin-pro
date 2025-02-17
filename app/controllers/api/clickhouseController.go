package api

/*
	https://cloud.tencent.com/developer/search/article-%E5%A4%A7%E6%95%B0%E6%8D%AEClickHouse
	go get github.com/ClickHouse/clickhouse-go/v2@v2.1.1 适配clickhouse:21.3.20
*/

import (
	"context"
	"net/http"
	"shalabing-gin/global"
	"time"

	"github.com/gin-gonic/gin"
)

type ClickhouseController struct{}

func (con ClickhouseController) CreateTable(c *gin.Context) {
	query := `
		CREATE TABLE users
		(
			id UInt32,
			username String,
			created DateTime
		) ENGINE = MergeTree()
		ORDER BY id;
	`
	err := global.App.Clickhouse.Exec(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "failed to create table"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Table created"})
}

func (con ClickhouseController) InsertHandler(c *gin.Context) {
	type Request struct {
		ID       int       `form:"id"`
		Username string    `form:"username"`
		Created  time.Time `form:"created"`
	}
	var req Request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Created = time.Now()
	query := "INSERT INTO users (id, username, created) VALUES (?, ?, ?)"
	if err := global.App.Clickhouse.Exec(context.Background(), query, req.ID, req.Username, req.Created); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Record inserted"})
}

func (con ClickhouseController) QueryHandler(c *gin.Context) {
	rows, err := global.App.Clickhouse.Query(context.Background(), "SELECT id, username FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var results []gin.H
	for rows.Next() {
		var id uint32
		var username string
		if err := rows.Scan(&id, &username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, gin.H{"id": id, "username": username})
	}
	c.JSON(http.StatusOK, results)
}

func (con ClickhouseController) UpdateHandler(c *gin.Context) {
	type Request struct {
		ID       int    `form:"id"`
		Username string `form:"username"`
	}
	var req Request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "ALTER TABLE users UPDATE username = ? WHERE id = ?"
	if err := global.App.Clickhouse.Exec(context.Background(), query, req.Username, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Record updated"})
}

func (con ClickhouseController) DeleteHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	query := "ALTER TABLE users DELETE WHERE id = ?"
	if err := global.App.Clickhouse.Exec(context.Background(), query, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Record deleted"})
}
