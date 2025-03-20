#include <vmlinux.h>
#include <vmlinux-x86.h>
#include "bpf/bpf_helpers.h"
#include "bpf/bpf_core_read.h"
#include "bpf/bpf_tracing.h"
#include "bpf/bpf_endian.h"
#include "bpf/bpf_ipv6.h"

char LICENSE[] SEC("license") = "GPL";

// 事件类型定义
#define EVENT_TYPE_ADD 1    // 程序加载
#define EVENT_TYPE_UPDATE 2 // 程序更新
#define EVENT_TYPE_DELETE 3 // 程序卸载

// 程序状态结构
struct bpf_prog_state
{
    __u32 prog_id;   // 程序ID
    __u64 load_time; // 加载时间
    char comm[16];   // 加载程序的进程名
    __u32 pid;       // 加载程序的进程ID
};

struct bpf_prog_state *unused_bpf_prog_state __attribute__((unused));

// 发送到 ringbuffer 的事件结构
struct bpf_prog_event
{
    __u32 event_type;            // 事件类型(ADD/UPDATE/DELETE)
    struct bpf_prog_state state; // 程序状态
};

// 定义 ringbuffer map
struct
{
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24); // 16MB
} events SEC(".maps");

// 定义程序状态跟踪 hash map (prog_id -> state)
struct
{
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10000);
    __type(key, __u32);
    __type(value, struct bpf_prog_state);
} prog_states SEC(".maps");

// 定义复合键结构
struct pid_func_key
{
    __u32 pid;
    char func_name[16]; // 函数名最大长度
};

// 定义 pid + funcname 的映射
struct
{
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10000);
    __type(key, struct pid_func_key);
    __type(value, struct bpf_prog_state);
} pid_func_name_states SEC(".maps");

// 辅助函数: 发送事件到 ringbuffer
static int send_event(__u32 event_type, struct bpf_prog_state *state)
{
    struct bpf_prog_event *event;

    // 从 ringbuffer 分配内存
    event = bpf_ringbuf_reserve(&events, sizeof(*event), 0);
    if (!event)
    {
        return -1;
    }

    // 填充事件信息
    event->event_type = event_type;
    __builtin_memcpy(&event->state, state, sizeof(*state));

    // 提交事件
    bpf_ringbuf_submit(event, 0);
    return 0;
}

// kprobe 处理函数：监控程序加载
SEC("kprobe/bpf_prog_kallsyms_add")
int BPF_KPROBE(trace_bpf_prog_kallsyms_add)
{
    // 获取返回的 bpf_prog 指针
    struct bpf_prog *prog = (struct bpf_prog *)PT_REGS_PARM1(ctx);
    if (!prog)
        return 0;

    // 准备程序状态结构
    struct bpf_prog_state state = {0};

    // 获取当前进程信息
    state.pid = bpf_get_current_pid_tgid() >> 32;
    bpf_get_current_comm(&state.comm, sizeof(state.comm));
    state.load_time = bpf_ktime_get_ns();

    // 获取 aux 指针以读取 id
    struct bpf_prog_aux *aux;
    bpf_probe_read_kernel(&aux, sizeof(aux), &prog->aux);
    if (!aux)
        return 0;

    // 读取程序 ID
    bpf_probe_read_kernel(&state.prog_id, sizeof(state.prog_id), &aux->id);
    if (state.prog_id == 0)
        return 0;

    char func_name[16];
    bpf_probe_read_kernel(&func_name, sizeof(func_name), &aux->name);
    if (func_name == NULL)
        return 0;

    struct pid_func_key key = {0};
    key.pid = state.pid;
    __builtin_memcpy(key.func_name, func_name, sizeof(func_name));

    // 检查是否存在，确定是ADD还是UPDATE
    struct bpf_prog_state *existing;
    existing = bpf_map_lookup_elem(&pid_func_name_states, &key);

    __u32 event_type;
    if (existing)
    {
        event_type = EVENT_TYPE_UPDATE;
    }
    else
    {
        event_type = EVENT_TYPE_ADD;
    }

    // 更新程序状态map
    bpf_map_update_elem(&prog_states, &state.prog_id, &state, BPF_ANY);
    bpf_map_update_elem(&pid_func_name_states, &key, &state, BPF_ANY);

    // 发送事件到ringbuffer
    send_event(event_type, &state);

    // 打印基本信息
    bpf_printk("BPF program %s: id=%u pid=%d comm=%s func_name=%s\n",
               event_type == EVENT_TYPE_ADD ? "loaded" : "updated",
               state.prog_id, state.pid, state.comm, func_name);

    return 0;
}

// kprobe 处理函数：监控程序释放
SEC("kprobe/bpf_prog_kallsyms_del_all")
int BPF_KPROBE(trace_bpf_prog_put)
{
    __u32 pid = bpf_get_current_pid_tgid() >> 32;
    char comm[16];
    bpf_get_current_comm(&comm, sizeof(comm));

    struct bpf_prog *prog;
    prog = (struct bpf_prog *)PT_REGS_PARM1(ctx);
    if (!prog)
    {
        bpf_printk("fail to get bpf_prog: pid=%u comm=%s\n", pid, comm);
        return 0;
    }

    // 读取程序信息
    struct bpf_prog_aux *aux;
    bpf_probe_read_kernel(&aux, sizeof(aux), &prog->aux);
    if (!aux)
    {
        bpf_printk("fail to get bpf_prog_aux: pid=%u comm=%s\n", pid, comm);
        return 0;
    }

    char func_name[16];
    bpf_probe_read_kernel(&func_name, sizeof(func_name), &aux->name);

    // 尝试从状态map中查找
    struct bpf_prog_state *state;
    struct pid_func_key key = {0};
    key.pid = pid;
    __builtin_memcpy(key.func_name, func_name, sizeof(func_name));
    state = bpf_map_lookup_elem(&pid_func_name_states, &key);
    if (!state)
    {
        bpf_printk("BPF program deleted (unknown load): id=%u\n", state->prog_id);
        return 0;
    }

    // 读取当前引用计数
    atomic_t ref_cnt;
    bpf_probe_read_kernel(&ref_cnt, sizeof(ref_cnt), &aux->refcnt);
    int count = 0;
    bpf_probe_read_kernel(&count, sizeof(count), &ref_cnt);

    // 更新状态中的引用计数
    struct bpf_prog_state updated_state = *state;

    // 检查是否为最后一个引用（释放）
    if (count <= 0)
    {
        // 发送删除事件
        send_event(EVENT_TYPE_DELETE, &updated_state);

        // 从map中删除
        bpf_map_delete_elem(&prog_states, &state->prog_id);
        bpf_map_delete_elem(&pid_func_name_states, &key);

        bpf_printk("BPF program deleted: id=%u\n",
                   state->prog_id);
    }

    return 0;
}
