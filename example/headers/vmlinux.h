#if defined(__TARGET_ARCH_x86)
#include "vmlinux-x86.h"
#elif defined(__TARGET_ARCH_arm64)
#include "vmlinux-arm64.h"
#else
#error "Unknown architecture"
#endif

#define MSG_OOB 0x1       /* process out-of-band data */
#define MSG_PEEK 0x2      /* peek at incoming message */
#define MSG_DONTROUTE 0x4 /* send without using routing tables */
#define MSG_EOR 0x8       /* data completes record */
#define MSG_TRUNC 0x10    /* data discarded before delivery */
#define MSG_CTRUNC 0x20   /* control data lost before delivery */
#define MSG_WAITALL 0x40  /* wait for full request or error */
#if !defined(_POSIX_C_SOURCE) || defined(_DARWIN_C_SOURCE)
#define MSG_DONTWAIT 0x80 /* this message should be nonblocking */
#define MSG_EOF 0x100     /* data completes connection */
#ifdef __APPLE__
#ifdef __APPLE_API_OBSOLETE
#define MSG_WAITSTREAM 0x200 /* wait up to full request.. may return partial */
#endif
#define MSG_FLUSH 0x400     /* Start of 'hold' seq; dump so_temp, deprecated */
#define MSG_HOLD 0x800      /* Hold frag in so_temp, deprecated */
#define MSG_SEND 0x1000     /* Send the packet in so_temp, deprecated */
#define MSG_HAVEMORE 0x2000 /* Data ready to be read */
#define MSG_RCVMORE 0x4000  /* Data remains in current pkt */
#endif
#define MSG_NEEDSA 0x10000 /* Fail receive if socket address cannot be allocated */
#endif                     /* (!_POSIX_C_SOURCE || _DARWIN_C_SOURCE) */