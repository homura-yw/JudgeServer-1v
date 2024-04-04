#include <bits/stdc++.h>

using namespace std;

int main() {
    ios::sync_with_stdio(false);
    cin.tie(0), cout.tie(0);
    int ans = 114514;
    for(int i = 1; i <= 30; ++i) {
        int op; cin >> op;
        if(op > ans) cout << 1 << endl;
        if(op < ans) cout << -1 << endl;
        if(op == ans) {
            cout << 0 << endl;
            break;
        }
    }
    return 0;
}