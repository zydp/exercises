#include <iostream>
#include <thread>
#include <stdio.h>//getchar()
#include <string.h> //strerror
#include <sys/socket.h>
#include <sys/types.h> 
#include <netinet/in.h>
#include <arpa/inet.h>
#include <errno.h>
#include <unistd.h>//close
#include <condition_variable>
#include <mutex>
#include <atomic>
#define THREAD_COUNT 8
std::atomic_bool g_exit{false};
void thread_func(int id, void* args);

class Semaphore{
public:
	Semaphore(unsigned long count = 0) : count_(count) {}
	Semaphore(const Semaphore&) = delete;
	Semaphore& operator=(const Semaphore&) = delete;
	void Signal(){
		{
			std::unique_lock<std::mutex> lock(mutex_);
			++count_;
		}
		cv_.notify_one();
	}
	
	void Wait(){
		std::unique_lock<std::mutex> lock(mutex_);
		while (count_ == 0) {
			cv_.wait(lock);
		}
		--count_;
	}
private:
	std::mutex mutex_;
	std::condition_variable cv_;
	std::atomic_ulong count_;
};


int main(int argc, char** argv){
	std::cout << __func__ << std::endl;
	Semaphore cnd;
	std::thread ths[THREAD_COUNT];
	for(int i=0; i<THREAD_COUNT; i++){
		ths[i] = std::thread(thread_func, i, &cnd);	
	}
	
	for(int i=0; i<50000; i++){
		cnd.Signal();
	}
	
	std::cout << "exit the threads?" << std::endl;
	//getchar();
	sleep(2);
	g_exit = true;
	
	for(int i=0; i<THREAD_COUNT; i++){
		cnd.Signal();
		usleep(500);
	}
	
	for(auto& th :ths){
		th.join();
	}

	std::cout << "Ready to exit ?" << std::endl;
	getchar();
	return 0;
}


void thread_func(int id, void* args){
	Semaphore* cnd = (Semaphore*)args;
	while(true)
	{
		cnd->Wait();
		if(g_exit){
			break;
		}
		std::cout << "I'm thread " << id << std::endl;
		std::this_thread::yield();
	}
	std::cout << "Thread " << id << " exit" << std::endl;
}

