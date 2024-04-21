package tinysearch

import (
	"bufio"
	"io"
)

type Indexer struct {
	index     *Index
	tokenizer *Tokenizer
}

func NewIndexer(tokenizer *Tokenizer) *Indexer {
	return &Indexer{
		index:     NewIndex(),
		tokenizer: tokenizer,
	}
}

// ドキュメントをインデックスに追加する処理
func (idxr *Indexer) update(docID DocumentID, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(idxr.tokenizer.SplitFunc)

	var position int
	for scanner.Scan() {
		term := scanner.Text()
		// ポスティングリストの更新
		if postingsList, ok := idxr.index.Dictionary[term]; !ok {
			idxr.index.Dictionary[term] = NewPostingList(NewPosting(docID, position))
		} else {
			postingsList.Add(NewPosting(docID, position))
		}
		position++
	}

	idxr.index.TotalDocsCount++
}
