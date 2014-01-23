#include <stdio.h>
#include <stdlib.h>

int main(int argc, char **argv) {
    int i, N;
    char *kv = "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890";
    FILE *fp = fopen("test.txt", "wb");
    N = atoi(argv[1]);

    for (i = 0; i < N; i++) {
        fwrite(&i, sizeof(int), 1, fp);
        fwrite(&i, sizeof(int), 1, fp);
        fwrite(kv, 90, 1, fp);
    }

    fclose(fp);
}
