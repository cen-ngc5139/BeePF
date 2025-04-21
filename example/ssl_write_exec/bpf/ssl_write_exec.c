// go:build ignore

#include "vmlinux.h"
#include "vmlinux-x86.h"
#include "bpf/bpf_helpers.h"
#include "bpf/bpf_core_read.h"
#include "bpf/bpf_tracing.h"
#include "bpf/bpf_endian.h"
#include "bpf/bpf_ipv6.h"

char __license[] SEC("license") = "Dual MIT/GPL";

SEC("uprobe/SSL_write")
int ssl_write(struct pt_regs *ctx)
{
    bpf_printk("SSL_write: ctx=%p\n", ctx);
    // 使用 PT_REGS_PARM 宏获取参数，而不是直接定义函数参数
    void *ssl = (void *)PT_REGS_PARM1(ctx);
    const void *buf = (void *)PT_REGS_PARM2(ctx);
    size_t num = (size_t)PT_REGS_PARM3(ctx);

    bpf_printk("SSL_write_ex: ssl=%p, buf=%p, num=%zu\n", ssl, buf, num);
    return 0;
}
