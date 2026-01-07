// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>

int i = 0;

// Note the return type: void*
void* incrementingThreadFunction(){
    for (int j = 0; j < 1000000; j++){
        i++;
    }
    // TODO: increment i 1_000_000 times
    return NULL;
}

void* decrementingThreadFunction(){
    for (int j = 0; j < 1000000; j++){
        i--;
    }
    // TODO: decrement i 1_000_000 times
    return NULL;
}


int main(){
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


    printf("The magic number is: %d\n", i);
    return 0;
}
