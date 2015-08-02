package bridge

/*
#cgo pkg-config: scouter-core
#include <stdlib.h>
#include "instances_visualizer_bridge.h"
*/
import "C"
import (
	"unsafe"
)

type InstancesVisualizer struct {
	p C.InstancesVisualizer
}

// NewInstancesVisualizer create InstancesVisualizer. InstanceManager is
// necessary to create a visualizer, but not use the manager in drawing.
func NewInstancesVisualizer(im *InstanceManager, config string) InstancesVisualizer {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return InstancesVisualizer{p: C.InstancesVisualizer_New(im.p, cConfig)}
}

func (v *InstancesVisualizer) Delete() {
	C.InstancesVisualizer_Delete(v.p)
	v.p = nil
}

func (v *InstancesVisualizer) Draw(frames map[int]MatVec3b, states []InstanceState,
	trackees []Trackee) MatVec3b {

	// MatMapPtr
	fLength := len(frames)
	framesPtr := []C.struct_MatWithCameraID{}
	for k, v := range frames {
		matWithID := C.struct_MatWithCameraID{
			cameraID: C.int(k),
			mat:      v.p,
		}
		framesPtr = append(framesPtr, matWithID)
	}

	// C.InstanceStates
	ss := []C.InstanceState{}
	for _, is := range states {
		ss = append(ss, is.p)
	}
	iss := C.struct_InstanceStates{
		instanceStates: (*C.InstanceState)(&ss[0]),
		length:         C.int(len(ss)),
	}

	// *C.Trackee
	tLength := len(trackees)
	trs := []C.struct_Trackee{}
	for _, t := range trackees {
		var interpo int
		if t.Interpolated {
			interpo = 1
		} else {
			interpo = 0
		}
		trackee := C.struct_Trackee{
			colorID:      C.ulonglong(t.ColorID),
			mvCandidate:  t.MVCandidate.p,
			interpolated: C.int(interpo),
		}

		trs = append(trs, trackee)
	}

	img := C.InstancesVisualizer_Draw(v.p, (*C.MatWithCameraID)(&framesPtr[0]),
		C.int(fLength), iss, (*C.Trackee)(&trs[0]), C.int(tLength))
	return MatVec3b{p: img}
}
