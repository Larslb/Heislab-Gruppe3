#include <pthread.h>
#include <stdio.h>

int i = 0; //Global variable

void* threadFunc1(){
	for (int j = 0; j < 1000000; j++){
		i++;
	}
	
	return NULL;
}


void* threadFunc2(){
	for (int j = 0; j < 1000001; j++){
		i--;
	}
	
	return NULL;
}

int main(){
	pthread_t thread1;
	pthread_t thread2;
	
	pthread_create(&thread1, NULL, threadFunc1, NULL);
	pthread_create(&thread2, NULL, threadFunc2, NULL);
	
	pthread_join(thread1, NULL);
	pthread_join(thread2, NULL);
	
	printf(" %d \n\n ", i);
	
	return 0;
}
