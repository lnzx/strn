package cron

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/lnzx/strn/api"
	"log"
	"os"
)

func Start() error {
	doc, err := GetHtml(STRN_STATS_URL)
	if err != nil {
		return err
	}

	var nodes []*node
	doc.Find("tbody").Children().Each(func(i int, tr *goquery.Selection) {
		tds := tr.Children()
		id, down := ParseId(tds)
		if !down {
			isp := ParseISP(tds)
			city, country := ParseLocation(tds)
			disk := ParseDisk(tds)
			ttfb := ParseTTFB(tds)
			traffic, bandwidth := ParseUpload(tds)
			weight := ParseWeight(tds)
			nodes = append(nodes, &node{
				id:        id,
				isp:       isp,
				country:   country,
				city:      city,
				disk:      disk,
				ttfb:      ttfb,
				traffic:   traffic,
				bandwidth: bandwidth,
				weight:    weight,
			})
		}
	})
	log.Println("active node count:", len(nodes))
	return batchInsert(nodes)
}

type node struct {
	id        string
	isp       string
	country   string
	city      string
	disk      int
	ttfb      int
	traffic   int
	bandwidth int
	weight    int
}

func batchInsert(nodes []*node) error {
	sql := "INSERT INTO nodes(id,isp,country,city,disk,ttfb,traffic,bandwidth,weight) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)"
	batch := &pgx.Batch{}
	for _, n := range nodes {
		batch.Queue(sql, n.id, n.isp, n.country, n.city, n.disk, n.ttfb, n.traffic, n.bandwidth, n.weight)
	}
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("unable to connect to database")
	}
	defer dbpool.Close()

	br := dbpool.SendBatch(context.Background(), batch)
	if err != nil {
		return err
	}
	err = br.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	return br.Close()
}
