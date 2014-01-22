#include <stdlib.h>
#include "listsort.h"

static int cmpfn(element **a, element **b) {
    return (*a)->i - (*b)->i;
}

element *qsort_list(element *head) {
    int l = 0;
    int i = 0 ;
    element *t = head;
    element **holder;
    while (t) {
        l++;
        t = t->next;
    }

    holder  = calloc(l, sizeof(element *));

    t = head;
    while (t) {
        holder[i++] = t;
        t = t->next;
    }

    qsort(holder, l, sizeof(element *), &cmpfn);

    for (i=0; i < l-1; i++) {
        holder[i]->next = holder[i+1];
    }
    holder[i]->next = NULL;

    return holder[0];
}
