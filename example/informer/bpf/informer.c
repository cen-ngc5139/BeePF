#include <vmlinux.h>
#include <vmlinux-x86.h>
#include "bpf/bpf_helpers.h"
#include "bpf/bpf_core_read.h"
#include "bpf/bpf_tracing.h"
#include "bpf/bpf_endian.h"
#include "bpf/bpf_ipv6.h"

char LICENSE[] SEC("license") = "GPL";

// 定义 BPF 命令常量
#define BPF_PROG_LOAD 5
#define BPF_MAP_CREATE 0
#define BPF_PROG_ATTACH 8
#define BPF_PROG_DETACH 9
#define BPF_LINK_CREATE 18
#define BPF_LINK_DETACH 26
#define BPF_OBJ_PIN 6
#define BPF_OBJ_GET 7

// 对象类型常量
#define OBJ_TYPE_PROG 1
#define OBJ_TYPE_MAP 2
#define OBJ_TYPE_LINK 3

// 新的结构体，用于存储 BPF 对象信息（包括 prog 和 map）
struct bpf_obj_info_t
{
    __u32 obj_type;  // 1: prog, 2: map, 3: link 等
    __u32 pid;       // 进程 ID
    __u32 cmd;       // BPF 命令
    __s32 fd;        // 文件描述符
    __u32 id;        // 对象 ID
    __u64 timestamp; // 时间戳
    char comm[16];   // 进程名
};

// 新增: 使用组合键 (pid + command) 来存储进行中的 BPF 操作
struct
{
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10000);
    __type(key, __u64); // pid 和 cmd 的组合键
    __type(value, struct bpf_obj_info_t);
} pending_bpf_ops SEC(".maps");

// 新增: 使用 fd 来跟踪 BPF 对象
struct
{
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10000);
    __type(key, __u64); // pid 和 fd 的组合键
    __type(value, struct bpf_obj_info_t);
} fd_bpf_map SEC(".maps");

// 修改后的结构体，匹配 tracepoint 格式
struct syscalls_enter_bpf_args
{
    __u16 common_type;
    __u8 common_flags;
    __u8 common_preempt_count;
    __s32 common_pid;

    __s32 __syscall_nr;
    __u64 cmd;   // 注意：size 从 4 改为 8
    __u64 uattr; // 用户空间指针，改名为 uattr
    __u64 size;  // 注意：size 从 4 改为 8
};

// 在 BPF 系统调用出口处记录 fd
struct trace_event_sys_exit
{
    __u16 common_type;
    __u8 common_flags;
    __u8 common_preempt_count;
    __s32 common_pid;

    __s32 __syscall_nr;
    __s64 ret;
};

// 创建组合键函数
static inline __u64 make_key(__u32 pid, __u32 cmd_or_fd)
{
    return ((__u64)pid << 32) | cmd_or_fd;
}

SEC("tracepoint/syscalls/sys_enter_bpf")
int trace_bpf_syscall(struct syscalls_enter_bpf_args *ctx)
{
    __u32 pid = bpf_get_current_pid_tgid() >> 32;
    __u32 cmd = ctx->cmd;
    __u32 obj_type = 0;

    // 只处理我们关心的命令
    switch (cmd)
    {
    case BPF_PROG_LOAD:
        obj_type = OBJ_TYPE_PROG;
        bpf_printk("BPF Program Load by PID: %d\n", pid);
        break;
    case BPF_MAP_CREATE:
        obj_type = OBJ_TYPE_MAP;
        bpf_printk("BPF Map Create by PID: %d\n", pid);
        break;
    case BPF_LINK_CREATE:
        obj_type = OBJ_TYPE_LINK;
        bpf_printk("BPF Link Create by PID: %d\n", pid);
        break;
    default:
        // 其他命令我们不关心
        return 0;
    }

    // 如果是我们关心的命令，创建一个对象信息记录
    if (obj_type > 0)
    {
        struct bpf_obj_info_t info = {};

        info.obj_type = obj_type;
        info.pid = pid;
        info.cmd = cmd;
        info.fd = -1; // fd 将在 sys_exit_bpf 中设置
        info.timestamp = bpf_ktime_get_ns();
        bpf_get_current_comm(&info.comm, sizeof(info.comm));

        // 使用 pid 和 cmd 的组合键存储到 pending_bpf_ops 映射中
        __u64 key = make_key(pid, cmd);
        bpf_map_update_elem(&pending_bpf_ops, &key, &info, BPF_ANY);
    }

    return 0;
}

SEC("tracepoint/syscalls/sys_exit_bpf")
int trace_bpf_syscall_ret(struct trace_event_sys_exit *ctx)
{
    __s64 ret = ctx->ret;
    __u32 pid = bpf_get_current_pid_tgid() >> 32;

    // 只有当返回成功的 fd 时才处理
    if (ret > 0)
    {
        __s32 fd = ret;

        // 遍历所有可能的 BPF 命令来查找 pending_bpf_ops
        __u32 cmds[] = {BPF_PROG_LOAD, BPF_MAP_CREATE, BPF_LINK_CREATE};

#pragma unroll
        for (int i = 0; i < sizeof(cmds) / sizeof(cmds[0]); i++)
        {
            __u64 key = make_key(pid, cmds[i]);
            struct bpf_obj_info_t *info = bpf_map_lookup_elem(&pending_bpf_ops, &key);

            if (info)
            {
                // 找到了对应的 pending 操作
                info->fd = fd; // 更新 fd 信息

                // 创建一个新的条目存储到 fd_prog_map 中
                __u64 fd_key = make_key(pid, fd);
                bpf_map_update_elem(&fd_bpf_map, &fd_key, info, BPF_ANY);

                // 从 pending 映射中删除此条目
                bpf_map_delete_elem(&pending_bpf_ops, &key);

                // 根据对象类型记录日志
                if (info->obj_type == OBJ_TYPE_PROG)
                {
                    bpf_printk("BPF Program loaded: PID=%d, FD=%d, proc=%s\n",
                               pid, fd, info->comm);
                }
                else if (info->obj_type == OBJ_TYPE_MAP)
                {
                    bpf_printk("BPF Map created: PID=%d, FD=%d, proc=%s\n",
                               pid, fd, info->comm);
                }
                else if (info->obj_type == OBJ_TYPE_LINK)
                {
                    bpf_printk("BPF Link created: PID=%d, FD=%d, proc=%s\n",
                               pid, fd, info->comm);
                }

                break; // 找到了匹配项，不需要继续循环
            }
        }
    }

    return 0;
}

// 跟踪 close 系统调用以检测对象释放
struct syscalls_enter_close_args
{
    __u16 common_type;
    __u8 common_flags;
    __u8 common_preempt_count;
    __s32 common_pid;

    __s32 __syscall_nr;
    __u64 fd;
};

// SEC("tracepoint/syscalls/sys_enter_close")
// int trace_close(struct syscalls_enter_close_args *ctx)
// {
//     __u32 pid = bpf_get_current_pid_tgid() >> 32;
//     __u32 fd = ctx->fd;

//     bpf_printk("pid: %d, close fd: %d\n", pid, fd);
//     // 查找这个 fd 是否对应 BPF 对象
//     __u64 key = make_key(pid, fd);
//     struct bpf_obj_info_t *info = bpf_map_lookup_elem(&fd_bpf_map, &key);

//     if (info)
//     {
//         // 根据对象类型记录释放事件
//         if (info->obj_type == OBJ_TYPE_PROG)
//         {
//             bpf_printk("BPF Program released: PID=%d, FD=%d, proc=%s\n",
//                        pid, fd, info->comm);
//         }
//         else if (info->obj_type == OBJ_TYPE_MAP)
//         {
//             bpf_printk("BPF Map released: PID=%d, FD=%d, proc=%s\n",
//                        pid, fd, info->comm);
//         }
//         else if (info->obj_type == OBJ_TYPE_LINK)
//         {
//             bpf_printk("BPF Link released: PID=%d, FD=%d, proc=%s\n",
//                        pid, fd, info->comm);
//         }

//         // 从映射中删除这个 fd
//         bpf_map_delete_elem(&fd_bpf_map, &key);
//     }

//     return 0;
// }

SEC("kprobe/filp_close")
int filp_close(struct pt_regs *ctx)
{
    __u32 pid = bpf_get_current_pid_tgid() >> 32;
    char comm[16];
    bpf_get_current_comm(&comm, sizeof(comm));

    // 获取 file 结构体指针（第一个参数）
    struct file *file;
    file = (struct file *)PT_REGS_PARM1(ctx);

    // 读取文件名
    char filename[256] = {0};
    struct dentry *dentry;
    const unsigned char *name;

    // 安全地读取内核内存
    bpf_probe_read_kernel(&dentry, sizeof(dentry), &file->f_path.dentry);
    if (dentry)
    {
        bpf_probe_read_kernel(&name, sizeof(name), &dentry->d_name.name);
        if (name)
        {
            bpf_probe_read_kernel_str(filename, sizeof(filename), name);
            bpf_printk("comm=%-16s pid=%-6d filename=%s\n",
                       comm, pid, filename);
        }
    }

    return 0;
}