```bash
bpf_prog_load_check_attach args=((enum bpf_prog_type)prog_type=BPF_PROG_TYPE_SOCKET_FILTER, (enum bpf_attach_type)expected_attach_type=BPF_CGROUP_INET_INGRESS, (struct btf *)attach_btf=0x0, (u32)btf_id=0x0/0, (struct bpf_prog *)dst_prog=0x0) retval=(int)0 cpu=0 process=(93479:fentry)

bpf_prog_alloc_no_stats args=((unsigned int)size=0x58/88, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac6744405000 cpu=0 process=(93479:fentry)

bpf_prog_alloc args=((unsigned int)size=0x58/88, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac6744405000 cpu=0 process=(93479:fentry)

bpf_prog_kallsyms_del_all args=((struct bpf_prog *)fp=0xffffac6744405000) retval=(void) cpu=0 process=(93479:fentry)

bpf_prog_free args=((struct bpf_prog *)fp=0xffffac6744405000) retval=(void) cpu=0 process=(93479:fentry)

bpf_prog_load args=((union bpf_attr *)attr=0xffffac674ab2fcd0, (bpfptr_t)uattr={{kernel=0xc000c71d38,user=0xc000c71d38},is_kernel=0x0}, (u32)uattr_size=0x90/144) retval=(int)-9 cpu=0 process=(93479:fentry)

bpf_prog_free_deferred args=((struct work_struct *)work=0xffff8e9935871c10) retval=(void) cpu=0 process=(92158:kworker/0:2)

bpf_prog_load_check_attach args=((enum bpf_prog_type)prog_type=BPF_PROG_TYPE_KPROBE, (enum bpf_attach_type)expected_attach_type=BPF_CGROUP_INET_INGRESS, (struct btf *)attach_btf=0x0, (u32)btf_id=0x0/0, (struct bpf_prog *)dst_prog=0x0) retval=(int)0 cpu=0 process=(93479:fentry)

bpf_prog_alloc_no_stats args=((unsigned int)size=0x78/120, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac674789d000 cpu=0 process=(93479:fentry)

bpf_prog_alloc args=((unsigned int)size=0x78/120, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac674789d000 cpu=0 process=(93479:fentry)

bpf_prog_calc_tag args=((struct bpf_prog *)fp=0xffffac674789d000) retval=(int)0 cpu=0 process=(93479:fentry)

bpf_prog_alloc_jited_linfo args=((struct bpf_prog *)prog=0xffffac674789d000) retval=(int)0 cpu=0 process=(93479:fentry)

bpf_prog_pack_alloc args=((u32)size=0x40/64, (bpf_jit_fill_hole_t)bpf_fill_ill_insns=0xffffffffb66f2820(jit_fill_hole)) retval=(void *)0xffffffffc0068900 cpu=0 process=(93479:fentry)

bpf_prog_fill_jited_linfo args=((struct bpf_prog *)prog=0xffffac674789d000, (u32 *)insn_to_jit_off=0xffff8e98d4253e24(0x15/21)) retval=(void) cpu=0 process=(93479:fentry)

bpf_prog_jit_attempt_done args=((struct bpf_prog *)prog=0xffffac674789d000) retval=(void) cpu=0 process=(93479:fentry)

bpf_prog_select_runtime args=((struct bpf_prog *)fp=0xffffac674789d000, (int *)err=0xffffac674ab2fcac(0)) retval=(struct bpf_prog *)0xffffac674789d000 cpu=0 process=(93479:fentry)

bpf_prog_kallsyms_add args=((struct bpf_prog *)fp=0xffffac674789d000) retval=(void) cpu=0 process=(93479:fentry)

bpf_prog_load args=((union bpf_attr *)attr=0xffffac674ab2fde8, (bpfptr_t)uattr={{kernel=0xc000c719e8,user=0xc000c719e8},is_kernel=0x0}, (u32)uattr_size=0x90/144) retval=(int)9 cpu=0 process=(93479:fentry)

bpf_prog_kallsyms_del_all args=((struct bpf_prog *)fp=0xffffac674789d000) retval=(void) cpu=0 process=(93479:fentry)

bpf_prog_put_deferred args=((struct work_struct *)work=0xffff8e9935871c10) retval=(void) cpu=0 process=(93479:fentry)

bpf_prog_release args=((struct inode *)inode=0xffff8e984046e300, (struct file *)filp=0xffff8e98c6dc6600) retval=(int)0 cpu=0 process=(93479:fentry)

bpf_prog_free args=((struct bpf_prog *)fp=0xffffac674789d000) retval=(void) cpu=0 process=(93479:fentry)

bpf_prog_load_check_attach args=((enum bpf_prog_type)prog_type=BPF_PROG_TYPE_TRACING, (enum bpf_attach_type)expected_attach_type=BPF_TRACE_FENTRY, (struct btf *)attach_btf=0xffff8e984104af00, (u32)btf_id=0x1460b/83467, (struct bpf_prog *)dst_prog=0x0) retval=(int)0 cpu=0 process=(93479:fentry)

bpf_prog_alloc_no_stats args=((unsigned int)size=0x120/288, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac67478b5000 cpu=0 process=(93479:fentry)

bpf_prog_alloc args=((unsigned int)size=0x120/288, (gfp_t)gfp_extra_flags=0x100cc0/1051840) retval=(struct bpf_prog *)0xffffac67478b5000 cpu=0 process=(93479:fentry)

bpf_prog_calc_tag args=((struct bpf_prog *)fp=0xffffac67478b5000) retval=(int)0 cpu=0 process=(93479:fentry)

bpf_prog_alloc_jited_linfo args=((struct bpf_prog *)prog=0xffffac67478b5000) retval=(int)0 cpu=0 process=(93479:fentry)

bpf_prog_pack_alloc args=((u32)size=0x1c0/448, (bpf_jit_fill_hole_t)bpf_fill_ill_insns=0xffffffffb66f2820(jit_fill_hole)) retval=(void *)0xffffffffc016bf80 cpu=0 process=(93479:fentry)

bpf_prog_fill_jited_linfo args=((struct bpf_prog *)prog=0xffffac67478b5000, (u32 *)insn_to_jit_off=0xffff8e97477eb284(0x12/18)) retval=(void) cpu=0 process=(93479:fentry)

bpf_prog_jit_attempt_done args=((struct bpf_prog *)prog=0xffffac67478b5000) retval=(void) cpu=0 process=(93479:fentry)

bpf_prog_select_runtime args=((struct bpf_prog *)fp=0xffffac67478b5000, (int *)err=0xffffac674ab2fc2c(0)) retval=(struct bpf_prog *)0xffffac67478b5000 cpu=0 process=(93479:fentry)

bpf_prog_kallsyms_add args=((struct bpf_prog *)fp=0xffffac67478b5000) retval=(void) cpu=0 process=(93479:fentry)

bpf_prog_load args=((union bpf_attr *)attr=0xffffac674ab2fd68, (bpfptr_t)uattr={{kernel=0xc000c72ba8,user=0xc000c72ba8},is_kernel=0x0}, (u32)uattr_size=0x90/144) retval=(int)9 cpu=0 process=(93479:fentry)

bpf_prog_pack_alloc args=((u32)size=0xa4/164, (bpf_jit_fill_hole_t)bpf_fill_ill_insns=0xffffffffb66f2820(jit_fill_hole)) retval=(void *)0xffffffffc0087280 cpu=0 process=(93479:fentry)

bpf_prog_pack_free args=((void *)ptr=0xffffffffc0068900, (u32)size=0x40/64) retval=(void) cpu=0 process=(92158:kworker/0:2)

bpf_prog_free_deferred args=((struct work_struct *)work=0xffff8e9935871c10) retval=(void) cpu=0 process=(92158:kworker/0:2)

bpf_prog_get_target_btf args=((const struct bpf_prog *)prog=0xffffac67478b5000) retval=(struct btf *)0xffff8e984104af00 cpu=1 process=(93479:fentry)

bpf_prog_get_stats args=((const struct bpf_prog *)prog=0xffffac67478b5000, (struct bpf_prog_kstats *)stats=0xffffac674ab2faf0) retval=(void) cpu=1 process=(93479:fentry)

bpf_prog_get_info_by_fd args=((struct file *)file=0xffff8e98c6dc6600, (struct bpf_prog *)prog=0xffffac67478b5000, (const union bpf_attr *)attr=0xffffac674ab2fd38, (union bpf_attr *)uattr=0xc000c72ee8) retval=(int)0 cpu=1 process=(93479:fentry)

bpf_prog_get_target_btf args=((const struct bpf_prog *)prog=0xffffac67478b5000) retval=(struct btf *)0xffff8e984104af00 cpu=1 process=(93479:fentry)

bpf_prog_get_stats args=((const struct bpf_prog *)prog=0xffffac67478b5000, (struct bpf_prog_kstats *)stats=0xffffac674ab2f7f0) retval=(void) cpu=1 process=(93479:fentry)

bpf_prog_get_info_by_fd args=((struct file *)file=0xffff8e98c6dc6600, (struct bpf_prog *)prog=0xffffac67478b5000, (const union bpf_attr *)attr=0xffffac674ab2fa38, (union bpf_attr *)uattr=0xc000c72ee8) retval=(int)0 cpu=1 process=(93479:fentry)

bpf_prog_kallsyms_del_all args=((struct bpf_prog *)fp=0xffffac67478b5000) retval=(void) cpu=3 process=(93479:fentry)

bpf_prog_put_deferred args=((struct work_struct *)work=0xffff8e9935875c10) retval=(void) cpu=3 process=(93479:fentry)

bpf_prog_release args=((struct inode *)inode=0xffff8e984046e300, (struct file *)filp=0xffff8e98c6dc6600) retval=(int)0 cpu=3 process=(93479:fentry)

bpf_prog_free args=((struct bpf_prog *)fp=0xffffac67478b5000) retval=(void) cpu=3 process=(750:containerd)

bpf_prog_pack_free args=((void *)ptr=0xffffffffc016bf80, (u32)size=0x1c0/448) retval=(void) cpu=3 process=(88444:kworker/3:2)

bpf_prog_free_deferred args=((struct work_struct *)work=0xffff8e9935875c10) retval=(void) cpu=3 process=(88444:kworker/3:2)

bpf_prog_pack_free args=((void *)ptr=0xffffffffc0087280, (u32)size=0xa4/164) retval=(void) cpu=0 process=(92158:kworker/0:2)
```