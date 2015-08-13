#include <thread>
#include <vector>
#include <iostream>

void modify(volatile unsigned int *arr, int offset) {
    for (int i=0; i<100000000; i++) {
        arr[offset]++;
    }
}

int main(int argc, char *argv[]) {
    std::vector<std::thread> threads;
    int width=1;
    if (argc > 1) {
        width=8;
        std::cout<<"false sharing won't occur this time"<<std::endl;
    }
    volatile unsigned int array[8*8];

    for (int i=0; i < 8; i++) {
        array[i*width] = 0;
        threads.push_back(std::thread(&modify, array, i*width));
    }

    for (auto& th : threads) th.join();


    for (int i=0; i < 8; i++) {
        std::cout<<array[i*width]<<std::endl;
    }
    return 0;
}
