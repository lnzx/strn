package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const STRN_STATS_URL = "https://orchestrator.strn.pl/stats"

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("cron") == "1" {
		if err := cron(); err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("ok"))
		}
	} else {
		w.Header().Set("content-type", "text/json")

		data, err := getData()
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			bytes, err := json.Marshal(data)
			if err != nil {
				w.Write([]byte(err.Error()))
			} else {
				w.Write(bytes)
			}
		}
	}
}

func getData() (data *Data, err error) {
	doc, err := GetHtml(STRN_STATS_URL)
	if err != nil {
		return nil, err
	}

	var nodes []*Node
	isps := make(map[string]int)
	regions := make(map[string]int)
	count := 0

	doc.Find("tbody").Children().Each(func(i int, tr *goquery.Selection) {
		html, err := tr.Html()
		if err != nil {
			log.Println(err)
			return
		}
		tds := tr.Children()

		id, down := ParseId(tds)
		if !down {
			count++
		}
		nodes = append(nodes, &Node{Id: id, Html: "<tr>" + html + "</tr>"})

		isp := ParseISP(tds)
		if n, ok := isps[isp]; ok {
			n++
			isps[isp] = n
		} else {
			isps[isp] = 1
		}

		_, country := ParseLocation(tds)
		if n, ok := regions[country]; ok {
			n++
			regions[country] = n
		} else {
			regions[country] = 1
		}
	})

	return &Data{
		Count:   count,
		Nodes:   nodes,
		Isps:    sortIsps(isps),
		Regions: sortRegions(regions),
	}, nil
}

func cron() error {
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

type Data struct {
	Count   int       `json:"count"`
	Nodes   []*Node   `json:"nodes"`
	Regions []*Region `json:"regions"`
	Isps    []*Isp    `json:"isps"`
}

type Node struct {
	Id   string `json:"id,omitempty"`
	Html string `json:"html,omitempty"`
}

type Isp struct {
	Name  string `json:"name,omitempty"`
	Count int    `json:"count,omitempty"`
}

type Region struct {
	Name  string `json:"name,omitempty"`
	Count int    `json:"count,omitempty"`
}

// GetHtml request the HTML document
func GetHtml(url string) (doc *goquery.Document, err error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Println(err)
			return
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return goquery.NewDocumentFromReader(res.Body)
}

func ParseId(tds *goquery.Selection) (string, bool) {
	child := tds.First().Children()
	return strings.TrimSpace(child.First().Text()), child.Eq(2).HasClass("down")
}

func ParseISP(tds *goquery.Selection) (isp string) {
	ret, err := tds.Eq(2).Html()
	if err != nil {
		return
	}
	ret = strings.Split(ret, "<br")[1]
	ret = clear(ret)
	return ret
}

func ParseLocation(tds *goquery.Selection) (city string, country string) {
	ret, err := tds.Eq(3).Html()
	if err != nil {
		return
	}
	rets := strings.Split(ret, "<br")
	city = rets[0]
	country = clear(rets[1])
	return
}

func ParseDisk(tds *goquery.Selection) int {
	ret, err := tds.Eq(4).Html()
	if err != nil {
		return 0
	}
	ret = strings.Split(ret, "<br")[0]
	ret = strings.Split(ret, "(")[0]
	ret = strings.TrimSpace(ret)
	val := toGiB(humanReadableToInt(ret))
	return val
}

func ParseTTFB(tds *goquery.Selection) int {
	text := tds.Eq(7).Text()
	h1 := strings.Split(text, "24h")[0]
	h1 = strings.Split(h1, ":")[1]
	h1 = strings.TrimSpace(h1)
	i, err := strconv.Atoi(h1)
	if err != nil {
		return 0
	}
	return i
}

func ParseUpload(tds *goquery.Selection) (traffic int, bandwidth int) {
	ret := tds.Eq(10).Text()
	rets := strings.Split(ret, "@")
	traffic = toGiB(humanReadableToInt(strings.TrimSpace(rets[0])))
	bandwidth = bandwidthToMbps(strings.TrimSpace(strings.Split(strings.TrimSpace(rets[1]), "(")[0]))
	return
}

func ParseWeight(tds *goquery.Selection) int {
	text := strings.TrimSpace(tds.Last().Text())
	if i, err := strconv.Atoi(text); err == nil {
		return i
	}
	return -1
}

func clear(str string) string {
	str = strings.ReplaceAll(str, "/>", "")
	str = strings.ReplaceAll(str, ">", "")
	str = strings.TrimSpace(str)
	return str
}

var sizeSuffix = map[string]int{
	"kb": 1024,
	"mb": 1 << (10 * 2),
	"gb": 1 << (10 * 3),
	"tb": 1 << (10 * 4),
	"pb": 1 << (10 * 5),
}

var bandwidthSuffix = map[string]int{
	"mbps": 1,
	"gbps": 1000,
}

// HumanReadableToInt Converts a human-readable size to int, param value: A string such as "10MB"
func humanReadableToInt(value string) int {
	value = strings.ToLower(value)
	value = strings.Replace(value, "i", "", 1)
	l := len(value)
	suffix := value[l-2 : l]
	if multiplier, ok := sizeSuffix[suffix]; ok && l > 2 {
		value = strings.TrimSpace(value[0 : l-2])
		if size, err := strconv.ParseFloat(value, 32); err == nil {
			n := int(math.Floor(float64(multiplier) * size))
			return n
		}
	}
	log.Println("Invalid size value:", value)
	return 0
}

func bandwidthToMbps(value string) int {
	value = strings.ToLower(value)
	l := len(value)
	suffix := value[l-4 : l]
	if multiplier, ok := bandwidthSuffix[suffix]; ok && l > 4 {
		value = strings.TrimSpace(value[0 : l-4])
		if size, err := strconv.ParseFloat(value, 32); err == nil {
			n := int(math.Floor(float64(multiplier) * size))
			return n
		}
	}
	return 0
}

func toGiB(bytes int) int {
	return bytes / sizeSuffix["gb"]
}

func sortIsps(m map[string]int) []*Isp {
	s := make([]*Isp, 0, len(m))
	for k, v := range m {
		s = append(s, &Isp{
			Name:  k,
			Count: v,
		})
	}

	sort.Slice(s, func(i, j int) bool {
		return s[i].Count > s[j].Count
	})

	return s
}

func sortRegions(m map[string]int) []*Region {
	s := make([]*Region, 0, len(m))
	for k, v := range m {
		s = append(s, &Region{
			Name:  k,
			Count: v,
		})
	}

	sort.Slice(s, func(i, j int) bool {
		return s[i].Count > s[j].Count
	})

	return s
}
