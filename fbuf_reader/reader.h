#ifndef BUFREADER_H
#define BUFREADER_H

#define BUF_SIZE 1048576

typedef struct {
    FILE *f;
    char buf[BUF_SIZE];
    int pos;
    int end;
    int iseof;
} bufreader_t;


bufreader_t *bufreader_create(char *filename);

size_t bufreader_read(bufreader_t *buf, size_t len, char *out);
#endif
