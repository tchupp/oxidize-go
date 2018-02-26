package utxo

import "github.com/tclchiam/oxidize-go/blockchain/entity"

type OutputSet struct {
	set map[*entity.Hash][]*entity.Output
}

func NewOutputSet() *OutputSet {
	return &OutputSet{
		set: make(map[*entity.Hash][]*entity.Output, 0),
	}
}

func (s *OutputSet) Add(txId *entity.Hash, output *entity.Output) *OutputSet {
	newSet := copySet(s)
	newSet[txId] = append(newSet[txId], output)
	return &OutputSet{newSet}
}

func (s *OutputSet) Remove(txId *entity.Hash, toRemove *entity.Output) *OutputSet {
	newSet := copySet(s)
	if index, ok := indexOf(newSet[txId], toRemove); ok {
		newSet[txId] = remove(newSet, txId, index)
	}
	return &OutputSet{newSet}
}

func (s *OutputSet) FilterByOutput(predicate func(output *entity.Output) bool) *OutputSet {
	newSet := make(map[*entity.Hash][]*entity.Output, len(s.set))

	for txId, outputs := range s.set {
		for _, output := range outputs {
			if predicate(output) {
				newSet[txId] = append(newSet[txId], output)
			}
		}
	}
	return &OutputSet{set: newSet}
}

func (s *OutputSet) ForEach(consumer func(*entity.Hash, *entity.Output)) {
	for txId, outputs := range s.set {
		for _, output := range outputs {
			consumer(txId, output)
		}
	}
}

func (s *OutputSet) ForEachOutput(consumer func(*entity.Output)) {
	for _, outputs := range s.set {
		for _, output := range outputs {
			consumer(output)
		}
	}
}

func (s *OutputSet) Copy() *OutputSet {
	return &OutputSet{copySet(s)}
}

func copySet(s *OutputSet) map[*entity.Hash][]*entity.Output {
	set := make(map[*entity.Hash][]*entity.Output, len(s.set))
	for txId, outputs := range s.set {
		set[txId] = outputs
	}
	return set
}

func remove(set map[*entity.Hash][]*entity.Output, txId *entity.Hash, index int) []*entity.Output {
	outputs := set[txId]

	return append(outputs[:index], outputs[index+1:]...)
}

func indexOf(outputs []*entity.Output, toRemove *entity.Output) (int, bool) {
	for index, output := range outputs {
		if output.IsEqual(toRemove) {
			return index, true
		}
	}

	return -1, false
}
