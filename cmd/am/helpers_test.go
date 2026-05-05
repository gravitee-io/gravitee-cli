package am

type discardWriter struct{}

func (w *discardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
