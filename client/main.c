#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <poll.h>
#include <netinet/ip_icmp.h>

#include "config.h"
#include "icmp/icmp.h"
#include "utils/utils.h"

static void keepalive(int s) {
    for (int i = 0; i < SERVER_COUNT; i++) {
        send_icmp_echo(s, servers[i], KEEPALIVE_ECHO_SEQ, id, strlen((char *)id));
    }
}

int main(void) {
    int s = socket(AF_INET, SOCK_RAW, IPPROTO_ICMP);
    if (s < 0) return -1;

    struct pollfd pfd = {0};
    pfd.fd = s;
    pfd.events = POLLIN;
    
    int r = 0;
    uint8_t buf[1500] = {0}, *data = NULL;
    while (1) {
        r = poll(&pfd, 1, 5000);
        if (r > 0) {
            ssize_t n = recv(s, buf, sizeof(buf), 0);
            if (n < 0) break;

            struct iphdr *iph = (struct iphdr *)buf;
            struct icmphdr *icmph = (struct icmphdr *)(buf+(iph->ihl*4));
            if (icmph->type == ICMP_ECHOREPLY) {
                if (ntohs(icmph->un.echo.id) != DEFAULT_ECHO_ID) continue;
                
                // soon...
            } else if (icmph->type == ICMP_DEST_UNREACH) {
                int len = parse_icmp_unreach(buf, n, MESSAGE_ECHO_SEQ, &data);
                if (len < 0) continue;
                
                data[len] = '\0';
                printf("message = %s\r\n", data);
                
                free(data);
                data = NULL;
            }
        } else if (r < 0) {
            break;
        }

        keepalive(s);
    }

    close(s);
    return 0;
}
