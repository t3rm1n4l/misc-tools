all: test.o test_driver gen_data reader.o

reader.o: reader.c
	gcc -O3 -c reader.c

gen_data: gen_data.c
	gcc -O3 -o gen_data gen_Data.c

test_driver: test.o reader.o
	gcc -O3 -o test test.o reader.o

test.o: test.c
	gcc -O3 -c test.c

clean:
	rm -rf *.o
