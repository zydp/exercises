#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <stdlib.h>
#include <signal.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <map>
#include <utility>
#include <iostream>
#include <atomic>
#include <mutex>
using namespace std;
/*funcs ..*/
void CatchSig();
void _sigHandler(int sig);
string GetBinName(const string& file);
string GetBinPath(const string& file); //default cwd
void CreateChildProcess(int argc, char** argv);
void BecomeEternal(const char* file); 
void PrintUsing(const char* argv);

/*global ..*/
std::mutex MMutex;
std::map<int, string> MChilds; /*pid, path*/
std::atomic_bool bexit(false);
std::atomic<unsigned int> waitCount(0);

/*main ..*/
int main(int argc ,char** argv)
{
	if(argc<2){
		PrintUsing(argv[0]);
		return 0;
	}
#ifdef DEBUG
	daemon(true, true);
#else
	daemon(true, false);
#endif
	CatchSig();
	CreateChildProcess(argc, argv);

	while(true){
		sleep(5);
		{
			if(bexit && (0==waitCount)){
				break;
			}
		}
	}
	return 0;
}

/*------------------func------------------*/




void CreateChildProcess(int argc, char** argv)
{
	for(int index = 1; index<argc; ++index){
		fprintf(stderr, "[MASTER] Ready to start child process '%s'.\n", argv[index]);
		BecomeEternal(argv[index]);
	}
}

void BecomeEternal(const char* file)
{
	pid_t pid;
	if ((pid = fork()) == 0){
		execlp(GetBinPath(file).c_str(), GetBinName(file).c_str(), NULL);
		fprintf(stderr, "[MASTER] '%s' start msg:%s\n", file,  strerror(errno));
		exit(-10);
	}
	MChilds[pid] = file;
}


void CatchSig()
{
	sigset(SIGINT, _sigHandler);
	sigset(SIGQUIT, _sigHandler);
	sigset(SIGTERM, _sigHandler);
	sigset(SIGHUP, _sigHandler);
	sigset(SIGCHLD, _sigHandler);
	sigset(SIGUSR1, _sigHandler);
	sigset(SIGUSR2, _sigHandler);
}

void _sigHandler(int sig)
{
	switch(sig)
	{
		case SIGCHLD:
		{
			std::lock_guard<std::mutex> lock(MMutex);
			int exitStatus = 0;
			pid_t child = -1;
			while((child = wait(&exitStatus))!=-1){
				fprintf(stderr, "[MASTER] child:%d exit, stauts:%d\n", child, WEXITSTATUS(exitStatus));
				if(!bexit){
					string file = MChilds[child];
					if(0==WEXITSTATUS(exitStatus)&&!file.empty()) {
						fprintf(stderr, "[MASTER] ready to restart child '%s'\n", file.c_str());
						BecomeEternal(file.c_str());
					}else if(246==WEXITSTATUS(exitStatus)){ //246 == exit(-10);
						fprintf(stderr, "[MASTER] gosh! the child can't alive, i'm [%d] going to die with he\n", getpid());
						exit(0);
					}
				}else{
					waitCount-=1;
				}
				MChilds.erase(child);
			}
		}
		break;
		case SIGTERM:
		{
			std::lock_guard<std::mutex> lock(MMutex);
			fprintf(stderr, "[MASTER] got signal %d, we will kill all childrens\n", sig);
			for(auto child = MChilds.begin(); child!=MChilds.end(); ++child){
				fprintf(stderr, "[MASTER] kill child %d : '%s'\n", child->first, child->second.c_str());
				waitCount+=1;
				bexit = true;
				kill(child->first, SIGKILL);		
			}
		}
		break;
		case SIGUSR1:
		case SIGUSR2:
		case SIGINT:
		case SIGQUIT:
		case SIGHUP:
		default:
			fprintf(stderr, "[MASTER] got %d signal, ignored\n", sig);
		break;
	}
}

string GetBinName(const string& file)
{
	auto pos = file.rfind('/');
	if(pos==std::string::npos){
		return file;
	}
	return file.substr(pos+1);
}
string GetBinPath(const string& file)
{
	auto pos = file.rfind('/');
	if(pos==std::string::npos){
		char buf[512] ={0};
		getcwd(buf, 512);
		return buf+string("/")+file;
	}
	//return file.substr(0, pos+1);
	return file;
}

void PrintUsing(const char* argv){
	fprintf(stderr, "\n-------------------------------------------\n");
	fprintf(stderr, "Did you forget something?\n");
	fprintf(stderr, "Using:\n");
	fprintf(stderr, "    %s a.out /user/b.out ./c.out ...\n\n", argv);
}
