package types

type KVMap map[string]string

// IsAlreadyComparedFlag to indicate whether the information of this entity was compared
type IsAlreadyComparedFlag struct {
	Index int
	Check bool
}
