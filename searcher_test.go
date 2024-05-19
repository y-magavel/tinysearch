package tinysearch

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSearchTopK(t *testing.T) {

	s := NewSearcher("testdata/index")
	actual := s.SearchTopK([]string{"quarrel", "sir"}, 1)

	expected := &TopDocs{2,
		[]*ScoreDoc{
			{2, 1.9657842846620868},
		},
	}

	if diff := cmp.Diff(actual, expected, cmpopts.EquateApprox(0, 1e-9), cmp.AllowUnexported(TopDocs{}, ScoreDoc{})); diff != "" {
		t.Fatalf("got:\n%v\nexpected:%v\n", actual, expected)
	}
}
