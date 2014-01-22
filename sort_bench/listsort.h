#ifndef LISTSORT_H
#define LISTSORT_H

typedef struct element element;
struct element {
    element *next, *prev;
    int i;
};

element *listsort(element *list, int is_circular, int is_double);
element *qsort_list(element *head);
#endif
