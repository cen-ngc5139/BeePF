// go:build ignore

#include "vmlinux.h"
#include "vmlinux-x86.h"
#include "bpf/bpf_helpers.h"
#include "bpf/bpf_core_read.h"
#include "bpf/bpf_tracing.h"
#include "bpf/bpf_endian.h"
#include "bpf/bpf_ipv6.h"

// 定义存储捕获数据的结构体
struct event
{
    u32 pid;
    u32 tgid;
    u64 timestamp;
    u32 pidfd;
    u32 fd;
    u32 flags;
    char comm[16];
};

// 定义一个性能事件映射用于向用户空间传输数据
struct
{
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(int));
    __uint(value_size, sizeof(int));
    __uint(max_entries, 1024);
} events SEC(".maps");

struct
{
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10240);
    __type(key, u32);
    __type(value, u32);
} pidfd_map SEC(".maps");

struct event *unused_event_t __attribute__((unused));

// kprobe 挂载点
SEC("kprobe/__x64_sys_pidfd_getfd")
int BPF_KPROBE(kprobe__pidfd_getfd)
{
    struct event e = {0};

    // 获取当前进程信息
    e.pid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
    e.tgid = bpf_get_current_pid_tgid() >> 32;
    e.timestamp = bpf_ktime_get_ns();
    bpf_get_current_comm(&e.comm, sizeof(e.comm));

    // 获取函数参数
    e.pidfd = (u32)PT_REGS_PARM1(ctx);
    e.fd = (u32)PT_REGS_PARM2(ctx);
    e.flags = (u32)PT_REGS_PARM3(ctx);

    // 将事件发送到用户空间
    bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &e, sizeof(e));
    bpf_map_update_elem(&pidfd_map, &e.pid, &e.fd, BPF_ANY);

    // 打印调试信息
    bpf_printk("pidfd_getfd: pid=%d, tgid=%d, timestamp=%lld, pidfd=%d, fd=%d, flags=%d, comm=%s\n",
               e.pid, e.tgid, e.timestamp, e.pidfd, e.fd, e.flags, e.comm);

    return 0;
}

SEC("kprobe/sys_mmap")
int BPF_KPROBE(kprobe__memset)
{
    u32 pid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
    u32 *fd = (u32 *)bpf_map_lookup_elem(&pidfd_map, &pid);
    if (!fd)
    {
        return 0;
    }
    bpf_printk("memset: pid=%d, fd=%d\n", pid, fd);

    return 0;
}

char LICENSE[] SEC("license") = "GPL";
