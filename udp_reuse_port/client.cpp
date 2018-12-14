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

#define BUF_SIZE 256
#define RECV_PORT 4432
#define RECV_ADDR "192.168.6.109"
//#define RECV_ADDR INADDR_ANY



int main(int argc, char** argv){
	std::cout << __func__ << std::endl;
	int fd = socket(AF_INET, SOCK_DGRAM, 0);
	if(fd<0){
		perror("cli socket:");  
		exit(-1);
	}
	struct sockaddr_in ser_addr;
	memset(&ser_addr, 0, sizeof(sockaddr_in));
	ser_addr.sin_family = AF_INET;
	ser_addr.sin_port = htons(RECV_PORT);
	if(inet_pton(AF_INET, RECV_ADDR, &ser_addr.sin_addr) == -1){
		printf("inet_pton error for %s\n", RECV_ADDR);
		exit(-2);
	}

	int len = sizeof(ser_addr);
	const char* buf = "drift_wonder_intelligence_diva";
	printf("[%s]\n", buf);
	short send_count = 5;
	while(send_count--){
		int send_len = sendto(fd, buf, strlen(buf), 0, (struct sockaddr*)&ser_addr, len);
		if(send_len<0){
			perror("cli sendto:");
		}else{
			printf("send msg successful [%d]\n", send_len);
		}
		sleep(1);
	}
	close(fd);
	std::cout << "Ready to exit ?" << std::endl;
	getchar();
	return 0;
}

