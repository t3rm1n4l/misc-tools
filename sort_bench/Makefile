all: sortdriver.o listsort.o qsort.o
	gcc -O3 -o test sortdriver.o listsort.o qsort.o

sortdriver.o: sortdriver.c
	gcc -O3 -c sortdriver.c

listsort.o: listsort.c
	gcc -O3 -c listsort.c

qsort.o: qsort.c
	gcc -O3 -c qsort.c

clean:
	rm -f *.o test
