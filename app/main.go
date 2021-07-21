package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
)

type Candidate struct {
	Cpf   string
	Name  string
	Score float64
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("./database"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/approvals/1")
	})

	router.GET("/approvals/:page", func(c *gin.Context) {
		page, err := strconv.Atoi(c.Param("page"))
		if err != nil {
			c.String(http.StatusOK, "Invalid page.")
		} else {
			candidates := []Candidate{}

			_ = db.View(func(txn *badger.Txn) error {
				opts := badger.DefaultIteratorOptions
				opts.PrefetchSize = 10
				it := txn.NewIterator(opts)
				defer it.Close()

				i := 0
				for it.Rewind(); it.Valid(); it.Next() {
					i++
					if i < (page-1)*10 {
						continue
					}
					if i >= page*10 {
						break
					}

					item := it.Item()
					k := item.Key()
					candidates = append(candidates, Candidate{Cpf: string(k)})
				}
				return nil
			})

			if len(candidates) > 0 {
				c.HTML(http.StatusOK, "approvals.tmpl", gin.H{
					"title":      "Approved candidates",
					"candidates": candidates,
					"nextPage":   page + 1,
				})
			} else {
				c.String(http.StatusOK, "Invalid page.")
			}
		}
	})

	router.GET("/candidate/:cpf", func(c *gin.Context) {
		cpf := c.Param("cpf")
		var name string
		var score float64

		err := db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(cpf))
			if err != nil {
				return err
			}

			err = item.Value(func(val []byte) error {
				var candidate Candidate
				err = json.Unmarshal(val, &candidate)
				if err != nil {
					return err
				}
				name = candidate.Name
				score = math.Round(candidate.Score*100) / 100
				return nil
			})
			return nil
		})

		if err != nil {
			c.String(http.StatusOK, "Invalid page.")
		} else {
			c.HTML(http.StatusOK, "candidate.tmpl", gin.H{
				"title": "Candidate",
				"name":  name,
				"score": score,
			})
		}
	})

	router.Run(":8080")
}
