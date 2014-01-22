#include <stdio.h>
#include <stdlib.h>
#include "listsort.h"
#include <time.h>

static element *gen_workload(size_t n) {
    int i = 0;
    srand(time(NULL));
    element *head = NULL, *tmp, *last;
    for (i = 0; i < n; i++) {
        tmp = calloc(1, sizeof(element));
        tmp->i = rand() % n;
        if (head) {
            last->next = tmp;
        } else {
            head = tmp;
        }

        last = tmp;
    }

    return head;
}


static void print(element *head) {

    while (head) {
        printf("%d\n", head->i);
        head = head->next;
    }
}


int main(int argc, char **argv) {

    if (argc != 3) {
        fprintf(stderr, "Usage: %s count q/l\n", argv[0]);
        exit(1);
    }

    element *in, *out;
    in = gen_workload(atoi(argv[1]));
    if (argv[2][0] == 'l') {
        out = listsort(in, 0, 0);
    } else {
        out = qsort_list(in);
    }

    print(out);

    return 0;
}
