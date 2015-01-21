#include <pthread.h>
#include <stdio.h>

int i = 0;

void* threadfunction_1(){
	i++;
	return NULL;
}

void* threadfunction_2(){
	i--;
	return NULL;
}

int main(){
	pthread_t Thread_1;
	pthread_t Thread_2;

	pthread_create(&Thread_1, NULL, threadfunction_1(), NULL);
	pthread_create(&Thread_2, NULL, threadfunction_2(), NULL);

	pthread_join(Thread_1, NULL);
	pthread_join(Thread_2, NULL); 

	return 0;
}
