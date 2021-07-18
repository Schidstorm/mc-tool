package job

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
)

func CopyAllWithProgress(title string, totalSize int64, dst io.Writer, src io.Reader) error {
	var bar *pb.ProgressBar = nil

	bar = pb.Start64(totalSize)
	bar.Set(pb.Bytes, true)
	bar.SetTemplateString(fmt.Sprintf(`%s: %s`, title, pb.Default))

	if srcCloser, ok := src.(io.Closer); ok {
		defer srcCloser.Close()
	}

	for {
		copied, err := io.CopyN(dst, src, 32*1024)
		bar.Add64(copied)

		if err == nil {
			continue
		}
		if err == io.EOF {
			bar.Finish()
			break
		}
		return err
	}

	return nil
}
