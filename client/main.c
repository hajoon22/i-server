#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <netinet/ip_icmp.h>
#include <stdlib.h>
#include <signal.h>
#include <sys/wait.h>

#include "config.h"
#include "icmp/icmp.h"
#include "utils/utils.h"

static void keepalive(int s) {
    for (int i = 0; i < SERVER_COUNT; i++) {
        send_icmp_echo(s, servers[i], KEEPALIVE_ECHO_SEQ, id, strlen((char *)id));
    }
}

static int init_keepalive(int s) {
    int pid = fork();
    if (pid == 0) {
        while (1) {
            keepalive(s);
            sleep(5);
        }
    }

    return pid;
}

int main(void) {
    int s = socket(AF_INET, SOCK_RAW, IPPROTO_ICMP);
    if (s < 0) return -1;

    int pid = init_keepalive(s);
    if (pid < 0) {
        close(s);
        return -1;
    }
    
    uint8_t buf[1500] = {0}, *data = NULL;
    while (1) {
        ssize_t n = recv(s, buf, sizeof(buf), 0);
        if (n < 0) break;

        struct iphdr *iph = (struct iphdr *)buf;
        struct icmphdr *icmph = (struct icmphdr *)(buf+(iph->ihl*4));
        
        int len = -1;
        if (icmph->type == ICMP_ECHOREPLY) {
            len = parse_icmp_echo(buf, n, MESSAGE_ECHO_SEQ, &data);
        } else if (icmph->type == ICMP_DEST_UNREACH) {
            len = parse_icmp_unreach(buf, n, MESSAGE_ECHO_SEQ, &data);
        }
        if (len < 0) continue;
            
        data[len] = '\0';
        printf("message = %s\r\n", data);
        
        free(data);
        data = NULL;
    }

    kill(pid, SIGTERM);
    waitpid(pid, NULL, 0);

    close(s);

    return 0;
}
