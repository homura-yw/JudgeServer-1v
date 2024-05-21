#include "../runtimeUtil/runtime.h"

char pipe_u2j[50];
char pipe_j2u[50];

char rm_pipe_u2j[50];
char rm_pipe_j2u[50];

void run_judge(int offset) {
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

void run_user() {
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
    setup_sandbox();
    execl(user_path, user_path, NULL);
}

int main(int args, char **argc) {
    // 创建管道文件
    int offset = atoi(argc[1]), time = atoi(argc[2]), memory = atoi(argc[3]);
    int res = 0, costTime = 0, costMemory = 0, tottime = 0;
    sprintf(pipe_j2u, "/app/InteractiveTest/j2u%d.fifo", offset);
    sprintf(pipe_u2j, "/app/InteractiveTest/u2j%d.fifo", offset);
    mkfifo(pipe_j2u, 0644);
    mkfifo(pipe_u2j, 0644);

    pid_t pid_j, pid_u;

    // 创建裁判进程
    pid_j = fork();
    if (pid_j < 0) {
        printf("-1 0 0\n");
        return 1;
    } else if (pid_j == 0) {
        run_judge(offset);
        return 0;
    }

    // 创建用户进程
    pid_u = fork();
    if (pid_u < 0) {
        printf("-1 0 0\n");
        return 1;
    } else if (pid_u == 0) {
        run_user();
        return 0;
    }
    
    // 等待进程运行结束，并判定结果
    int stat_j, stat_u;
    while(1){
        if(waitpid(pid_j, &stat_j, WNOHANG) && waitpid(pid_u, &stat_u, WNOHANG)) {
            break;
        }
        load(&costMemory, get_proc_mem(pid_u));
        load(&costTime, get_pro_cpu_time(pid_u));
        if(tottime > time + 100 || costTime > time) {

            kill(pid_j, SIGKILL);
            kill(pid_u, SIGKILL);
            waitpid(pid_j, &stat_j, 0);
            waitpid(pid_u, &stat_u, 0);

            res = TLE;
            costTime = time;

            printf("%d %d %d\n", res, costTime, costMemory);
            return res;

        }
        if(costMemory > memory) {

            kill(pid_j, SIGKILL);
            kill(pid_u, SIGKILL);
            waitpid(pid_j, &stat_j, 0);
            waitpid(pid_u, &stat_u, 0);

            res = MLE;

            printf("%d %d %d\n", res, costTime, costMemory);
            return res;

        }
        tottime += 4;
        usleep(40000);
    }

    usleep(20000);
    
    waitpid(pid_j, &stat_j, 0);
    waitpid(pid_u, &stat_u, 0);

    if(costTime > time) {
        printf("%d %d %d\n", TLE, costTime, costMemory);
        return TLE;
    }

    if(costMemory > memory) {
        printf("%d %d %d\n", MLE, costTime, costMemory);
        return MLE;
    }

    if (WIFEXITED(stat_u) || (WIFSIGNALED(stat_u) && WTERMSIG(stat_u) == SIGPIPE)) {
        // 用户程序正常退出，或由于 SIGPIPE 退出，需要裁判程序判定
        if (WIFEXITED(stat_j)) {
            // 裁判程序正常退出
            switch (WEXITSTATUS(stat_j)) {
            case OK:
                // printf("Accepted\n");
                res = OK;
                break;
            case WRONG_ANSWER:
                // printf("Wrong answer\n");
                res = WRONG_ANSWER;
                break;
            default:
                // printf("Invalid judge exit code\n");
                res = WRONG_ANSWER;
                break;
            }
        } else if (WIFSIGNALED(stat_j) && WTERMSIG(stat_j) == SIGPIPE) {
            // 裁判程序由于 SIGPIPE 退出
            // printf("Wrong answer\n");
            res = WRONG_ANSWER;
        } else {
            // 裁判程序异常退出
            // printf("Judge exit abnormally\n");
            res = WRONG_ANSWER;
        }
    } else {
        // 用户程序运行时错误
        // printf("Runtime error\n");
        res = RE;
    }
    printf("%d %d %d\n", res, costTime, costMemory);

    sprintf(rm_pipe_j2u, "rm /app/InteractiveTest/j2u%d.fifo", offset);
    sprintf(rm_pipe_u2j, "rm /app/InteractiveTest/u2j%d.fifo", offset);
    system(rm_pipe_j2u);
    system(rm_pipe_u2j);
    
    return 0;
}