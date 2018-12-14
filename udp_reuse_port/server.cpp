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

#define THREAD_COUNT 3
#define BUF_SIZE 256
#define RECV_PORT 4432
#define RECV_ADDR "192.168.6.109"
//#define RECV_ADDR INADDR_ANY

void thread_func(int id);


int main(int argc, char** argv){
	std::cout << __func__ << std::endl;
	
	std::thread ths[THREAD_COUNT];
	for(int i=0; i<THREAD_COUNT; i++){
		ths[i] = std::thread(thread_func, i);	
	}
	for(auto& th :ths){
		th.join();
	}

	std::cout << "Ready to exit ?" << std::endl;
	getchar();
	return 0;
}


void thread_func(int id){
	std::cout << "I'm thread " << id << std::endl;
	int fd = socket(AF_INET, SOCK_DGRAM, 0);
	if(fd < 0){
		fprintf(stderr, "Thread[%d] socket error:%s\n", id, strerror(errno));
		return;
	}
	int len = -1;
	struct sockaddr_in ser_addr;
	memset(&ser_addr, 0, sizeof(sockaddr_in));
	ser_addr.sin_family = AF_INET;
	ser_addr.sin_port = htons(RECV_PORT);
	//ser_addr.sin_addr.s_addr = htonl(RECV_ADDR);
	if(inet_pton(AF_INET, RECV_ADDR, &ser_addr.sin_addr) == -1){
		printf("inet_pton error for[%s], errmsg[%s]\n", RECV_ADDR, strerror(errno));
		exit(-2);
	}
	len = sizeof(ser_addr);

	int opt = 1;
	if(setsockopt(fd, SOL_SOCKET, SO_REUSEPORT, &opt, sizeof(opt))<0){
		fprintf(stderr, "Thread[%d] setsockopt SO_REUSEPORT error:%s\n", id, strerror(errno));
		return;
	}
	
	opt = 1;
	if(setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt))<0){
		fprintf(stderr, "Thread[%d] setsockopt SO_REUSEADDR error:%s\n", id, strerror(errno));
		return;
	}
	
	if(bind(fd, (struct sockaddr*)&ser_addr, sizeof(ser_addr))){
		fprintf(stderr, "Thread[%d] bind error:%s\n", id, strerror(errno));
		return;
	}
	
	int recv_count = BUF_SIZE;
	char buf[BUF_SIZE] = {};
	struct sockaddr_in cli_addr;
	int cli_addr_len = sizeof(cli_addr);
	while(recv_count--){
		memset(&cli_addr, 0, cli_addr_len);
		memset(buf, 0, BUF_SIZE);
		int recv_len = recvfrom(fd, buf, BUF_SIZE-1, 0, (struct sockaddr*)&cli_addr,(socklen_t*)&cli_addr_len);
		if(recv_len < 0){
			fprintf(stderr, "Thread[%d] recvfrom error:%s\n", id, strerror(errno));
			continue;
		}else{
			fprintf(stdout, "Thread[%d] recv a msg:[%s][%s]\n", id,  inet_ntoa(cli_addr.sin_addr), buf);
		}
		
	}
	close(fd);
}

