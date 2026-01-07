Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> Concurrency: Doing (but not executing) multiple tasks at the same time.
> Parallelism: Executing multiple tasks at the same time.

What is the difference between a *race condition* and a *data race*? 
> Race condition: Result depends on unpredictable timing of concurrent operations.
> Data race: Multiple threads attempts to access the same memory location concurrently, where alteast one is a write, and no synchronization.
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> A schedueler schedules which threads should run at any one time. It may be a cooperative scheduler where each thread much yield access by itself (compiler automatically inserts yields or the programmer must manually insert yields), and there is preemptive scheduling where the scheduler stops threads after a specific time.


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
* To do tasks faster.
* For the code's functions to be independent for more readable and easier to change code.
* For a program's GUI to not freeze while doing heavy computations.

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> Threads are managed by the OS kernel's scheduler, are preemptive scheduled, and many threads can run in parallel.
> Fibers are cooperative scheduled, managed by sofware on the OS, and can't run in parallel if all fibers are on the same thread.
> Fibers are much cheaper than threads.

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> Easier: when not sharing resources.
> Harder: when trying to access the same resource.

What do you think is best - *shared variables* or *message passing*?
> Shared variables are used under the hood (e.g. semaphores, mutexes). But we prefer message passing where 


