#ifndef ICMP_H
#define ICMP_H

int send_icmp_echo(int s, uint32_t dst, uint16_t seq, uint8_t *data, size_t len);

#endif
