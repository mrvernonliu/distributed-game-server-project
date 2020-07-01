package actioninterfaces

type Action int

type ActionUpdate struct {
	Action Action
	Payload int
}