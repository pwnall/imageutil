#include <stdint.h>

// Accelerates MaskRgba.
void GoMaskRgba(void* bytes, int byteCount, uint64_t mask) {
  uint64_t* words = (uint64_t*)bytes;
  for (int wordCount = (byteCount >> 3); wordCount > 0; --wordCount) {
    *words &= mask;
    ++words;
  }
}
