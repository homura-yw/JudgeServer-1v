#include <unistd.h>
#include <signal.h>
#include <fcntl.h>
#include <sys/time.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define OK 0
#define WRONG_ANSWER 1
#define TLE 2
#define MLE 3
#define nullptr NULL
#define VMRSS_LINE 17
#define VMSIZE_LINE 13
#define PROCESS_ITEM 14
#define nullptr NULL

int max1(int a, int b) {
    return a > b ? a : b;
}

//获取进程占用内存
unsigned int get_proc_mem1(unsigned int p){
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

void runUser1(int offset) {
    char input_path[50];
    char output_path[50];
    char user_code[50] = "/app/NormalTest/user";
    sprintf(input_path, "/app/NormalTest/input%d", offset);
    sprintf(output_path, "/app/NormalTest/output%d", offset);
    freopen(input_path, "r", stdin);
    freopen(output_path, "w", stdout);
    execl(user_code, user_code, NULL);
}

void runJudge1(int offset) {
    char judge_code[50];
    sprintf(judge_code, "/app/NormalTest/judge%d", offset);
    execl(judge_code, judge_code, NULL);
}
int NormalTest(int time, int memory, int offset, int *res, int *costTime, int *costMemory){
    int pidu = fork();
    if(pidu == 0) {
        runUser1(offset);
        return 0;
    }
	int tottime = 0, maxMemory = 0, statu, statj;
	while(1) {
		if(waitpid(pidu, &statu, WNOHANG)) {
			break;
		}
		int nowmemory = get_proc_mem1(pidu);
		maxMemory = max1(maxMemory, nowmemory);
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
        runJudge1(offset);
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
            printf("sigpipe\n");
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