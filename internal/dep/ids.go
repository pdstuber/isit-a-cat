package dep

type IDGenerator interface {
	GenerateID() string
}

type HasIDGenerator interface {
	IDGenerator() IDGenerator
}

func (d AppDependencies) IDGenerator() IDGenerator {
	return d.idGenerator
}
