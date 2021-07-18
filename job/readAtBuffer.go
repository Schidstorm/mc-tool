package job

import "io"

type ReadAtBuffer struct {
	buffer []byte
}

func (r *ReadAtBuffer) ReadAt(p []byte, off int64) (n int, err error) {
	end := minInt64(int64(len(r.buffer)), off+int64(len(p)))
	copied := copy(p, r.buffer[off:end])
	if copied != len(p) {
		return copied, io.EOF
	}
	return copied, nil
}

func (r *ReadAtBuffer) Write(p []byte) (n int, err error) {
	r.buffer = append(r.buffer, p...)
	return len(p), nil
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
