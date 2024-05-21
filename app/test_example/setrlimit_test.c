#include <stdio.h>
#include <sys/resource.h>

int main() {
    struct rlimit rl;
    
    // 获取当前的文件大小限制
    if (getrlimit(RLIMIT_FSIZE, &rl) == -1) {
        perror("getrlimit");
        return 1;
    }

    printf("Current soft limit: %llu\n", rl.rlim_cur);
    printf("Current hard limit: %llu\n", rl.rlim_max);

    // 设置新的软限制为当前硬限制
    rl.rlim_cur = rl.rlim_max;

    // 尝试设置新的限制
    if (setrlimit(RLIMIT_FSIZE, &rl) == -1) {
        perror("setrlimit");
        return 1;
    }

    // 再次获取限制以验证更改
    if (getrlimit(RLIMIT_FSIZE, &rl) == -1) {
        perror("getrlimit");
        return 1;
    }

    printf("New soft limit: %llu\n", rl.rlim_cur);
    printf("New hard limit: %llu\n", rl.rlim_max);

    return 0;
}