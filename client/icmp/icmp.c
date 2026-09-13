#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <netinet/ip_icmp.h>

#include "../config.h"
#include "../utils/utils.h"

int send_icmp_echo(int s, uint32_t dst, uint16_t seq, const uint8_t *data, size_t len) {
    if (s < 0 || (len > 0 && data == NULL)) {
        return -1;
    }

    size_t total = sizeof(struct icmphdr)+len;
    uint8_t *buf = calloc(total, sizeof(uint8_t));
    if (!buf) return -1;

    struct icmphdr *icmph = (struct icmphdr *)buf;
    icmph->type = ICMP_ECHO;
    icmph->un.echo.id = htons(DEFAULT_ECHO_ID);
    icmph->un.echo.sequence = htons(seq);

    memcpy(buf+sizeof(struct icmphdr), data, len);
    icmph->checksum = htons(checksum(buf, total));

    
    struct sockaddr_in sin;
    sin.sin_family = AF_INET;
    sin.sin_addr.s_addr = htonl(dst);

    ssize_t ret = sendto(s, buf, total, 0, (struct sockaddr *)&sin, sizeof(sin));
    free(buf);
    return (int)ret;
}

int parse_icmp_echo(uint8_t *buf, size_t len, uint16_t seq, uint8_t **output) {
    size_t offset = 0;
    struct iphdr *iph = (struct iphdr *)buf;
    offset += iph->ihl*4;

    struct icmphdr *icmph = (struct icmphdr *)(buf+offset);
    if (ntohs(icmph->un.echo.id) != DEFAULT_ECHO_ID) {
        return -1;
    } else if (ntohs(icmph->un.echo.sequence) != seq) {
        return -1;
    }
    offset += sizeof(struct icmphdr);

    size_t data_len = len-offset;
    if (data_len == 0) return -1;

    *output = calloc(data_len+1, sizeof(uint8_t));
    if (!*output) return -1;

    memcpy(*output, buf+offset, data_len);
    return (int)data_len;
}

int parse_icmp_unreach(uint8_t *buf, size_t len, uint16_t seq, uint8_t **output) {
    size_t offset = 0;

    // outer packet
    struct iphdr *iph = (struct iphdr *)buf;
    offset += iph->ihl*4;

    struct icmphdr *icmph = (struct icmphdr *)(buf+offset);
    if (icmph->type != ICMP_DEST_UNREACH) return -1;
    offset += sizeof(struct icmphdr);

    // inner packet's header size
    int header_size = sizeof(struct iphdr)+sizeof(struct icmphdr);
    if (len < offset+header_size) {
        return -1;
    }

    // inner packet
    return parse_icmp_echo(buf+offset, len-offset, seq, output);
}

