// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>

pthread_mutex_t i_lock;
int i = 0;

// Note the return type: void*
void* incrementingThreadFunction(){
    for (int j = 0; j < 1000000; j++){
        pthread_mutex_lock(&i_lock);
        i++;
        pthread_mutex_unlock(&i_lock);
    }
    return NULL;
}

void* decrementingThreadFunction(){
    for (int j = 0; j < 1000000; j++){
        pthread_mutex_lock(&i_lock);
        i--;
        pthread_mutex_unlock(&i_lock);
    }
    return NULL;
}


int main(){
    pthread_mutex_init(&i_lock, NULL);

    // TODO: 
    // start the two functions as their own threads using `pthread_create`
    // Hint: search the web! Maybe try "pthread_create example"?
    pthread_t thread_1, thread_2;
    int return_code_1, return_code_2;

    return_code_1 = pthread_create(&thread_1, NULL, incrementingThreadFunction, NULL);
    if (return_code_1 != 0) {
        perror("pthread_create() error for thread 0");
    }

    return_code_2 = pthread_create(&thread_2, NULL, decrementingThreadFunction, NULL);
    if (return_code_2 != 0) {
        perror("pthread_create() error for thread 1");
    }
    
    // TODO:
    // wait for the two threads to be done before printing the final result
    // Hint: Use `pthread_join`
    return_code_1 = pthread_join(thread_1, NULL);
    if (return_code_1 != 0) {
        perror("pthread_join() error");
        // exit(EXIT_FAILURE);
    }
    return_code_2 = pthread_join(thread_2, NULL);
    if (return_code_2 != 0) {
        perror("pthread_join() error");
        // exit(EXIT_FAILURE);
    }

    pthread_mutex_destroy(&i_lock);

    printf("The magic number is: %d\n", i);
    return 0;
}

// Use mutex instead of semaphore because we are only
// controlling access to one resource (the global variable i),
// and the one locking the resource should only be the one
// allowed to unlock it.