//go:build darwin
// +build darwin

package common

/*
#cgo CFLAGS: -x objective-c -mmacosx-version-min=10.15
#cgo LDFLAGS: -framework Vision -framework Foundation -framework CoreGraphics -framework CoreImage
#include <Vision/Vision.h>
#include <Foundation/Foundation.h>
#include <CoreGraphics/CoreGraphics.h>
#include <CoreImage/CoreImage.h>

// OCR识别函数
char* recognizeTextInImage(unsigned char* imageData, int dataLength) {
    @autoreleasepool {
        // 从字节数据创建 NSData
        NSData *data = [NSData dataWithBytes:imageData length:dataLength];
        if (data == nil) {
            return strdup("");
        }

        // 创建 CIImage
        CIImage *ciImage = [CIImage imageWithData:data];
        if (ciImage == nil) {
            return strdup("");
        }

        // 创建 VNImageRequestHandler
        VNImageRequestHandler *handler = [[VNImageRequestHandler alloc] initWithCIImage:ciImage options:@{}];
        if (handler == nil) {
            return strdup("");
        }

        // 创建文字识别请求
        VNRecognizeTextRequest *request = [[VNRecognizeTextRequest alloc] init];

        // 设置识别级别（准确度优先）
        request.recognitionLevel = VNRequestTextRecognitionLevelAccurate;

        // 启用语言校正（提高准确性）
        request.usesLanguageCorrection = YES;

        // 设置识别语言（支持中文和英文）
        request.recognitionLanguages = @[@"zh-Hans", @"zh-Hant", @"en-US"];

        // 执行请求（同步执行）
        NSError *error = nil;
        BOOL success = [handler performRequests:@[request] error:&error];

        if (!success || error != nil) {
            return strdup("");
        }

        // 提取识别结果
        NSMutableString *result = [NSMutableString string];
        NSArray<VNRecognizedTextObservation*> *observations = request.results;

        if (observations.count == 0) {
            return strdup("");
        }

        // 改进的排序算法：使用容差判断同一行，然后按列排序
        // 1. 计算平均文本块高度（用于判断容差）
        CGFloat totalHeight = 0;
        for (VNRecognizedTextObservation *obs in observations) {
            totalHeight += obs.boundingBox.size.height;
        }
        CGFloat avgHeight = totalHeight / observations.count;
        CGFloat rowTolerance = avgHeight * 0.4; // 行容差：平均高度的 40%

        // 2. 改进的排序：先按行分组（使用容差），再在行内按列排序
        NSArray<VNRecognizedTextObservation*> *sortedObservations = [observations sortedArrayUsingComparator:^NSComparisonResult(VNRecognizedTextObservation *obs1, VNRecognizedTextObservation *obs2) {
            CGRect rect1 = obs1.boundingBox;
            CGRect rect2 = obs2.boundingBox;

            // 计算文本块的中心 Y 坐标
            CGFloat centerY1 = rect1.origin.y + rect1.size.height / 2.0;
            CGFloat centerY2 = rect2.origin.y + rect2.size.height / 2.0;

            // 计算 Y 坐标差值
            CGFloat yDiff = fabs(centerY1 - centerY2);

            // 如果 Y 坐标差值小于容差，认为在同一行，按 X 坐标排序
            if (yDiff <= rowTolerance) {
                // 同一行，比较 X 坐标（从左到右）
                CGFloat centerX1 = rect1.origin.x + rect1.size.width / 2.0;
                CGFloat centerX2 = rect2.origin.x + rect2.size.width / 2.0;

                if (centerX1 < centerX2) {
                    return NSOrderedAscending;  // X 值小的排在前面（左）
                } else if (centerX1 > centerX2) {
                    return NSOrderedDescending; // X 值大的排在后面（右）
                }
                return NSOrderedSame;
            } else {
                // 不同行，按 Y 坐标排序（从上到下）
                // 由于原点在左下角，Y 值越大位置越靠上
                if (centerY1 > centerY2) {
                    return NSOrderedAscending;  // Y 值大的排在前面（上）
                } else {
                    return NSOrderedDescending; // Y 值小的排在后面（下）
                }
            }
        }];

        // 3. 输出结果，同一行的文本用空格分隔，不同行用换行符
        CGFloat lastY = -1;
        for (VNRecognizedTextObservation *observation in sortedObservations) {
            NSArray<VNRecognizedText*> *topCandidates = [observation topCandidates:1];
            if (topCandidates.count > 0) {
                VNRecognizedText *recognizedText = topCandidates[0];

                // 可选：过滤低置信度结果（置信度阈值 0.5）
                // if (recognizedText.confidence < 0.5) {
                //     continue;
                // }

                CGRect rect = observation.boundingBox;
                CGFloat currentY = rect.origin.y + rect.size.height / 2.0;

                // 判断是否需要换行或添加空格
                if (result.length > 0) {
                    if (lastY >= 0 && fabs(currentY - lastY) > rowTolerance) {
                        // 不同行，添加换行符
                        [result appendString:@"\n"];
                    } else if (lastY >= 0) {
                        // 同一行，添加空格分隔
                        [result appendString:@" "];
                    }
                }

                [result appendString:recognizedText.string];
                lastY = currentY;
            }
        }

        // 返回结果
        const char *cString = [result UTF8String];
        if (cString == nil) {
            return strdup("");
        }
        return strdup(cString);
    }
}
*/
import "C"
import (
	"unsafe"
)

// RecognizeTextInImage 识别图片中的文字（macOS Vision 框架）
func RecognizeTextInImage(imageData []byte) string {
	if len(imageData) == 0 {
		return ""
	}

	cstr := C.recognizeTextInImage((*C.uchar)(unsafe.Pointer(&imageData[0])), C.int(len(imageData)))
	defer C.free(unsafe.Pointer(cstr))

	return C.GoString(cstr)
}
