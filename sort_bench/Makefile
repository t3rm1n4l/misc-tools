all: sortdriver.o listsort.o qsort.o
	gcc -o test sortdriver.o listsort.o qsort.o

sortdriver.o: sortdriver.c
	gcc -c sortdriver.c

listsort.o: listsort.c
	gcc -c listsort.c

qsort.o: qsort.c
	gcc -c qsort.c

clean:
	rm -f *.o test
