#include <stdint.h>
#include <stddef.h>

#include "utils.h"

uint16_t checksum(uint8_t *buf, size_t len) {
    unsigned int sum = 0;
    if (len%2 != 0) {
        sum += (buf[len-1]<<8);
    }
    
    for (int i = 0; i+1 < len; i += 2) {
        sum += (buf[i] << 8) | buf[i + 1];
    }

    while (sum >> 16) {
        sum = (sum&0xFFFF)+(sum >> 16);
    }

    return ~sum;
}
