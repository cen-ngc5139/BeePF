// go:build ignore

#include "vmlinux.h"
#include "vmlinux-x86.h"
#include "bpf/bpf_helpers.h"

char __license[] SEC("license") = "Dual MIT/GPL";

struct pkt_count
{
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, u32);
    __type(value, u64);
} pkt_count SEC(".maps");

struct span_info
{
    __u64 timestamp;      // 操作时间
    __u64 start_time;     // 操作开始时间
    __u64 end_time;       // 操作结束时间
    __u64 trace_id;       // trace id 高 64 位
    __u64 span_id;        // span id
    __u64 parent_span_id; // 父 span id
    __u32 pid;            // 进程 ID
    __u32 tid;            // 线程 ID
    char name[32];        // span 名称
    char pod[100];        // pod 名称
    char container[100];  // 容器名称
    __u64 dev_file_id;    // 设备文件 ID
} __attribute__((packed));

struct
{
    __uint(type, BPF_MAP_TYPE_HASH);
    __type(key, u32);
    __type(value, struct span_info);
    __uint(max_entries, 1024);
} link_begin SEC(".maps");

// 需要通过 bpf2go 将该结构体转换为 Go 结构体，同时在 ebpf.Collection.Variables 中添加该结构体
struct span_info *unused_span_info __attribute__((unused));

SEC("cgroup_skb/egress")
int count_egress_packets(struct __sk_buff *skb)
{
    u32 key = 0;
    u64 init_val = 1;

    u64 *count = bpf_map_lookup_elem(&pkt_count, &key);
    if (!count)
    {
        bpf_map_update_elem(&pkt_count, &key, &init_val, BPF_ANY);
        return 1;
    }
    __sync_fetch_and_add(count, 1);

    return 1;
}