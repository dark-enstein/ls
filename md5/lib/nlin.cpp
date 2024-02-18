#include <iostream>
#include <string>
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

void pad(string secret) {
	int secret_len = secret.length();
	char* bank = new char[secret_len + 1];

	strcpy(bank, secret.c_str());
	printf("original secret byte: %s\n", bank);
	// get the mod rem of the length of input string in bits
	
	char* original_bit = bin_str(secret);
	printf("original bit len: %s\n", original_bit);
	
	int original_bitlen = secret_len * 8;
	printf("original bit len: %d\n", original_bitlen);

	int rem = MODULO_BITS - original_bitlen % MODULO_BITS;
	printf("modulo reminder; %d - %dmod%d: %d\n", MODULO_BITS, original_bitlen, MODULO_BITS, rem);
	
	int constant_padding = rem - 64;
	printf("constant padding: %d\n", constant_padding);
	// add 1 bit
	// pad with zeros until the reminder remains 64 bits == 8 bytes
	// return

}

void bitm(int shift) {
	int with = ADD_B;
	
	// print original value without conversion
	printf("%d\n", with);

	// convert original value into binary
	const char* binary = bin(with);
	printf("%s\n", binary);

	// bit shift the original value by shift
	with >>= shift;
		
	printf("%d\n", with);
	binary = bin(with);
	printf("%s\n", binary);
	
	// cleanup
	delete[] binary;
}

char* bin(int x) {

	int bitsize = sizeof(x);
	string bundle = "";
	for (int i = bitsize; i >= 0; --i){
		int bbit = ((x >> i) & 1);
		bundle += to_string(bbit);
	}

	char* binary = new char[bundle.length() + 1];

	strcpy(binary, bundle.c_str());

	return binary;
}

// handle converting a string into its binary representation
char* bin_str(string c) {
	char root [c.length() * 8];	
	// iterate over string, converting every char into binary, then appending that to the return string
	// const char* rune = c.c_str();

	for (int i = 0; i < c.length(); i++) {
		strcat(root, bin(c[i]));
		printf("g: %s\n", root);
	}

	char* root_ptr = &(root[0]);

	return root_ptr;
}

void mdbuf() {}
	
int main() {
	// core compute
	nlin(0, 4);
	// bitshifting
	bitm(1);

	string secret = "hkbdfsdsed";
	pad(secret);
	
	return 0;
}
