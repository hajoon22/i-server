#include <stdint.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <poll.h>

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
    while (1) {
        r = poll(&pfd, 1, 5000);
        if (r > 0) {
            
        } else if (r < 0) {
            break;
        }

        keepalive(s);
    }

    close(s);
    return 0;
}
