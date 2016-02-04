#ifndef _OPENCV_BRIDGE_H_
#define _OPENCV_BRIDGE_H_

#include "util.h"

#ifdef __cplusplus
#include <opencv2/opencv.hpp>
extern "C" {
#endif

typedef struct RawData {
  int width;
  int height;
  struct ByteArray data;
} RawData;
typedef struct Rect {
  int x;
  int y;
  int width;
  int height;
} Rect;
typedef struct Rects {
  Rect* rects;
  int length;
} Rects;

#ifdef __cplusplus
typedef cv::Mat_<cv::Vec3b>* MatVec3b;
typedef cv::VideoCapture* VideoCapture;
typedef cv::VideoWriter* VideoWriter;
typedef cv::CascadeClassifier* CascadeClassifier;
#else
typedef void* MatVec3b;
typedef void* VideoCapture;
typedef void* VideoWriter;
typedef void* CascadeClassifier;
#endif

MatVec3b MatVec3b_New();
struct ByteArray MatVec3b_ToJpegData(MatVec3b m, int quality);
void MatVec3b_Delete(MatVec3b m);
void MatVec3b_CopyTo(MatVec3b src, MatVec3b dst);
int MatVec3b_Empty(MatVec3b m);
struct RawData MatVec3b_ToRawData(MatVec3b m);
MatVec3b RawData_ToMatVec3b(struct RawData r);

VideoCapture VideoCapture_New();
void VideoCapture_Delete(VideoCapture v);
int VideoCapture_Open(VideoCapture v, const char* uri);
int VideoCapture_OpenDevice(VideoCapture v, int device);
void VideoCapture_Release(VideoCapture v);
void VideoCapture_Set(VideoCapture v, int prop, int param);
int VideoCapture_IsOpened(VideoCapture v);
int VideoCapture_Read(VideoCapture v, MatVec3b buf);
void VideoCapture_Grab(VideoCapture v, int skip);

VideoWriter VideoWriter_New();
void VideoWriter_Delete(VideoWriter vw);
void VideoWriter_Open(VideoWriter vw, const char* name, double fps, int width,
  int height);
void VideoWriter_OpenWithMat(VideoWriter vw, const char* name, double fps,
  MatVec3b img);
int VideoWriter_IsOpened(VideoWriter vw);
void VideoWriter_Write(VideoWriter vw, MatVec3b img);

CascadeClassifier CascadeClassifier_New();
void CascadeClassifier_Delete(CascadeClassifier cs);
int CascadeClassifier_Load(CascadeClassifier cs, const char* name);
struct Rects CascadeClassifier_DetectMultiScale(CascadeClassifier cs, MatVec3b img);
void Rects_Delete(struct Rects rs);
void DrawRectsToImage(MatVec3b img, struct Rects rects);

#ifdef __cplusplus
}
#endif

#endif //_OPENCV_BRIDGE_H_
