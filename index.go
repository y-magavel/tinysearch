package tinysearch

import (
	"container/list"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type DocumentID int

type Posting struct {
	DocID         DocumentID // ドキュメントID
	Positions     []int      // 用語の出現位置
	TermFrequency int        // ドキュメント内の用語の出現回数
}

func NewPosting(docID DocumentID, positions ...int) *Posting {
	return &Posting{docID, positions, len(positions)}
}

type PostingList struct {
	*list.List
}

func NewPostingList(postings ...*Posting) PostingList {
	l := list.New()
	for _, posting := range postings {
		l.PushBack(posting)
	}
	return PostingList{l}
}

func (pl *PostingList) add(p *Posting) {
	pl.PushBack(p)
}

func (pl *PostingList) last() *Posting {
	e := pl.List.Back()
	if e == nil {
		return nil
	}

	return e.Value.(*Posting)
}

func (pl *PostingList) Add(new *Posting) {
	last := pl.last()
	if last == nil || last.DocID != new.DocID {
		pl.add(new)
		return
	}
	last.Positions = append(last.Positions, new.Positions...)
	last.TermFrequency++
}

func (pl PostingList) MarshalJSON() ([]byte, error) {
	postings := make([]*Posting, 0, pl.Len())

	for e := pl.Front(); e != nil; e = e.Next() {
		postings = append(postings, e.Value.(*Posting))
	}
	return json.Marshal(postings)
}

func (pl *PostingList) UnmarshalJSON(b []byte) error {
	var postings []*Posting
	if err := json.Unmarshal(b, &postings); err != nil {
		return err
	}
	pl.List = list.New()
	for _, posting := range postings {
		pl.add(posting)
	}

	return nil
}

type Index struct {
	Dictionary     map[string]PostingList
	TotalDocsCount int
}

func NewIndex() *Index {
	dict := make(map[string]PostingList)
	return &Index{
		Dictionary:     dict,
		TotalDocsCount: 0,
	}
}

func (p *Posting) String() string {
	return fmt.Sprintf("(DocID: %v, TermFrequency: %v, Positions: %v)", p.DocID, p.TermFrequency, p.Positions)
}

func (pl *PostingList) String() string {
	str := make([]string, 0, pl.Len())
	for e := pl.Front(); e != nil; e = e.Next() {
		str = append(str, e.Value.(*Posting).String())
	}
	return strings.Join(str, "=>")
}

func (idx *Index) String() string {
	var padding int
	keys := make([]string, 0, len(idx.Dictionary))
	for k := range idx.Dictionary {
		l := utf8.RuneCountInString(k)
		if padding < l {
			padding = l
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	strs := make([]string, len(keys))
	format := "  [%-" + strconv.Itoa(padding) + "s] -> %s"
	for i, k := range keys {
		if postingList, ok := idx.Dictionary[k]; ok {
			strs[i] = fmt.Sprintf(format, k, postingList.String())
		}
	}
	return fmt.Sprintf("total documents : %v\ndictionary:\n%v\n", idx.TotalDocsCount, strings.Join(strs, "\n"))
}

type Cursor struct {
	postingsList *PostingList  // cursorがたどっているポスティングリストへの参照
	current      *list.Element // 現在の読み込み位置
}

func (pl PostingList) OpenCursor() *Cursor {
	return &Cursor{
		postingsList: &pl,
		current:      pl.Front(),
	}
}

func (c *Cursor) Next() {
	c.current = c.current.Next()
}

func (c *Cursor) NextDoc(id DocumentID) {
	for !c.Empty() && c.DocId() < id {
		c.Next()
	}
}

func (c *Cursor) Empty() bool {
	if c.current == nil {
		return true
	}
	return false
}

func (c *Cursor) Posting() *Posting {
	return c.current.Value.(*Posting)
}

func (c *Cursor) DocId() DocumentID {
	return c.current.Value.(*Posting).DocID
}

func (c *Cursor) String() string {
	return fmt.Sprint(c.Posting())
}
