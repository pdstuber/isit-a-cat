package dep

type CanForwardDependencies interface {
	Forward() *AppDependencies
}

func (d *AppDependencies) Forward() *AppDependencies {
	return d
}
