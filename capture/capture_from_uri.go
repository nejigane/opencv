package capture

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync/atomic"
	"time"
)

// CaptureFromURICreator is a creator of a capture from URI.
type CaptureFromURICreator struct{}

func (c *CaptureFromURICreator) TypeName() string {
	return "capture_from_uri"
}

// CreateSource creates a frame generator using OpenCV video capture.
// URI can be set HTTP address or file path.
//
// Usage of WITH parameters:
//  uri:              [required] a capture data's URI (e.g. /data/test.avi)
//  frame_skip:       the number of frame skip, if set empty or "0"
//                    then read all frames
//  camera_id:        the unique ID of this source if set empty then the ID will be 0
//  next_frame_error: when this source cannot read a new frame, occur error or not
//                    decided by the flag. if the flag set `true` then return error.
//                    default parameter is true.
func (c *CaptureFromURICreator) CreateSource(ctx *core.Context, ioParams *bql.IOParams,
	params data.Map) (core.Source, error) {

	cs, err := createCaptureFromURI(ctx, ioParams, params)
	if err != nil {
		return nil, err
	}
	return core.NewRewindableSource(cs), nil
}

func createCaptureFromURI(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (
	core.Source, error) {

	uri, err := params.Get("uri")
	if err != nil {
		return nil, fmt.Errorf("capture source needs URI")
	}
	uriStr, err := data.AsString(uri)
	if err != nil {
		return nil, err
	}

	fs, err := params.Get("frame_skip")
	if err != nil {
		fs = data.Int(0) // will be ignored
	}
	frameSkip, err := data.AsInt(fs)
	if err != nil {
		return nil, err
	}

	cid, err := params.Get("camera_id")
	if err != nil {
		cid = data.Int(0)
	}
	cameraID, err := data.AsInt(cid)
	if err != nil {
		return nil, err
	}

	endErrFlag, err := params.Get("next_frame_error")
	if err != nil {
		endErrFlag = data.True
	}
	endErr, err := data.AsBool(endErrFlag)
	if err != nil {
		return nil, err
	}

	cs := &captureFromURI{}
	atomic.StoreInt32(&(cs.stop), int32(1))
	cs.uri = uriStr
	cs.frameSkip = frameSkip
	cs.cameraID = cameraID
	cs.endErrFlag = endErr
	return cs, nil
}

type captureFromURI struct {
	vcap bridge.VideoCapture
	// stop is used as atomic bool
	// stop set 0 then means false, set other then means true
	stop int32

	uri        string
	frameSkip  int64
	cameraID   int64
	endErrFlag bool
}

// GenerateStream streams video capture datum. OpenCV video capture read
// frames from URI, user can control frame streaming frequency use
// FrameSkip.
//
// When a capture source is a file-style (e.g. AVI file) and complete to read
// all frames, video capture cannot read a new frame. User can determine to
// occur an error or not to set "next_frame_error". User can also count total
// frame to confirm complete of read file. The number of count is logged.
func (c *captureFromURI) GenerateStream(ctx *core.Context, w core.Writer) error {
	if atomic.LoadInt32(&(c.stop)) == 0 {
		atomic.StoreInt32(&(c.stop), int32(1))
		ctx.Log().Infof("interrupt reading video stream or file and reset: %v",
			c.uri)
		c.vcap.Release()
		c.vcap.Delete()
	}

	c.vcap = bridge.NewVideoCapture()
	if ok := c.vcap.Open(c.uri); !ok {
		return fmt.Errorf("error opening video stream or file: %v", c.uri)
	}

	buf := bridge.NewMatVec3b()
	defer buf.Delete()

	cnt := 0
	ctx.Log().Infof("start reading video stream of file: %v", c.uri)
	atomic.StoreInt32(&(c.stop), int32(0))
	defer atomic.StoreInt32(&(c.stop), int32(1))
	for atomic.LoadInt32(&(c.stop)) == 0 {
		cnt++
		if ok := c.vcap.Read(buf); !ok {
			ctx.Log().Infof("total read frames count is %d", cnt-1)
			if c.endErrFlag {
				return fmt.Errorf("cannot reed a new frame")
			}
			break
		}
		if c.frameSkip > 0 {
			c.vcap.Grab(int(c.frameSkip))
		}

		var m = data.Map{
			"capture":   data.Blob(buf.Serialize()),
			"camera_id": data.Int(c.cameraID),
		}
		now := time.Now()
		t := core.Tuple{
			Data:          m,
			Timestamp:     now,
			ProcTimestamp: now,
			Trace:         []core.TraceEvent{},
		}
		err := w.Write(ctx, &t)
		if err == core.ErrSourceRewound || err == core.ErrSourceStopped {
			return err
		}
	}
	return nil
}

func (c *captureFromURI) Stop(ctx *core.Context) error {
	atomic.StoreInt32(&(c.stop), int32(1))
	c.vcap.Delete()
	return nil
}
