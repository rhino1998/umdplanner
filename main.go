package main

import (
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
	// f, err := os.Open("testudo.json")
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }
	// store, err := testudo.LoadStore(f)

	fmt.Println(testudo.ParseClass("https://ntst.umd.edu/soc/201801/CMSC/CMSC250"))

	/*ch := (testudo.QueryWithExcludedTimes(store.QueryAll(),
		testudo.Duration{
			Start: time.Time{}.AddDate(-1, 0, 2).Add(10 * time.Hour),
			End:   time.Time{}.AddDate(-1, 0, 2).Add(10*time.Hour + 50*time.Minute),
		},
		testudo.Duration{
			Start: time.Time{}.AddDate(-1, 0, 2).Add(8 * time.Hour),
			End:   time.Time{}.AddDate(-1, 0, 2).Add(11*time.Hour + 50*time.Minute),
		},
	)).Evaluate()

	// ch := store.QueryAll().Evaluate()
	for class := range ch {
		fmt.Println(class)
	}*/
}
