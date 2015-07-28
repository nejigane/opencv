package mjpegserv

import (
	"fmt"
	"io/ioutil"
	"os"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// DebugJPEGWriterCreator is a creator of JPEG Writer.
type DebugJPEGWriterCreator struct{}

// CreateSink creates a JPEG output sink, which output converted JPEG from
// `cv::Mat_<cv::Vec3b>`.
// Usage of WITH parameters:
//  output:  output directory, if empty then files are output to the current
//           directory.
//  quality: the quality of conversion, if empty then set 50
//
// Example:
//  when a creation query is
//    `CREATE SINK jpeg_files TYPE jpeg_debug WITH output='temp', quality=50`
//  then JPEG files are output to "temp" directory.
func (c *DebugJPEGWriterCreator) CreateSink(ctx *core.Context, ioParams *bql.IOParams,
	params data.Map) (core.Sink, error) {

	output, err := params.Get("output")
	if err != nil {
		output = data.String(".")
	}
	outputDir, err := data.AsString(output)
	if err != nil {
		return nil, err
	}

	quality, err := params.Get("quality")
	if err != nil {
		quality = data.Int(50)
	}
	q, err := data.AsInt(quality)
	if err != nil {
		return nil, err
	}

	s := &debugJPEGSink{}
	s.outputDir = outputDir
	s.jpegQuality = int(q)
	return s, nil
}

func (c *DebugJPEGWriterCreator) TypeName() string {
	return "jpeg_debug"
}

type debugJPEGSink struct {
	outputDir   string
	jpegQuality int
}

// Write output JPEG files to the directory which is set `WITH` "output"
// parameter. Input tuple is required to have follow `data.Map`
//
//  data.Map{
//    "name": [output file name] (`data.String`),
//    "img" : [image binary data] (`data.Blob`),
//  }
func (s *debugJPEGSink) Write(ctx *core.Context, t *core.Tuple) error {
	name, err := t.Data.Get("name")
	if err != nil {
		return err
	}
	nameStr, err := data.ToString(name)
	if err != nil {
		return err
	}

	img, err := t.Data.Get("img")
	if err != nil {
		return err
	}
	imgByte, err := data.AsBlob(img)
	if err != nil {
		return err
	}
	imgp := bridge.DeserializeMatVec3b(imgByte)
	defer imgp.Delete()

	fileName := fmt.Sprintf("%v/%v.jpg", s.outputDir, nameStr)
	ioutil.WriteFile(fileName, imgp.ToJpegData(s.jpegQuality), os.ModePerm)
	return nil
}

func (s *debugJPEGSink) Close(ctx *core.Context) error {
	return nil
}
