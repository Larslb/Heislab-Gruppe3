#include <pthread.h>
#include <stdio.h>

int i = 0; //Global variable
pthread_mutex_t mtx;

void* threadFunc1(){
	pthread_mutex_lock(&mtx);	
	for (int j = 0; j <= 1000000; j++){
		i++;
	}
	pthread_mutex_unlock(&mtx);
	return NULL;
}


void* threadFunc2(){
	pthread_mutex_lock(&mtx);
	for (int j = 0; j <= 1000000; j++){
		i--;
	}
	pthread_mutex_unlock(&mtx);
	return NULL;
}

int main(){

	pthread_mutex_init(&mtx, NULL);
	
	pthread_t thread1;
	pthread_t thread2;
		
	pthread_create(&thread1, NULL, threadFunc1, NULL);
	pthread_create(&thread2, NULL, threadFunc2, NULL);
		
	pthread_join(thread1, NULL);
	pthread_join(thread2, NULL);
	
	pthread_mutex_destroy(&mtx);	
	printf(" %d \n\n ", i);
	
	return 0;
}
