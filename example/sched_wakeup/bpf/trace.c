// SPDX-License-Identifier: (GPL-2.0-only OR BSD-2-Clause)
/* Copyright Martynas Pumputis */
/* Copyright Authors of Cilium */

#include "vmlinux.h"
#include "vmlinux-x86.h"
#include "bpf/bpf_helpers.h"
#include "bpf/bpf_core_read.h"
#include "bpf/bpf_tracing.h"
#include "bpf/bpf_endian.h"
#include "bpf/bpf_ipv6.h"

extern int LINUX_KERNEL_VERSION __kconfig;

// 定义数据结构来存储调度延迟信息
struct sched_latency_t
{
    __u32 pid;               // 进程ID
    __u32 tid;               // 线程ID
    __u64 delay_ns;          // 调度延迟(纳秒)
    __u64 ts;                // 时间戳
    __u32 preempted_pid;     // 被抢占的进程ID
    char preempted_comm[16]; // 被抢占的进程名
    __u64 is_preempt;        // 是否抢占(0: 否, 1: 是)
    char comm[16];           // 进程名
} __attribute__((packed));

struct sched_latency_t *unused_sched_latency_t __attribute__((unused));

// 定义 ring buffer 用于传输数据到用户空间
struct
{
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(max_entries, 256 * 1024);
} sched_events SEC(".maps");

// 用于临时存储唤醒时间的 hash map
struct
{
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10240);
    __type(key, __u32);
    __type(value, __u64);
} wakeup_times SEC(".maps");

struct trace_event_raw_sched_wakeup
{
    /* common fields */
    __u16 common_type;         /* offset: 0, size: 2 */
    __u8 common_flags;         /* offset: 2, size: 1 */
    __u8 common_preempt_count; /* offset: 3, size: 1 */
    __s32 common_pid;          /* offset: 4, size: 4 */

    /* event specific fields */
    char comm[16];         /* offset: 8, size: 16 */
    __s32 pid;             /* offset: 24, size: 4 */
    __s32 prio;            /* offset: 28, size: 4 */
    __s32 target_cpu;      /* offset: 32, size: 4 */
} __attribute__((packed)); /* 确保结构体紧凑，没有额外的填充字节 */

// 添加 tracepoint 事件结构体定义
struct trace_event_raw_sched_wakeup_new
{
    /* common fields */
    __u16 common_type;
    __u8 common_flags;
    __u8 common_preempt_count;
    __s32 common_pid;

    /* event specific fields */
    char comm[16];
    __s32 pid;
    __s32 prio;
    __s32 target_cpu;
} __attribute__((packed));

// 公共函数：处理进程唤醒
static __always_inline void handle_wakeup(u32 pid)
{
    if (pid == 0)
    {
        return;
    }

    __u64 ts = bpf_ktime_get_ns();
    bpf_map_update_elem(&wakeup_times, &pid, &ts, BPF_ANY);
}

#if LINUX_KERNEL_VERSION >= KERNEL_VERSION(5, 10, 0)
SEC("tp_btf/sched_wakeup")
int sched_wakeup(u64 *ctx)
{
    struct task_struct *task = (void *)ctx[0];
    handle_wakeup(task->pid);
    return 0;
}

SEC("tp_btf/sched_wakeup_new")
int sched_wakeup_new(u64 *ctx)
{
    struct task_struct *task = (void *)ctx[0];
    handle_wakeup(task->pid);
    return 0;
}
#else
SEC("tp/sched/sched_wakeup")
int sched_wakeup(struct trace_event_raw_sched_wakeup *ctx)
{
    handle_wakeup(ctx->pid);
    return 0;
}

SEC("tp/sched/sched_wakeup_new")
int sched_wakeup_new(struct trace_event_raw_sched_wakeup_new *ctx)
{
    handle_wakeup(ctx->pid);
    return 0;
}
#endif

// 定义流控相关的常量和map
#define SAMPLING_RATIO 100   // 采样率 1/100
#define THRESHOLD_NS 1000000 // 延迟阈值 1ms

// 用于记录每个 CPU 的最后一次采样时间
struct
{
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __uint(max_entries, 1);
    __type(key, u32);
    __type(value, u64);
} last_sample SEC(".maps");

#define TASK_RUNNING 0

u64 get_task_cgroup_id(struct task_struct *task)
{
    u64 cgroup_id = 0;
    struct css_set *cgroups;

    // 使用 BPF_CORE_READ 安全地读取
    cgroups = BPF_CORE_READ(task, cgroups);
    if (cgroups)
    {
        cgroup_id = BPF_CORE_READ(cgroups, dfl_cgrp, kn, id);
    }

    return cgroup_id;
}

// 公共函数：处理调度切换事件
static __always_inline void handle_sched_switch(u32 prev_pid, u32 prev_tgid,
                                                u32 next_pid, u32 next_tgid, __u32 prev_state,
                                                const char *prev_comm, const char *next_comm, void *ctx)
{
    __u64 *wakeup_ts;
    __u64 now = bpf_ktime_get_ns();

    if (prev_pid == 0 || next_pid == 0)
    {
        return;
    }

    // 查找进程的唤醒时间
    wakeup_ts = bpf_map_lookup_elem(&wakeup_times, &next_pid);
    if (!wakeup_ts)
        return;

    // 计算调度延迟
    __u64 delay = now - *wakeup_ts;

    // 流控逻辑开始
    __u32 key = 0;
    __u64 *last_ts = bpf_map_lookup_elem(&last_sample, &key);
    if (!last_ts)
        return;

    // 基于时间的流控
    if ((now - *last_ts) < THRESHOLD_NS)
    {
        if (bpf_get_prandom_u32() % SAMPLING_RATIO != 0)
        {
            bpf_map_delete_elem(&wakeup_times, &next_pid);
            return;
        }
    }

    // 更新最后采样时间
    bpf_map_update_elem(&last_sample, &key, &now, BPF_ANY);

    // 延迟阈值过滤
    if (delay < THRESHOLD_NS)
    {
        bpf_map_delete_elem(&wakeup_times, &next_pid);
        return;
    }

    // 准备输出数据
    struct sched_latency_t latency = {
        .pid = next_tgid ? next_tgid : next_pid, // 如果有 tgid 则使用 tgid，否则使用 pid
        .tid = next_pid,
        .delay_ns = delay,
        .ts = now,
    };

    bpf_probe_read_kernel_str(&latency.comm, sizeof(latency.comm), next_comm);

    // 如果前一个状态是 TASK_RUNNING，则认为是抢占
    if (prev_state == TASK_RUNNING)
    {
        latency.is_preempt = 1;
        latency.preempted_pid = prev_tgid ? prev_tgid : prev_pid;
        bpf_probe_read_kernel_str(&latency.preempted_comm, sizeof(latency.preempted_comm), prev_comm);
    }

    bpf_printk("pid: %d, delay: %llu ns, is_preempt: %d\n",
               latency.pid, latency.delay_ns, latency.is_preempt);

    // 输出到 perf event
    bpf_perf_event_output(ctx, &sched_events, BPF_F_CURRENT_CPU, &latency, sizeof(latency));

    // 删除已处理的唤醒时间记录
    bpf_map_delete_elem(&wakeup_times, &next_pid);
}

#if LINUX_KERNEL_VERSION >= KERNEL_VERSION(5, 10, 0)
SEC("tp_btf/sched_switch")
int sched_switch(u64 *ctx)
{
    struct task_struct *prev = (struct task_struct *)ctx[1];
    struct task_struct *next = (struct task_struct *)ctx[2];

    u32 prev_pid = BPF_CORE_READ(prev, pid);
    u32 prev_tgid = BPF_CORE_READ(prev, tgid);
    u32 next_pid = BPF_CORE_READ(next, pid);
    u32 next_tgid = BPF_CORE_READ(next, tgid);

#if defined(__TARGET_ARCH_x86)
    __u32 state = BPF_CORE_READ(prev, __state);
#elif defined(__TARGET_ARCH_arm64)
    __u32 state = BPF_CORE_READ(prev, state);
#else
    __u32 state = BPF_CORE_READ(prev, __state);
#endif

    handle_sched_switch(prev_pid, prev_tgid, next_pid, next_tgid,
                        state, prev->comm, next->comm, ctx);
    return 0;
}
#else
SEC("tp/sched/sched_switch")
int sched_switch(struct trace_event_raw_sched_switch *ctx)
{
    handle_sched_switch(ctx->prev_pid, 0, ctx->next_pid, 0,
                        ctx->prev_state, ctx->prev_comm, ctx->next_comm, ctx);
    return 0;
}
#endif

char __license[] SEC("license") = "Dual BSD/GPL";
