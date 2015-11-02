package bridge

/*
#cgo linux pkg-config: opencv pficommon
#cgo darwin pkg-config: opencv pficommon
#include <stdlib.h>
#include "util.h"
#include "opencv_bridge.h"
*/
import "C"
import (
	"sync"
	"unsafe"
)

const (
	// CvCapPropFrameWidth is OpenCV parameter of Frame Width
	CvCapPropFrameWidth = 3
	// CvCapPropFrameHeight is OpenCV parameter of Frame Height
	CvCapPropFrameHeight = 4
	// CvCapPropFps is OpenCV parameter of FPS
	CvCapPropFps = 5
)

// MatVec3b is a bind of `cv::Mat_<cv::Vec3b>`
type MatVec3b struct {
	p C.MatVec3b
}

// NewMatVec3b returns a new MatVec3b.
func NewMatVec3b() MatVec3b {
	return MatVec3b{p: C.MatVec3b_New()}
}

// ToJpegData convert to JPEG data.
func (m *MatVec3b) ToJpegData(quality int) []byte {
	b := C.MatVec3b_ToJpegData(m.p, C.int(quality))
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

// Serialize object.
func (m *MatVec3b) Serialize() []byte {
	b := C.MatVec3b_Serialize(m.p)
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

// DeserializeMatVec3b deserializes object.
func DeserializeMatVec3b(m []byte) MatVec3b {
	b := toByteArray(m)
	return MatVec3b{p: C.MatVec3b_Deserialize(b)}
}

// Delete object.
func (m *MatVec3b) Delete() {
	C.MatVec3b_Delete(m.p)
	m.p = nil
}

// CopyTo copies MatVec3b.
func (m *MatVec3b) CopyTo(dst MatVec3b) {
	C.MatVec3b_CopyTo(m.p, dst.p)
}

// Empty returns the MatVec3b is empty or not.
func (m *MatVec3b) Empty() bool {
	isEmpty := C.MatVec3b_Empty(m.p)
	return isEmpty != 0
}

// VideoCapture is a bind of `cv::VideoCapture`.
type VideoCapture struct {
	p C.VideoCapture
}

// NewVideoCapture returns a new video capture.
func NewVideoCapture() VideoCapture {
	return VideoCapture{p: C.VideoCapture_New()}
}

// Delete object.
func (v *VideoCapture) Delete() {
	C.VideoCapture_Delete(v.p)
	v.p = nil
}

// Open a video data and prepares to start capturing.
func (v *VideoCapture) Open(uri string) bool {
	cURI := C.CString(uri)
	defer C.free(unsafe.Pointer(cURI))
	return C.VideoCapture_Open(v.p, cURI) != 0
}

// OpenDevice opens a video device and prepares to start capturing.
func (v *VideoCapture) OpenDevice(device int) bool {
	return C.VideoCapture_OpenDevice(v.p, C.int(device)) != 0
}

// Release video capture object.
func (v *VideoCapture) Release() {
	C.VideoCapture_Release(v.p)
}

// Set parameter with property (=key).
func (v *VideoCapture) Set(prop int, param int) {
	C.VideoCapture_Set(v.p, C.int(prop), C.int(param))
}

// IsOpened returns the video capture opens a file(or device) or not.
func (v *VideoCapture) IsOpened() bool {
	isOpened := C.VideoCapture_IsOpened(v.p)
	return isOpened != 0
}

// Read set frame to argument MatVec3b, returns `false` when the video capture
// cannot read frame.
func (v *VideoCapture) Read(m MatVec3b) bool {
	return C.VideoCapture_Read(v.p, m.p) != 0
}

// Grab `skip` count frames.
func (v *VideoCapture) Grab(skip int) {
	C.VideoCapture_Grab(v.p, C.int(skip))
}

// VideoWriter is a bind of `cv::VideoWriter`.
type VideoWriter struct {
	mu sync.RWMutex
	p  C.VideoWriter
}

// NewVideoWriter returns a new video writer.
func NewVideoWriter() VideoWriter {
	return VideoWriter{p: C.VideoWriter_New()}
}

// Delete object.
func (vw *VideoWriter) Delete() {
	C.VideoWriter_Delete(vw.p)
	vw.p = nil
}

// Open a video writer.
func (vw *VideoWriter) Open(name string, fps float64, width int, height int) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.VideoWriter_Open(vw.p, cName, C.double(fps), C.int(width), C.int(height))
}

// OpenWithMat opens video writer.
func (vw *VideoWriter) OpenWithMat(name string, fps float64, img MatVec3b) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.VideoWriter_OpenWithMat(vw.p, cName, C.double(fps), img.p)
}

// IsOpened returns the video writer opens a file or not.
func (vw *VideoWriter) IsOpened() bool {
	isOpend := C.VideoWriter_IsOpened(vw.p)
	return isOpend != 0
}

// Write the image to file.
func (vw *VideoWriter) Write(img MatVec3b) {
	vw.mu.Lock()
	defer vw.mu.Unlock()
	C.VideoWriter_Write(vw.p, img.p)
}
