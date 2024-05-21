#include <stdio.h>
#include <unistd.h>
#include <sys/prctl.h>
#include <sys/syscall.h>
#include <seccomp.h>
#include <stdlib.h>

// void setup_sandbox() {
//     scmp_filter_ctx ctx = seccomp_init(SCMP_ACT_ALLOW); // 默认允许所有系统调用

//     // 禁止一些基本的系统调用
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(fork), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(vfork), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(clone), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(setuid), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(setgid), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(chroot), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(pivot_root), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(mount), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(umount), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(open), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(truncate), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(ftruncate), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(unlink), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(rmdir), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(rename), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(mkdir), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(socket), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(bind), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(listen), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(accept), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(sendto), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(sendmsg), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(recvfrom), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(recvmsg), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(dup), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(dup2), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(fcntl), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(sigaction), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(sigprocmask), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(sigreturn), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(times), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(getrusage), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(sysinfo), 0);
//     seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(nanosleep), 0);

//     if (seccomp_load(ctx) < 0) {
//         perror("seccomp_load failed");
//         exit(EXIT_FAILURE);
//     }
    
// }

int main() {
    // 配置seccomp
    // setup_sandbox();
    printf("%d\n", sizeof(long));
    // 尝试调用sleep，将被禁用
    usleep(5000000);
    return 0;
}