package actions

type Action interface {
	Execute() (bool, error)
}
