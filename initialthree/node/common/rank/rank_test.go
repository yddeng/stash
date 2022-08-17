package rank

//go test -covermode=count -v -coverprofile=coverage.out -run=. -cpuprofile=rank.p
//go tool cover -html=coverage.out
//go tool pprof rank.p
//go test -v -run=^$ -bench BenchmarkRank -count 10

import (
	"fmt"
	"github.com/schollz/progressbar"
	"github.com/stretchr/testify/assert"
	"github.com/yddeng/sortedset"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func init() {
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

}

func TestBenchmarkRank2(t *testing.T) {
	var r *Rank = NewRank()
	testCount := 50000000
	idRange := 10000000
	{
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			item := r.getRankItem(uint64(idx))
			var score int
			if nil == item {
				score = rand.Int() % 1000000
			} else {
				score = item.value + rand.Int()%10000
			}
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg), len(r.id2Item))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
	}
}

func TestBenchmarkRank1(t *testing.T) {
	var r *Rank = NewRank()
	fmt.Println("TestBenchmarkRank")

	testCount := 10000000
	idRange := 10000000

	{
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			score := rand.Int()%1000000 + 1
			//fmt.Println(i, idx, score)
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	{
		testCount := 10000000
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			score := rand.Int()%1000000 + 1
			//fmt.Println(idx, score)
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	{

		testCount := 10000000

		bar := progressbar.New(int(testCount))
		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := (rand.Int() % len(r.id2Item)) + 1
			item := r.id2Item[uint64(idx)]
			score := rand.Int()%10000 + 1
			score = item.value + score
			//fmt.Println(idx, score)
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	{

		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			r.GetRankPercent(uint64(idx))
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
	}

	{
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			r.GetRank(uint64(idx))
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
	}
}

func TestRank2(t *testing.T) {
	fmt.Println("TestRank")

	var r *Rank = NewRank()
	fmt.Println("TestBenchmarkRank")

	f := 0.0099
	fmt.Println(int32(f * 100))

	testCount := 200
	idRange := 1000
	ids := make([]uint64, 0, testCount)

	for i := 0; i < testCount; i++ {
		idx := i%idRange + 1
		ids = append(ids, uint64(idx))
		score := rand.Int() % 10000
		r.UpdateScore(uint64(idx), score)
	}

	r.Show()

	for _, id := range ids {
		score := r.GetScore(id)
		idx := r.GetRank(id)
		perc := r.GetRankPercent(id)
		fmt.Println(id, score, idx, perc)
	}
}

func BenchmarkRank(b *testing.B) {
	var r *Rank = NewRank()
	for i := 0; i < b.N; i++ {
		idx := (i % 1000000) + 1
		score := rand.Int()
		r.UpdateScore(uint64(idx), score)
	}
}

func TestRank_GetTopN(t *testing.T) {
	var r *Rank = NewRank()
	testCount := 200
	idRange := 1000

	for i := 0; i < testCount; i++ {
		idx := i%idRange + 1
		score := rand.Int() % 10000
		r.UpdateScore(uint64(idx), score)
	}

	r.Show()

	ids := r.GetTopN(50)
	fmt.Println(ids, "--", len(ids))

	ids = r.GetTopN(220)
	fmt.Println(ids, "--", len(ids))

}

func TestRank_GetScoreByIdx(t *testing.T) {
	var r *Rank = NewRank()
	testCount := 200
	idRange := 1000

	for i := 0; i < testCount; i++ {
		idx := i%idRange + 1
		score := rand.Int() % 10000
		r.UpdateScore(uint64(idx), score)
	}

	r.Show()

	fmt.Println(r.GetScoreByIdx(5))
	fmt.Println(r.GetScoreByIdx(200))
	fmt.Println(r.GetScoreByIdx(220))
}

func TestRank_Delete(t *testing.T) {
	var r *Rank = NewRank()
	testCount := 20
	idRange := 1000

	for i := 0; i < testCount; i++ {
		idx := i%idRange + 1
		score := rand.Int() % 10000
		r.UpdateScore(uint64(idx), score)
	}

	r.Show()

	id, _ := r.GetScoreByIdx(2)
	r.Delete(id)

	id, _ = r.GetScoreByIdx(1)
	r.Delete(id)

	id, _ = r.GetScoreByIdx(18)
	r.Delete(id)

	r.Show()
}

func TestRank222(t *testing.T) {
	var r *Rank = NewRank()

	r.UpdateScore(1, 1)
	r.UpdateScore(2, 1)
	r.UpdateScore(3, 2)

	fmt.Println(r.GetRank(1))
	fmt.Println(r.GetRank(2))
	fmt.Println(r.GetRank(3))
	r.Show()
}

type Score int

func (this Score) Less(other interface{}) bool {
	return this >= other.(Score)
}

func TestRank3(t *testing.T) {
	l := sortedset.New()
	testCount := 10000000
	{
		bar := progressbar.New(int(testCount))
		now := time.Now()
		for i := 1; i <= testCount; i++ {
			score := rand.Int()
			l.Set(sortedset.Key(fmt.Sprintf("%d", i)), Score(score))
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(now).String())
	}
	{
		bar := progressbar.New(int(testCount))
		now := time.Now()
		for i := 1; i <= testCount; i++ {
			score := rand.Int()
			l.Set(sortedset.Key(fmt.Sprintf("%d", i)), Score(score))
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(now).String())
	}
	{
		bar := progressbar.New(int(testCount))
		now := time.Now()
		for i := 1; i <= testCount; i++ {
			l.GetRank(sortedset.Key(fmt.Sprintf("%d", i)))
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(now).String())
	}
}

func BenchmarkRank_Update(b *testing.B) {
	var r *Rank = NewRank()
	for i := 0; i < b.N; i++ {
		idx := (i % 1000000) + 1
		score := rand.Int()
		r.UpdateScore(uint64(idx), score)
	}
}
