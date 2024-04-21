package tinysearch

import (
	"reflect"
	"strings"
	"testing"
)

func TestUpdate(t *testing.T) {
	collection := []string{
		"Do you quarrel, sir?",
		"Quarrel sir! no, sir!",
		"No better.",
		"Well, sir",
	}

	indexer := NewIndexer(NewTokenizer())

	for i, doc := range collection {
		indexer.update(DocumentID(i), strings.NewReader(doc))
	}

	actual := indexer.index
	expected := &Index{
		Dictionary: map[string]PostingList{
			"do":     NewPostingList(NewPosting(0, 0)),
			"better": NewPostingList(NewPosting(2, 1)),
			"no": NewPostingList(
				NewPosting(1, 2),
				NewPosting(2, 0),
			),
			"quarrel": NewPostingList(
				NewPosting(0, 2),
				NewPosting(1, 0),
			),
			"sir": NewPostingList(
				NewPosting(0, 3),
				NewPosting(1, 1, 3),
				NewPosting(3, 1),
			),
			"well": NewPostingList(NewPosting(3, 0)),
			"you":  NewPostingList(NewPosting(0, 1)),
		},
		TotalDocsCount: 4,
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("wrong index. \n\nwant: \n%v\n\n got:\n%v\n", actual, expected)
	}
}
