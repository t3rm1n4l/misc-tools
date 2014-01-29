all:
	gcc -g -ltcmalloc reader.c -o reader


profile: clean
	HEAP_PROFILE_ALLOCATION_INTERVAL=1242880 HEAPPROFILE=heap.prof ./reader 1.sort

clean:
	rm -rf heap.prof*

