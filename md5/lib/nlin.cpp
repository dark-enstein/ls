#include <iostream>
#include <stdio.h>
#include <cmath>
#include <cstdio>
#include "nlin.h"

using namespace std;


void nlin(int min, int max) {
	for (; min <= max; min++){
		unsigned long long compute = floor(abs(UNSIGNED_LONG * sin(min + 1)));
	//	printf("%d: %llu\n", min, compute);
	}
}

void pad(char string) {}

void bitm(int shift) {
	int with = ADD_B;
	
	printf("%d\n", with);

	printf("%s\n", bin(with));
	
	with >>= shift;
		
	printf("%s\n", bin(with));
}

char* bin(int x) {

	int bitsize = sizeof(x) * 8;
	string bundle = "";
	for (int i = bitsize; i >= 0; --i){
		int bbit = ((ADD_B >> i) & 1);
		bundle += to_string(bbit);
	}

	printf("Binary: %s\n", bundle.c_str());

	char* binary = *bundle.c_str();

	return binary;
}

void mdbuf() {}
	
int main() {
	// nlin(0, 63);
	bitm(0);
	return 0;
}
