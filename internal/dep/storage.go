package dep

type StorageWriter interface {
	WriteToBucketObject(objectID string, data []byte) error
}

type StorageReader interface {
	ReadFromBucketObject(objectId string) ([]byte, error)
}

type StorageReaderWriter interface {
	StorageReader
	StorageWriter
}

type HasStorageReader interface {
	StorageReader() StorageReader
}

type HasStorageWriter interface {
	StorageWriter() StorageWriter
}

func (d AppDependencies) StorageReader() StorageReader {
	return d.storageService
}

func (d AppDependencies) StorageWriter() StorageWriter {
	return d.storageService
}
