#include <pthread.h>
#include <stdio.h>

int i = 0;

void* threadfunction_1(){
	for (int j=0; j<1000001; j++)
		i++;
	return NULL;
}

void* threadfunction_2(){
	for (int j=0; j<1000001; j++)
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

	printf("%d\n\n", i); 

	return 0;
}
