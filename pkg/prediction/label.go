package prediction

// The Label for a tensorflow prediction
type Label struct {
	Index     int    `csv:"index"`
	ClassName string `csv:"class_name"`
}
