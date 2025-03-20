```bash
bpf_prog_load_check_attach args=((enum bpf_prog_type)prog_type=BPF_PROG_TYPE_SOCKET_FILTER, (enum bpf_attach_type)expected_attach_type=BPF_CGROUP_INET_INGRESS, (struct btf *)attach_btf=0x0, (u32)btf_id=0x0/0, (struct bpf_prog *)dst_prog=0x0) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_alloc_no_stats args=((unsigned int)size=0x58/88, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac6742089000 cpu=3 process=(108906:kprobe)

bpf_prog_alloc args=((unsigned int)size=0x58/88, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac6742089000 cpu=3 process=(108906:kprobe)

bpf_prog_free args=((struct bpf_prog *)fp=0xffffac6742089000) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_load args=((union bpf_attr *)attr=0xffffac6743bd7cc8, (bpfptr_t)uattr={{kernel=0xc002799d38,user=0xc002799d38},is_kernel=0x0}, (u32)uattr_size=0x90/144) retval=(int)-9 cpu=3 process=(108906:kprobe)

bpf_prog_free_deferred args=((struct work_struct *)work=0xffff8e98c6c2c410) retval=(void) cpu=3 process=(103401:kworker/3:0)

bpf_prog_load_check_attach args=((enum bpf_prog_type)prog_type=BPF_PROG_TYPE_KPROBE, (enum bpf_attach_type)expected_attach_type=BPF_CGROUP_INET_INGRESS, (struct btf *)attach_btf=0x0, (u32)btf_id=0x0/0, (struct bpf_prog *)dst_prog=0x0) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_alloc_no_stats args=((unsigned int)size=0x78/120, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac67420c1000 cpu=3 process=(108906:kprobe)

bpf_prog_alloc args=((unsigned int)size=0x78/120, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac67420c1000 cpu=3 process=(108906:kprobe)

bpf_prog_calc_tag args=((struct bpf_prog *)fp=0xffffac67420c1000) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_alloc_jited_linfo args=((struct bpf_prog *)prog=0xffffac67420c1000) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_pack_alloc args=((u32)size=0x40/64, (bpf_jit_fill_hole_t)bpf_fill_ill_insns=0xffffffffb66f2820(jit_fill_hole)) retval=(void *)0xffffffffc0068900 cpu=3 process=(108906:kprobe)

bpf_prog_fill_jited_linfo args=((struct bpf_prog *)prog=0xffffac67420c1000, (u32 *)insn_to_jit_off=0xffff8e98488e00a4(0x15/21)) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_jit_attempt_done args=((struct bpf_prog *)prog=0xffffac67420c1000) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_select_runtime args=((struct bpf_prog *)fp=0xffffac67420c1000, (int *)err=0xffffac6743bd7a4c(0)) retval=(struct bpf_prog *)0xffffac67420c1000 cpu=3 process=(108906:kprobe)

bpf_prog_load args=((union bpf_attr *)attr=0xffffac6743bd7b88, (bpfptr_t)uattr={{kernel=0xc0027999e8,user=0xc0027999e8},is_kernel=0x0}, (u32)uattr_size=0x90/144) retval=(int)8 cpu=3 process=(108906:kprobe)

bpf_prog_put_deferred args=((struct work_struct *)work=0xffff8e98c6c2c410) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_release args=((struct inode *)inode=0xffff8e984046e300, (struct file *)filp=0xffff8e984d33b100) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_load_check_attach args=((enum bpf_prog_type)prog_type=BPF_PROG_TYPE_KPROBE, (enum bpf_attach_type)expected_attach_type=BPF_CGROUP_INET_INGRESS, (struct btf *)attach_btf=0x0, (u32)btf_id=0x0/0, (struct bpf_prog *)dst_prog=0x0) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_alloc_no_stats args=((unsigned int)size=0xf8/248, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac67420d9000 cpu=3 process=(108906:kprobe)

bpf_prog_alloc args=((unsigned int)size=0xf8/248, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac67420d9000 cpu=3 process=(108906:kprobe)

bpf_prog_calc_tag args=((struct bpf_prog *)fp=0xffffac67420d9000) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_realloc args=((struct bpf_prog *)fp_old=0xffffac67420d9000, (unsigned int)size=0x128/296, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac67420d9000 cpu=3 process=(108906:kprobe)

bpf_prog_alloc_jited_linfo args=((struct bpf_prog *)prog=0xffffac67420d9000) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_pack_alloc args=((u32)size=0xc0/192, (bpf_jit_fill_hole_t)bpf_fill_ill_insns=0xffffffffb66f2820(jit_fill_hole)) retval=(void *)0xffffffffc0084700 cpu=3 process=(108906:kprobe)

bpf_prog_fill_jited_linfo args=((struct bpf_prog *)prog=0xffffac67420d9000, (u32 *)insn_to_jit_off=0xffff8e9842e58304(0x15/21)) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_jit_attempt_done args=((struct bpf_prog *)prog=0xffffac67420d9000) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_select_runtime args=((struct bpf_prog *)fp=0xffffac67420d9000, (int *)err=0xffffac6743bd7c2c(0)) retval=(struct bpf_prog *)0xffffac67420d9000 cpu=3 process=(108906:kprobe)

bpf_prog_load args=((union bpf_attr *)attr=0xffffac6743bd7d68, (bpfptr_t)uattr={{kernel=0xc00279aba8,user=0xc00279aba8},is_kernel=0x0}, (u32)uattr_size=0x90/144) retval=(int)8 cpu=3 process=(108906:kprobe)

bpf_prog_load_check_attach args=((enum bpf_prog_type)prog_type=BPF_PROG_TYPE_KPROBE, (enum bpf_attach_type)expected_attach_type=BPF_CGROUP_INET_INGRESS, (struct btf *)attach_btf=0x0, (u32)btf_id=0x0/0, (struct bpf_prog *)dst_prog=0x0) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_alloc_no_stats args=((unsigned int)size=0x58/88, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac67420f9000 cpu=3 process=(108906:kprobe)

bpf_prog_alloc args=((unsigned int)size=0x58/88, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac67420f9000 cpu=3 process=(108906:kprobe)

bpf_prog_calc_tag args=((struct bpf_prog *)fp=0xffffac67420f9000) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_alloc_jited_linfo args=((struct bpf_prog *)prog=0xffffac67420f9000) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_pack_alloc args=((u32)size=0x40/64, (bpf_jit_fill_hole_t)bpf_fill_ill_insns=0xffffffffb66f2820(jit_fill_hole)) retval=(void *)0xffffffffc0087280 cpu=3 process=(108906:kprobe)

bpf_prog_fill_jited_linfo args=((struct bpf_prog *)prog=0xffffac67420f9000, (u32 *)insn_to_jit_off=0xffff8e98415ae634(0xd/13)) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_jit_attempt_done args=((struct bpf_prog *)prog=0xffffac67420f9000) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_select_runtime args=((struct bpf_prog *)fp=0xffffac67420f9000, (int *)err=0xffffac6743bd7a34(0)) retval=(struct bpf_prog *)0xffffac67420f9000 cpu=3 process=(108906:kprobe)

bpf_prog_load args=((union bpf_attr *)attr=0xffffac6743bd7b70, (bpfptr_t)uattr={{kernel=0xc00279a4c0,user=0xc00279a4c0},is_kernel=0x0}, (u32)uattr_size=0x90/144) retval=(int)9 cpu=3 process=(108906:kprobe)

bpf_prog_attach_check_attach_type args=((const struct bpf_prog *)prog=0xffffac67420f9000, (enum bpf_attach_type)attach_type=BPF_PERF_EVENT) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_put_deferred args=((struct work_struct *)work=0xffff8e98c58ec410) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_release args=((struct inode *)inode=0xffff8e984046e300, (struct file *)filp=0xffff8e984d33b500) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_attach_check_attach_type args=((const struct bpf_prog *)prog=0xffffac67420d9000, (enum bpf_attach_type)attach_type=BPF_PERF_EVENT) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_array_copy args=((struct bpf_prog_array *)old_array=0x0, (struct bpf_prog *)exclude_prog=0x0, (struct bpf_prog *)include_prog=0xffffac67420d9000, (u64)bpf_cookie=0x0/0, (struct bpf_prog_array **)new_array=0xffffac6743bd7a10) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_array_free_sleepable args=((struct bpf_prog_array *)progs=0x0) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_get_target_btf args=((const struct bpf_prog *)prog=0xffffac67420d9000) retval=(struct btf *)0x0 cpu=3 process=(108906:kprobe)

bpf_prog_get_stats args=((const struct bpf_prog *)prog=0xffffac67420d9000, (struct bpf_prog_kstats *)stats=0xffffac6743bd7910) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_get_info_by_fd args=((struct file *)file=0xffff8e984d33b100, (struct bpf_prog *)prog=0xffffac67420d9000, (const union bpf_attr *)attr=0xffffac6743bd7b58, (union bpf_attr *)uattr=0xc00279aee8) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_get_target_btf args=((const struct bpf_prog *)prog=0xffffac67420d9000) retval=(struct btf *)0x0 cpu=3 process=(108906:kprobe)

bpf_prog_get_stats args=((const struct bpf_prog *)prog=0xffffac67420d9000, (struct bpf_prog_kstats *)stats=0xffffac6743bd7918) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_get_info_by_fd args=((struct file *)file=0xffff8e984d33b100, (struct bpf_prog *)prog=0xffffac67420d9000, (const union bpf_attr *)attr=0xffffac6743bd7b60, (union bpf_attr *)uattr=0xc00279aee8) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_free args=((struct bpf_prog *)fp=0xffffac67420c1000) retval=(void) cpu=3 process=(0:swapper/3)

bpf_prog_pack_free args=((void *)ptr=0xffffffffc0068900, (u32)size=0x40/64) retval=(void) cpu=3 process=(103401:kworker/3:0)

bpf_prog_free_deferred args=((struct work_struct *)work=0xffff8e98c6c2c410) retval=(void) cpu=3 process=(103401:kworker/3:0)

bpf_prog_free args=((struct bpf_prog *)fp=0xffffac67420f9000) retval=(void) cpu=3 process=(0:swapper/3)

bpf_prog_pack_free args=((void *)ptr=0xffffffffc0087280, (u32)size=0x40/64) retval=(void) cpu=3 process=(103401:kworker/3:0)

bpf_prog_free_deferred args=((struct work_struct *)work=0xffff8e98c58ec410) retval=(void) cpu=3 process=(103401:kworker/3:0)

bpf_prog_array_copy args=((struct bpf_prog_array *)old_array=0xffff8e9841583360, (struct bpf_prog *)exclude_prog=0xffffac67420d9000, (struct bpf_prog *)include_prog=0x0, (u64)bpf_cookie=0x0/0, (struct bpf_prog_array **)new_array=0xffffac6743bd79d0) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_array_free_sleepable args=((struct bpf_prog_array *)progs=0xffff8e9841583360) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_put args=((struct bpf_prog *)prog=0xffffac67420d9000) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_put_deferred args=((struct work_struct *)work=0xffff8e98c6c2b410) retval=(void) cpu=3 process=(108906:kprobe)

bpf_prog_release args=((struct inode *)inode=0xffff8e984046e300, (struct file *)filp=0xffff8e984d33b100) retval=(int)0 cpu=3 process=(108906:kprobe)

bpf_prog_free args=((struct bpf_prog *)fp=0xffffac67420d9000) retval=(void) cpu=3 process=(0:swapper/3)

bpf_prog_pack_free args=((void *)ptr=0xffffffffc0084700, (u32)size=0xc0/192) retval=(void) cpu=3 process=(103401:kworker/3:0)

bpf_prog_free_deferred args=((struct work_struct *)work=0xffff8e98c6c2b410) retval=(void) cpu=3 process=(103401:kworker/3:0)
```