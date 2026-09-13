#ifndef CONFIG_H
#define CONFIG_H

#define DEFAULT_ECHO_ID 2222

#define MESSAGE_ECHO_SEQ 1111
#define KEEPALIVE_ECHO_SEQ 2222

#define IPV4(a,b,c,d) ((uint32_t)(a)<<24|(b)<<16|(c)<<8|(d))

#define SERVER_COUNT 3

static const uint32_t servers[SERVER_COUNT] = {
    IPV4(8,8,8,8),
    IPV4(1,1,1,1),
    IPV4(127,0,0,1),
};

static const uint8_t *id = "client1";

#endif
