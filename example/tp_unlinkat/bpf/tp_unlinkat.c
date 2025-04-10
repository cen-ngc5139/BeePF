// go:build ignore

#include "vmlinux.h"
#include "vmlinux-x86.h"
#include "bpf/bpf_helpers.h"
#include "bpf/bpf_core_read.h"
#include "bpf/bpf_tracing.h"

char __license[] SEC("license") = "Dual MIT/GPL";

struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, u32);
    __type(value, u64);
} kprobe_map SEC(".maps");

// tracepoint/syscalls/sys_enter_unlinkat 格式定义
struct sys_enter_unlinkat_args {
    unsigned short common_type;
    unsigned char common_flags;
    unsigned char common_preempt_count;
    int common_pid;
    
    int __syscall_nr;
    unsigned long dfd;
    const char *pathname;  // 文件路径指针
    unsigned long flag;
};

SEC("tracepoint/syscalls/sys_enter_unlinkat")
int handle_unlinkat(struct sys_enter_unlinkat_args *args) {
    char dummy;
    long ret = bpf_probe_read_user(&dummy, 1, args->pathname);
    if (ret < 0)
    {
        bpf_printk("Read userspace address false:%d\n", ret); //返回-14，查看资料说的是无法安全访问用户空间地址
        return 0;
    }
    
    char filename[64];
    ret = bpf_probe_read_user_str(filename, sizeof(filename), (const char*)args->pathname);
    if (ret < 0)
    {
        bpf_printk("bpf_probe_read_user_str failed:%d\n",ret);    //返回-14，查看资料说的是无法安全访问用户空间地址
        return 0;
    }

    bpf_printk("filename: %s\n", filename);
    
    return 0;
}

