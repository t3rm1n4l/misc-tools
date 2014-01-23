#include <stdio.h>
#include "reader.h"

int main(int argc, char **argv) {

    int x;
    char b[90];
    if (argv[1][0] == 'b') {
        bufreader_t *rdr = bufreader_create("test.txt");
        while (1) {
            if (bufreader_read(rdr, sizeof(int), &x) == 0) {
                break;
            }

            if (bufreader_read(rdr, sizeof(int), &x) == 0) {
                break;
            }

            if (bufreader_read(rdr, 90, b) == 0) {
                break;
            }


            printf("No is %d\n", x);

        }
    } else {
        FILE *fp = fopen("test.txt", "rb");
        while (1) {
            if (fread(&x, sizeof(int), 1, fp) != 1) {
                if (feof(fp)) {
                    break;
                }
            }

            if (fread(&x, sizeof(int), 1, fp) != 1) {
                if (feof(fp)) {
                    break;
                }
            }

            if (fread(b, 90, 1, fp) != 1) {
                if (feof(fp)) {
                    break;
                }
            }
            printf("No is %d\n", x);
        }
    }

}
