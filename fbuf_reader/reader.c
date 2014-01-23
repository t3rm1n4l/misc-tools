#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <fcntl.h>
#include "reader.h"

bufreader_t *bufreader_create(char *filename)
{
    bufreader_t *buf = calloc(1, sizeof(bufreader_t));
    buf->fd = open(filename, O_RDONLY);
    if (buf->fd < 0) {
        free(buf);
        return NULL;
    }

    buf->end = -1;

    return buf;
}


size_t bufreader_read(bufreader_t *buf, size_t len, char *out)
{
    int n, buffered, remaining;
    assert(buf);

    buffered = buf->end - buf->pos + 1;

    if (buffered < len) {
        if (buf->iseof) {
            return 0;
        }
        memmove(buf->buf, buf->buf + buf->pos, buffered);
        buf->pos = 0;
        buf->end = buffered - 1;
        remaining = BUF_SIZE - buffered;
        n = read(buf->fd, buf->buf + buf->end + 1, remaining);
        if (n == 0) {
                buf->iseof = 1;
        }

        buf->end += n;
    }

    buffered = buf->end - buf->pos + 1;
    if (buffered < len) {
        return -1;
    }

    memcpy(out, buf->buf + buf->pos, len);
    buf->pos += len;

    return 1;
}
