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

#define OK 0
#define WRONG_ANSWER 1
#define TLE 2
#define MLE 3
#define nullptr NULL
#define VMRSS_LINE 17
#define VMSIZE_LINE 13
#define PROCESS_ITEM 14
 

int max3(int a, int b) {
    return a > b ? a : b;
}

//获取进程占用内存
unsigned int get_proc_mem3(unsigned int p){
	char file[64] = { 0 };       //文件名
    FILE *fd;                  //定义文件指针fd
    char line_buff[256] = { 0 }; //读取行的缓冲区
    sprintf(file, "/proc/%d/status", p);
    fd = fopen(file, "r"); //以R读的方式打开文件再赋给指针fd
    //获取vmrss:实际物理内存占用
    char name[32]; //存放项目名称
    int vmrss;     //存放内存
    //读取VmRSS这一行的数据
    for (int i = 0; i < VMRSS_LINE - 1; i++){
        char *ret = fgets(line_buff, sizeof(line_buff), fd);
    }
    char *ret1 = fgets(line_buff, sizeof(line_buff), fd);
    sscanf(line_buff, "%s %d", name, &vmrss);
    //fprintf(stderr, "====%s：%d====\n", name, vmrss);
    fclose(fd); //关闭文件fd
    return vmrss;
}

void setup_sandbox2() {
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

    if (seccomp_load(ctx) < 0) {
        perror("seccomp_load failed");
        exit(EXIT_FAILURE);
    }
}

char pipe_u2j[50];
char pipe_j2u[50];

void run_judge3(int offset) {
    // 打开管道文件。注意打开顺序，否则会造成死锁！
    int in = open(pipe_u2j, O_RDONLY);
    int out = open(pipe_j2u, O_WRONLY);

    // 重定向标准输入输出
    dup2(in, 0);
    dup2(out, 1);
    close(in);
    close(out);

    char judge_path[50];
    sprintf(judge_path, "/app/InteractiveTest/judge%d", offset);
    // 执行裁判程序
    execl(judge_path, judge_path, NULL);
}

void run_user3() {
    // 打开管道文件。注意打开顺序，否则会造成死锁！
    int out = open(pipe_u2j, O_WRONLY);
    int in = open(pipe_j2u, O_RDONLY);

    // 重定向标准输入输出
    dup2(in, 0);
    dup2(out, 1);
    close(in);
    close(out);

    char user_path[50] = "/app/InteractiveTest/user";
    // 执行用户程序
    setup_sandbox2();
    execl(user_path, user_path, NULL);
}


int InteractiveTest(int time, int memory, int offset, int *res, int *costTime, int *costMemory) {
    // 创建管道文件
    sprintf(pipe_j2u, "/app/InteractiveTest/j2u%d.fifo", offset);
    sprintf(pipe_u2j, "/app/InteractiveTest/u2j%d.fifo", offset);
    mkfifo(pipe_j2u, 0644);
    mkfifo(pipe_u2j, 0644);

    pid_t pid_j, pid_u;

    // 创建裁判进程
    pid_j = fork();
    if (pid_j < 0) {
        printf("0 0\n");
        return 1;
    } else if (pid_j == 0) {
        run_judge3(offset);
        return 0;
    }

    // 创建用户进程
    pid_u = fork();
    if (pid_u < 0) {
        printf("0 0\n");
        return 1;
    } else if (pid_u == 0) {
        run_user3();
        return 0;
    }

    unsigned maxmemory = 0, tottime = 0;
    
    // 等待进程运行结束，并判定结果
    int stat_j, stat_u;
    while(1){
        if(waitpid(pid_j, &stat_j, WNOHANG) && waitpid(pid_u, &stat_u, WNOHANG)) {
            break;
        }
        maxmemory = max3(maxmemory, get_proc_mem3(pid_u));
        if(tottime > time) {
            kill(pid_j, SIGKILL);
            kill(pid_u, SIGKILL);
            waitpid(pid_j, &stat_j, 0);
            waitpid(pid_u, &stat_u, 0);
            *res = TLE;
            *costTime = tottime;
            *costMemory = maxmemory;
            printf("%d %d\nTLE\n", tottime, maxmemory);
            return TLE;
        }
        if(maxmemory > memory) {
            kill(pid_j, SIGKILL);
            kill(pid_u, SIGKILL);
            waitpid(pid_j, &stat_j, 0);
            waitpid(pid_u, &stat_u, 0);
            *res = MLE;
            *costTime = tottime;
            *costMemory = maxmemory;
            printf("%d %d\nMLE\n", tottime, maxmemory);
            return MLE;
        }
        tottime += 4;
        usleep(40000);
    }
    usleep(20000);
    waitpid(pid_j, &stat_j, 0);
    waitpid(pid_u, &stat_u, 0);
    printf("%d %d\n", tottime, maxmemory);
    *costTime = tottime;
    *costMemory = maxmemory;
    if (WIFEXITED(stat_u) || (WIFSIGNALED(stat_u) && WTERMSIG(stat_u) == SIGPIPE)) {
        // 用户程序正常退出，或由于 SIGPIPE 退出，需要裁判程序判定
        if (WIFEXITED(stat_j)) {
            // 裁判程序正常退出
            switch (WEXITSTATUS(stat_j)) {
            case OK:
                printf("Accepted\n");
                *res = OK;
                break;
            case WRONG_ANSWER:
                printf("Wrong answer\n");
                *res = WRONG_ANSWER;
                break;
            default:
                printf("Invalid judge exit code\n");
                *res = WRONG_ANSWER;
                break;
            }
        } else if (WIFSIGNALED(stat_j) && WTERMSIG(stat_j) == SIGPIPE) {
            // 裁判程序由于 SIGPIPE 退出
            printf("Wrong answer\n");
            *res = WRONG_ANSWER;
        } else {
            // 裁判程序异常退出
            printf("Judge exit abnormally\n");
            *res = WRONG_ANSWER;
        }
    } else {
        // 用户程序运行时错误
        printf("Runtime error\n");
        *res = WRONG_ANSWER;
    }
    char rm_pipe_u2j[50];
    char rm_pipe_j2u[50];
    sprintf(rm_pipe_j2u, "rm /app/InteractiveTest/j2u%d.fifo", offset);
    sprintf(rm_pipe_u2j, "rm /app/InteractiveTest/u2j%d.fifo", offset);
    system(rm_pipe_j2u);
    system(rm_pipe_u2j);
    return *res;
}