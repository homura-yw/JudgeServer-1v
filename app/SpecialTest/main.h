#include <unistd.h>
#include <signal.h>
#include <fcntl.h>
#include <sys/time.h>
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
#define nullptr NULL

int max2(int a, int b) {
    return a > b ? a : b;
}

//获取进程占用内存
unsigned int get_proc_mem2(unsigned int p){
	char file[64] = {0};       //文件名
    FILE *fd;                  //定义文件指针fd
    char line_buff[256] = {0}; //读取行的缓冲区
    sprintf(file, "/proc/%d/status", p);
    fd = fopen(file, "r"); //以R读的方式打开文件再赋给指针fd
    //获取vmrss:实际物理内存占用
    int i;
    char name[32]; //存放项目名称
    int vmrss;     //存放内存
    //读取VmRSS这一行的数据
    for (i = 0; i < VMRSS_LINE - 1; i++)
    {
        char *ret = fgets(line_buff, sizeof(line_buff), fd);
    }
    char *ret1 = fgets(line_buff, sizeof(line_buff), fd);
    sscanf(line_buff, "%s %d", name, &vmrss);
    //fprintf(stderr, "====%s：%d====\n", name, vmrss);
    fclose(fd); //关闭文件fd
    return vmrss;
}

void setup_sandbox3() {
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

void runUser2(int offset) {
    char input_path[50];
    char output_path[50];
    char user_code[50] = "/app/SpecialTest/user";
    sprintf(input_path, "/app/SpecialTest/input%d", offset);
    sprintf(output_path, "/app/SpecialTest/output%d", offset);
    freopen(input_path, "r", stdin);
    freopen(output_path, "w", stdout);
    setup_sandbox3();
    execl(user_code, user_code, NULL);
}

void runJudge2(int offset) {
    char judge_code[50];
    sprintf(judge_code, "/app/SpecialTest/judge%d", offset);
    execl(judge_code, judge_code, NULL);
}
int SpecialTest(int time, int memory, int offset, int *res, int *costTime, int *costMemory){
    int pidu = fork();
    if(pidu == 0) {
        runUser2(offset);
        return 0;
    }
	int tottime = 0, maxMemory = 0, statu, statj;
	while(1) {
		if(waitpid(pidu, &statu, WNOHANG)) {
			break;
		}
		int nowmemory = get_proc_mem2(pidu);
		maxMemory = max2(maxMemory, nowmemory);
		if(tottime > time) {
			kill(pidu, SIGKILL);
			waitpid(pidu, &statu, 0);
			printf("%d %d\nTLE\n", tottime, maxMemory);
            *res = TLE;
            *costTime = tottime;
            *costMemory = maxMemory;
			return TLE;
		}
		if(nowmemory > memory) {
			kill(pidu, SIGKILL);
			waitpid(pidu, &statu, 0);
			printf("%d %d\nMLE\n", tottime, maxMemory);
            *res = MLE;
            *costTime = tottime;
            *costMemory = maxMemory;
			return MLE;
		}
		usleep(40000);
		tottime += 4;
	}
    int pidj = fork();
    if(pidj == 0) {
        runJudge2(offset);
        return 0;
    }
    printf("%d %d\n", tottime, maxMemory);
    *costTime = tottime;
    *costMemory = maxMemory;
    waitpid(pidj, &statj, 0);
	if (WIFEXITED(statu) || (WIFSIGNALED(statu) && WTERMSIG(statu) == SIGPIPE)) {
        // 用户程序正常退出，或由于 SIGPIPE 退出，需要裁判程序判定
        if (WIFEXITED(statj)) {
            // 裁判程序正常退出
            switch (WEXITSTATUS(statj)) {
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
        } else if (WIFSIGNALED(statj) && WTERMSIG(statj) == SIGPIPE) {
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
    return *res;

}