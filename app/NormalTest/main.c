#include "../runtimeUtil/runtime.h"

void runUser(int offset) {
    char input_path[50];
    char output_path[50];
    char user_code[50] = "/app/NormalTest/user";
    
    sprintf(input_path, "/app/NormalTest/input%d", offset);
    sprintf(output_path, "/app/NormalTest/output%d", offset);
    freopen(input_path, "r", stdin);
    freopen(output_path, "w", stdout);
    
    setup_sandbox();
    
    execl(user_code, user_code, NULL);
}

void runJudge(int offset) {
    char judge_code[50];
    
    sprintf(judge_code, "/app/NormalTest/judge%d", offset);
    
    execl(judge_code, judge_code, NULL);
}


int main(int args, char **argc){
    int offset = atoi(argc[1]), time = atoi(argc[2]), memory = atoi(argc[3]);
    int res = 0, costTime = 0, costMemory = 0, tottime = 0;
    int pid_u = fork();
    if(pid_u < 0) {
        printf("-1 0 0\n");
        return 0;
    } else if(pid_u == 0) {
        runUser(offset);
        return 0;
    }

	int stat_u, stat_j;
	while(1) {
		if(waitpid(pid_u, &stat_u, WNOHANG)) {
			break;
		}
		load(&costMemory, get_proc_mem(pid_u));
        load(&costTime, get_pro_cpu_time(pid_u));
		if(tottime > time + 100 || costTime > time) {
			kill(pid_u, SIGKILL);
			waitpid(pid_u, &stat_u, 0);

            res = TLE;
            costTime = time;

			printf("%d %d %d\n", res, costTime, costMemory);
			return 0;
		}
		if(costMemory > memory) {
			kill(pid_u, SIGKILL);
			waitpid(pid_u, &stat_u, 0);
            
            res = MLE;
			
            printf("%d %d %d\n", res, costTime, costMemory);
			return 0;
		}
		usleep(40000);
		tottime += 4;
	}

    int pid_j = fork();
    if(pid_j < 0){
        printf("-1 0 0\n");
        return 0;
    }else if(pid_j == 0) {
        runJudge(offset);
        return 0;
    }

    waitpid(pid_j, &stat_j, 0);
	
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
    return 0;

}