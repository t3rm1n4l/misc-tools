#include <iostream>
#include <sstream>
#include <vector>
#include <iomanip>
using namespace std;

int main() {
    int N = 10000000;
    vector<string> items(N);
    int i;
    time_t start, end;

    start = time(NULL);
    for (i=0; i<N; i++) {
        stringstream ss;
        ss<<"document_"<<setfill('0')<<setw(7)<<i;
        items[i] = ss.str();
    }
    end = time(NULL);
    cout<<"Time taken:"<<end-start<<endl;

}
