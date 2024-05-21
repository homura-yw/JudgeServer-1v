#include <bits/stdc++.h>

using namespace std;

int main() {
    FILE *output = fopen("/app/NormalTest/output1", "r+");
    FILE *ans = fopen("/app/NormalTest/answer1", "r+");
    int a, b; fscanf(output, "%d", &a);
    fscanf(ans, "%d", &b);
    if(a == b) return 0;
    return 1;
}
