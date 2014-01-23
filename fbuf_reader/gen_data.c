#include <stdio.h>
#include <stdlib.h>

int main(int argc, char **argv) {
    int i, N;
    FILE *fp = fopen("test.txt", "wb");
    N = atoi(argv[1]);

    for (i = 0; i < N; i++) {
        fwrite(&i, sizeof(int), 1, fp);
    }

    fclose(fp);
}
