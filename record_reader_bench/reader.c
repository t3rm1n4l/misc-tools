#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>
#include <unistd.h>
#include <assert.h>
#include <google/heap-profiler.h>


#define INCR 0
#define BUFSIZE 67108864
#define BUCKETS 500000
#define BUCKET_INCR 100000

typedef struct {
    char *buf;
    size_t size;
} sized_buf;

typedef struct {
    uint8_t   op;
    uint16_t ksize;
    uint32_t vsize;
} view_file_merge_record_t;

int read_view_record(FILE *in, void **buf)
{
    uint32_t len, vlen;
    uint16_t klen;
    uint8_t op;
    view_file_merge_record_t *rec;

    /* On disk format is a bit weird, but it's compatible with what
       Erlang's file_sorter module requires. */

    if (fread(&len, sizeof(len), 1, in) != 1) {
        if (feof(in)) {
            return 0;
        } else {
            return 1;
        }
    }
    if (INCR) {
        if (fread(&op, sizeof(rec->op), 1, in) != 1) {
            return 1;
        }
    }
    if (fread(&klen, sizeof(klen), 1, in) != 1) {
        return 1;
    }

    klen = ntohs(klen);
    vlen = len - sizeof(klen) - klen;
    if (INCR) {
        vlen -= sizeof(op);
    }

    rec = (view_file_merge_record_t *) malloc(sizeof(*rec) + klen + vlen);
    if (rec == NULL) {
        return 1;
    }

    rec->op = op;
    rec->ksize = klen;
    rec->vsize = vlen;

    if (fread(((char *) rec) + sizeof(view_file_merge_record_t), klen + vlen, 1, in) != 1) {
        free(rec);
        return 1;
    }

    *buf = (void *) rec;

    return klen + vlen;
}


int main(int argc, char **argv)
{

    void **records;
    void *rec;
    int x = 1;
    int bufsize = 0;
    int record_count = BUCKETS;
    int c  = 0;
    FILE *f = fopen(argv[1], "rb");

    records = (void **) calloc(record_count, sizeof(void *));

    //bufsize = record_count * sizeof(void *);
    HeapProfilerStart("program");
    while (x) {
        x = read_view_record(f, (void **) &rec);
        records[c] = rec;
        c++;
        bufsize += x;
        if (bufsize >= BUFSIZE) {
            printf("BUFSIZE: read %d records\n", c);

            sleep(1000);
        }


        if (c == record_count) {
            record_count += BUCKET_INCR;
            //bufsize = record_count * sizeof(void *);
            records = realloc(records, record_count * sizeof(void *));
            printf("mem size %d\n", record_count * sizeof(void *));
            assert(records);
        }
    }

    HeapProfilerStop();


    printf("read %d records\n", c);


}
