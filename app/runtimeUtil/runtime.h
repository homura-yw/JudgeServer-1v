#include <unistd.h>
#include <signal.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <seccomp.h>
#include <assert.h>

#define OK                  0
#define WRONG_ANSWER        1
#define TLE                 2
#define MLE                 3
#define CE                  4
#define RE                  5
#define CPU_START_POS       14   
 

int max(int a, int b) {
    return a > b ? a : b;
}

//获取进程占用内存
unsigned int get_proc_mem(unsigned int pid){
	FILE *fd;
    int   vmrss;
    char *valid, sbuff[32], tbuff[1024];
    
    sprintf(tbuff, "/proc/%d/status", pid);                                    
    /* 在proc目录下查找进程对应文件 */
    
    fd = fopen(tbuff, "r");
    if(fd == NULL) {
        return -1;
    }
    
    while (1) {                                                                
        /* 对文件内容进行逐行搜索 */
        assert(fgets(tbuff, sizeof(tbuff), fd) != NULL);                       
        /* 文件读取出错 */
        valid = strstr(tbuff, "VmPeak");                                        
        /* 在该行内容中搜索关键词 */
        if (valid != NULL) {                                                   
            /* 结果非空则表示搜索成功 */
            break;
        }
    }
    
    sscanf(tbuff, "%s %d", sbuff, &vmrss);
    /* 将该行内容拆成符合需要的格式 */
    fclose(fd);
    
    return vmrss;
}

char *get_items_by_pos(char *buff, unsigned int numb) {
    char *crpos;
    int   i, ttlen, count;
    
    crpos = buff;
    ttlen = strlen(buff);
    count = 0;
    
    for (i = 0; i < ttlen; i++) {
        if (' ' == *crpos) {                                                   
            /* 以空格为标记符进行识别 */
            count++;
            if (count == (numb - 1)) {                                         
                /* 全部个数都找完了 */
                crpos++;
                break;
            }
        }
        crpos++;
    }
    
    return crpos;
}

long get_pro_cpu_time(unsigned int pid) {
    FILE   *fd;
    char   *vpos, buff[1024];
    long    utime, stime, cutime, cstime;
    
    sprintf(buff, "/proc/%d/stat", pid);                                       
    /* 读取进程的状态文件 */
    
    fd = fopen(buff, "r");
    if(fd == NULL) {
        return -1;
    }
    if(fgets(buff, sizeof(buff), fd) == NULL) {
        return -1;
    }                             
    /* 读取文件内容到缓冲区 */
    
    vpos = get_items_by_pos(buff, CPU_START_POS);                              
    /* 读取指定的条目内容 */
    sscanf(vpos, "%ld %ld %ld %ld", &utime, &stime, &cutime, &cstime);         
    /* 将条目内容拆分成实际的数据 */
    
    fclose(fd);
    
    return (utime + stime);
}

void setup_sandbox() {
    scmp_filter_ctx ctx = seccomp_init(SCMP_ACT_ALLOW); // 默认允许所有系统调用

    // 禁止一些基本的系统调用
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(fork), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(vfork), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(clone), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(setuid), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(setgid), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(chroot), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(pivot_root), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(mount), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(umount), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(open), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(truncate), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(ftruncate), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(unlink), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(rmdir), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(rename), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(mkdir), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(socket), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(bind), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(listen), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(accept), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(sendto), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(sendmsg), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(recvfrom), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(recvmsg), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(dup), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(dup2), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(fcntl), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(sigaction), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(sigprocmask), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(sigreturn), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(times), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(getrusage), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(sysinfo), 0);
    seccomp_rule_add(ctx, SCMP_ACT_KILL, SCMP_SYS(nanosleep), 0);

    if (seccomp_load(ctx) < 0) {
        perror("seccomp_load failed");
        exit(EXIT_FAILURE);
    }
}

void load(int *x, int y) {
    if(y >= 0) *x = y;
}