package bridge

/*
#cgo pkg-config: scouter-core
#include "moving_matcher_bridge.h"
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type RegionsWithCameraID struct {
	CameraID   int
	Candidates []Candidate
}

type MVCandidate struct {
	p C.MVCandidate
}

func (c MVCandidate) Serialize() []byte {
	b := C.MVCandidate_Serialize(c.p)
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

func DeserializeMVCandiate(c []byte) MVCandidate {
	b := toByteArray(c)
	return MVCandidate{p: C.MVCandidate_Deserialize(b)}
}

func (c MVCandidate) Delete() {
	C.MVCandidate_Delete(c.p)
	c.p = nil
}

func convertCandidatezToPointer(regions []RegionsWithCameraID) []C.struct_RegionsWithCameraID {
	regionsPointers := []C.struct_RegionsWithCameraID{}
	for _, r := range regions {
		candidatePointers := convertCandidatesToPointer(r.Candidates) // -> []C.Candidate
		candidates := C.InvertCandidates((*C.Candidate)(&candidatePointers[0]),
			C.int(len(candidatePointers))) // -> C.struct_Candidates
		f := C.struct_RegionsWithCameraID{
			candidates: candidates,
			cameraID:   C.int(r.CameraID),
		}
		regionsPointers = append(regionsPointers, f)
	}
	return regionsPointers
}

func GetMatching(kThreashold float32, regions []RegionsWithCameraID) []MVCandidate {
	regionsPointers := convertCandidatezToPointer(regions) // -> []C.struct_RegionsWithCameraID
	mvCandidatePointers := C.MVOM_GetMatching((*C.struct_RegionsWithCameraID)(&regionsPointers[0]),
		C.int(len(regions)), C.float(kThreashold)) // -> vector<vector<ObjectCandidate>>
	defer C.MVCandidates_Delete(mvCandidatePointers)

	var cArray *C.MVCandidate = mvCandidatePointers.mvCandidates
	length := int(mvCandidatePointers.length)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cArray)),
		Len:  length,
		Cap:  length,
	}
	goSlice := *(*[]C.MVCandidate)(unsafe.Pointer(&hdr))

	ret := make([]MVCandidate, length)
	for i, c := range goSlice {
		ret[i] = MVCandidate{p: c}
	}
	return ret
}
