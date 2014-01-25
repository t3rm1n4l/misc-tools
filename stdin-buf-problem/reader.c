#include <stdio.h>

int main() {

	char buf[4096];
	int l;
	int h, i1, i2;

	fgets(buf, sizeof(buf), stdin);
	l = strlen(buf);
	buf[l-1] = '\0';

	printf("got line %s\n", buf);
	fflush(stdout);

	fscanf(stdin, "%d\n", &h);
	printf("got h %d\n", h);
	fflush(stdout);

	fscanf(stdin, "%d\n", &i1);
	printf("got i1 %d\n", i1);
	fflush(stdout);

	fscanf(stdin, "%d\n", &i2);
	printf("got i2 %d\n", i2);
	fflush(stdout);

}
