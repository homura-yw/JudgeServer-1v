#include<bits/stdc++.h>

using namespace std;

int main() {
    ios::sync_with_stdio(false);
    cin.tie(0), cout.tie(0);
    int l = 0, r = 1e8, ans = 0;
    while(l <= r) {
        int mid = (l + r) / 2;
        cout << mid << endl;
        int op; cin >> op;
        if(op == 0) return 0; 
        if(op == 1) r = mid - 1, ans = mid;
        else l = mid + 1;
    }
    return 0;
}