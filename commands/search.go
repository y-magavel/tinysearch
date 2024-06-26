package commands

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
	"github.com/y-magavel/tinysearch"
)

var searchCommand = cli.Command{
	Name:      "search",
	Usage:     "search documents",
	ArgsUsage: `<query>`,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "number, n",
			Value: 10,
			Usage: "",
		},
	},
	Action: search,
}

func search(c *cli.Context) error {
	if err := checkArgs(c, 1, exactArgs); err != nil {
		return err
	}
	query := c.Args().Get(0)
	result, err := engine.Search(query, c.Int("number"))
	if err != nil {
		return err
	}
	printResult(result)
	return nil
}

func printResult(results []*tinysearch.SearchResult) {
	if len(results) == 0 {
		fmt.Println("0 match!!")
		return
	}
	s := make([]string, len(results))
	for i, result := range results {
		s[i] = fmt.Sprintf("rank: %3d, score: %4f, title: %s",
			i+1, result.Score, result.Title)
	}
	fmt.Println(strings.Join(s, "\n"))
}
