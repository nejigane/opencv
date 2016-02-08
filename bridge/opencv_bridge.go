package bridge

/*
#cgo linux pkg-config: opencv
#cgo darwin pkg-config: opencv
#include <stdlib.h>
#include "util.h"
#include "opencv_bridge.h"
*/
import "C"
import (
	"reflect"
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

// CMatVec3b is an alias for C pointer.
type CMatVec3b C.MatVec3b

// MatVec3b is a bind of `cv::Mat_<cv::Vec3b>`
type MatVec3b struct {
	p C.MatVec3b
}

// GetCPointer returns C pointer of MatVec3b.
func (m *MatVec3b) GetCPointer() C.MatVec3b {
	return m.p
}

// NewMatVec3b returns a new MatVec3b.
func NewMatVec3b() MatVec3b {
	return MatVec3b{p: C.MatVec3b_New()}
}

// NewMatVec3bWithCPointer return a new MatVec3b with argument C pointer.
func NewMatVec3bWithCPointer(p CMatVec3b) MatVec3b {
	return MatVec3b{p: C.MatVec3b(p)}
}

// ToJpegData convert to JPEG data.
func (m *MatVec3b) ToJpegData(quality int) []byte {
	b := C.MatVec3b_ToJpegData(m.p, C.int(quality))
	defer C.ByteArray_Release(b)
	return toGoBytes(b)
}

// Delete object.
func (m *MatVec3b) Delete() {
	C.MatVec3b_Delete(m.p)
	m.p = nil
}

// CopyTo copies MatVec3b.
func (m *MatVec3b) CopyTo(dst *MatVec3b) {
	C.MatVec3b_CopyTo(m.p, dst.p)
}

// Empty returns the MatVec3b is empty or not.
func (m *MatVec3b) Empty() bool {
	isEmpty := C.MatVec3b_Empty(m.p)
	return isEmpty != 0
}

// ToRawData converts MatVec3b to RawData.
func (m *MatVec3b) ToRawData() (int, int, []byte) {
	r := C.MatVec3b_ToRawData(m.p)
	return int(r.width), int(r.height), toGoBytes(r.data)
}

// ToMatVec3b converts RawData to MatVec3b. Returned MatVec3b is required to
// delete after using.
func ToMatVec3b(width int, height int, data []byte) MatVec3b {
	cr := C.struct_RawData{
		width:  C.int(width),
		height: C.int(height),
		data:   toByteArray(data),
	}
	return MatVec3b{p: C.RawData_ToMatVec3b(cr)}
}

// MatVec4b is a bind of `cv::Mat_<cv::Vec4b>`
type MatVec4b struct {
	p C.MatVec4b
}

// Delete object.
func (m *MatVec4b) Delete() {
	C.MatVec4b_Delete(m.p)
	m.p = nil
}

// ToRawData converts MatVec4b to RawData.
func (m *MatVec4b) ToRawData() (int, int, []byte) {
	r := C.MatVec4b_ToRawData(m.p)
	return int(r.width), int(r.height), toGoBytes(r.data)
}

// ToMatVec4b converts RawData to MatVec4b. Returned MatVec4b is required to
// delete after using.
func ToMatVec4b(width int, height int, data []byte) MatVec4b {
	cr := C.struct_RawData{
		width:  C.int(width),
		height: C.int(height),
		data:   toByteArray(data),
	}
	return MatVec4b{p: C.RawData_ToMatVec4b(cr)}
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

// CascadeClassifier is a bind of `cv::CascadeClassifier`
type CascadeClassifier struct {
	p C.CascadeClassifier
}

// NewCascadeClassifier returns a new CascadeClassifier.
func NewCascadeClassifier() CascadeClassifier {
	return CascadeClassifier{p: C.CascadeClassifier_New()}
}

// Delete CascadeClassifier's pointer.
func (c *CascadeClassifier) Delete() {
	C.CascadeClassifier_Delete(c.p)
	c.p = nil
}

// Load cascade configuration file to classifier.
func (c *CascadeClassifier) Load(name string) bool {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.CascadeClassifier_Load(c.p, cName) != 0
}

// Rect represents rectangle. X and Y is a start point of Width and Height.
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// DetectMultiScale detects something which is decided by loaded file. Returns
// multi results addressed with rectangle.
func (c *CascadeClassifier) DetectMultiScale(img MatVec3b) []Rect {
	ret := C.CascadeClassifier_DetectMultiScale(c.p, img.p)
	defer C.Rects_Delete(ret)

	cArray := ret.rects
	length := int(ret.length)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cArray)),
		Len:  length,
		Cap:  length,
	}
	goSlice := *(*[]C.Rect)(unsafe.Pointer(&hdr))

	rects := make([]Rect, length)
	for i, r := range goSlice {
		rects[i] = Rect{
			X:      int(r.x),
			Y:      int(r.y),
			Width:  int(r.width),
			Height: int(r.height),
		}
	}
	return rects
}

// DrawRectsToImage draws rectangle information to target image.
func DrawRectsToImage(img MatVec3b, rects []Rect) {
	cRectArray := make([]C.struct_Rect, len(rects))
	for i, r := range rects {
		cRect := C.struct_Rect{
			x:      C.int(r.X),
			y:      C.int(r.Y),
			width:  C.int(r.Width),
			height: C.int(r.Height),
		}
		cRectArray[i] = cRect
	}
	cRects := C.struct_Rects{
		rects:  (*C.Rect)(&cRectArray[0]),
		length: C.int(len(rects)),
	}
	C.DrawRectsToImage(img.p, cRects)
}

// LoadAlphaImage loads RGBA type image.
func LoadAlphaImage(name string) MatVec4b {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return MatVec4b{p: C.LoadAlphaImg(cName)}
}

// MountAlphaImage draws img on back leading to rects. img is required RGBA,
// TODO should be check file type.
func MountAlphaImage(img MatVec4b, back MatVec3b, rects []Rect) {
	cRectArray := make([]C.struct_Rect, len(rects))
	for i, r := range rects {
		cRect := C.struct_Rect{
			x:      C.int(r.X),
			y:      C.int(r.Y),
			width:  C.int(r.Width),
			height: C.int(r.Height),
		}
		cRectArray[i] = cRect
	}
	cRects := C.struct_Rects{
		rects:  (*C.Rect)(&cRectArray[0]),
		length: C.int(len(rects)),
	}
	C.MountAlphaImage(img.p, back.p, cRects)
}
