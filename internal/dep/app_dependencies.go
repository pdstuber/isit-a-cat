package dep

type AppDependencies struct {
	storageService StorageReaderWriter
	idGenerator    IDGenerator
	imagePredictor ImagePredictor
}

func NewAppDependencies() AppDependencies {
	return AppDependencies{}
}

func (d AppDependencies) WithStorageService(storageService StorageReaderWriter) AppDependencies {
	d.storageService = storageService
	return d
}

func (d AppDependencies) WithIDGenerator(idGenerator IDGenerator) AppDependencies {
	d.idGenerator = idGenerator
	return d
}

func (d AppDependencies) WithImagePredictor(imagePredictor ImagePredictor) AppDependencies {
	d.imagePredictor = imagePredictor
	return d
}
