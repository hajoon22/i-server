#ifndef ICMP_H
#define ICMP_H

#include <stdint.h>

int send_icmp_echo(int s, uint32_t dst, uint16_t seq, uint8_t *data, size_t len);
int parse_icmp_unreach(uint8_t *buf, size_t len, uint8_t **output);

#endif
