package sequence

import (
	"math/rand"
	"time"
)

type Sequence struct {
	Data []int
}

// {2, 3} -> [0, 0, 1, 1, 1]
func NewSequence(counts []int) *Sequence {
	totalOps := 0
	for _, count := range counts {
		totalOps += count
	}

	data := make([]int, 0, totalOps)
	for jobIdx, count := range counts {
		for i := 0; i < count; i++ {
			data = append(data, jobIdx)
		}
	}

	return &Sequence{
		Data: data,
	}
}

func (s *Sequence) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(s.Data), func(i, j int) {
		s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
	})
}

func (s *Sequence) Swap(i, j int) {
	if i >= 0 && i < len(s.Data) && j >= 0 && j < len(s.Data) {
		s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
	}
}

func (s *Sequence) Get(index int) int {
	return s.Data[index]
}

func (s *Sequence) Len() int {
	return len(s.Data)
}

func (s *Sequence) Clone() *Sequence {
	newData := make([]int, len(s.Data))
	copy(newData, s.Data)
	return &Sequence{Data: newData}
}
