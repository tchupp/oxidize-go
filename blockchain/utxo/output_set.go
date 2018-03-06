package utxo

import "github.com/tclchiam/oxidize-go/blockchain/entity"

type OutputSet struct {
	set map[string]entity.Outputs
}

func NewOutputSet() *OutputSet {
	return &OutputSet{
		set: make(map[string]entity.Outputs, 0),
	}
}

func (s *OutputSet) Add(txId *entity.Hash, output *entity.Output) *OutputSet {
	newSet := copySet(s)
	newSet[txId.String()] = newSet[txId.String()].Add(output)
	return &OutputSet{newSet}
}

func (s *OutputSet) AddMany(txId *entity.Hash, outputs []*entity.Output) *OutputSet {
	newSet := copySet(s)
	newSet[txId.String()] = newSet[txId.String()].Append(outputs)
	return &OutputSet{newSet}
}

func (s *OutputSet) Plus(other *OutputSet) *OutputSet {
	newSet := copySet(s)
	for txId, outputs := range other.set {
		newSet[txId] = newSet[txId].Append(outputs)
	}
	return &OutputSet{newSet}
}

func (s *OutputSet) Remove(txId *entity.Hash, toRemove *entity.Output) *OutputSet {
	newSet := copySet(s)
	if index, ok := indexOf(newSet[txId.String()], toRemove); ok {
		newSet[txId.String()] = remove(newSet, txId, index)
	}
	if 0 == len(newSet[txId.String()]) {
		delete(newSet, txId.String())
	}
	return &OutputSet{newSet}
}

func (s *OutputSet) FilterByOutput(predicate func(output *entity.Output) bool) *OutputSet {
	newSet := make(map[string]entity.Outputs, len(s.set))

	for txId, outputs := range s.set {
		for _, output := range outputs {
			if predicate(output) {
				newSet[txId] = newSet[txId].Add(output)
			}
		}
	}
	return &OutputSet{set: newSet}
}

func (s *OutputSet) ForEach(consumer func(*entity.Hash, *entity.Output)) {
	for txId, outputs := range s.set {
		for _, output := range outputs {
			consumer(entity.NewHashOrPanic(txId), output)
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

func copySet(s *OutputSet) map[string]entity.Outputs {
	set := make(map[string]entity.Outputs, len(s.set))
	for txId, outputs := range s.set {
		set[txId] = outputs
	}
	return set
}

func remove(set map[string]entity.Outputs, txId *entity.Hash, index int) []*entity.Output {
	outputs := set[txId.String()]

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
