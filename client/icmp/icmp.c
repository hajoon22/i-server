#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/ip_icmp.h>

#include "../config.h"
#include "../utils/utils.h"

int send_icmp_echo(int s, uint32_t dst, uint16_t seq, uint8_t *data, size_t len) {
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
