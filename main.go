package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rhino1998/umdplanner/testudo"
)

func main() {
	scrape := flag.Bool("scrape", false, "whether to scrape a new db from testudo")

	flag.Parse()

	if *scrape {
		store := testudo.NewStore()
		f, err := os.Create("testudo.json")
		if err != nil {
			log.Fatalf("%v", err)
		}
		testudo.ScrapeAll("https://ntst.umd.edu/soc", store)
		store.Dump(f)
		f.Close()
		return
	}
	f, err := os.Open("testudo.json")
	if err != nil {
		log.Fatalf("%v", err)
	}
	store, err := testudo.LoadStore(f)

	ch := store.QueryAll().Evaluate(context.Background())

	for class := range ch {
		prereqs := ""
		if len(class.Prereqs) == 0 {
			continue
		}

		for _, req := range class.Prereqs {
			prereqs += req.Code + ", "
		}
		fmt.Printf("[%s]: %s\n", class.Code, prereqs)
	}
	class, _ := store.Get("CMSC412")
	fmt.Printf("%s\n", class.Prerequisite)

}
