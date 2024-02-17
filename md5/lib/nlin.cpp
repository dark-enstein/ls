#include <iostream>
#include <stdio.h>
#include <cmath>
#include <cstdio>
#include "nlin.h"

using namespace std;


void nlin(int min, int max) {
	for (; min <= max; min++){
		unsigned long long compute = floor(abs(UNSIGNED_LONG * sin(min + 1)));
		printf("%d: %llu\n", min, compute);
	}
}

void pad(char string) {}

void mdbuf() {}
	
int main() {
	nlin(0, 63);
	return 0;
}
