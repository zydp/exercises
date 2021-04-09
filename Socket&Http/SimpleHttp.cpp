#include "SimpleHttp.h"
#include <string.h>
#include <thread>
namespace Exercises {
	namespace SimpleHttp {
		constexpr auto RECV_BUF_INIT_SIZE = 512;
		/* Init global environment */
		static CURLcode _init_r = curl_global_init(CURL_GLOBAL_DEFAULT);	
		/*curl call back*/
		unsigned int user_header_write(char* pContent, size_t size, size_t nmemb, void* pUsrData);
		unsigned int user_data_write(char* pContent, size_t size, size_t nmemb, void* pUsrData);
		
		/*tools function*/
		bool beautify_header_data(Response& resp, char* buf);
		char* Split(char *s, const char *delim, char **save_ptr);


		/* ------------------Request------------------ */
		void Request::_impl_init_() {
			if(isInit_.load()){
				return;
			}
			
			/*handle check*/
			pHandle_ = curl_easy_init();
			curl_easy_reset(pHandle_);
			
			/* set the method */
			switch (Method) {
				case PUT:	curl_easy_setopt(pHandle_, CURLOPT_PUT, 1L); 	break;	/*set the body*/
				case POST:	curl_easy_setopt(pHandle_, CURLOPT_POST, 1L); 	break;	/*set the body*/
				case DEL:	curl_easy_setopt(pHandle_, CURLOPT_CUSTOMREQUEST, "DELETE"); break;
				case GET:
				default:
					curl_easy_setopt(pHandle_, CURLOPT_HTTPGET, 1L);/* default */
				break;
			}
			/* set url */
			curl_easy_setopt(pHandle_, CURLOPT_URL, Url.data());

			/* set time out */
			curl_easy_setopt(pHandle_, CURLOPT_CONNECTTIMEOUT, TimeOut);
			if(filePath_.empty()){
				curl_easy_setopt(pHandle_, CURLOPT_TIMEOUT, TimeOut*6);
			}

			/* prevent CURLRESOLV_TIMEOUT*/
			curl_easy_setopt(pHandle_, CURLOPT_NOSIGNAL, 1L);

			/* set callback func*/
			curl_easy_setopt(pHandle_, CURLOPT_READFUNCTION, nullptr);
			curl_easy_setopt(pHandle_, CURLOPT_WRITEFUNCTION, user_data_write);


			/* header read function */
			curl_easy_setopt(pHandle_, CURLOPT_HEADERFUNCTION, user_data_write);
			//curl_easy_setopt(pHandle_, CURLOPT_HEADERFUNCTION, user_header_write);
			
			/* response buf */
			if(filePath_.empty()) {
				dataChunk_.memory = (char*)malloc((sizeof(char) * RECV_BUF_INIT_SIZE));
				memset(dataChunk_.memory, 0, RECV_BUF_INIT_SIZE);
				dataChunk_.cap = RECV_BUF_INIT_SIZE;
				curl_easy_setopt(pHandle_, CURLOPT_WRITEDATA, &dataChunk_);
			}
			
			/* pass in custom data to the callback */
			headerChunk_.memory = (char*)malloc((sizeof(char) * RECV_BUF_INIT_SIZE));
			memset(headerChunk_.memory, 0, RECV_BUF_INIT_SIZE/2);
			headerChunk_.cap = RECV_BUF_INIT_SIZE;
			curl_easy_setopt(pHandle_, CURLOPT_HEADERDATA, &headerChunk_);
			
			/* using Cert */
			if (!CertPath.empty()) {
				curl_easy_setopt(pHandle_, CURLOPT_SSL_VERIFYPEER, 1L); //in check
				//curl_easy_setopt(pHandle_, CURLOPT_SSL_VERIFYHOST, 0L);
				//curl_easy_setopt(pHandle_, CURLOPT_SSL_VERIFYPEER, 0L);
				//curl_easy_setopt(pHandle_, CURLOPT_CAPATH, CertPath.data());
				curl_easy_setopt(pHandle_, CURLOPT_CAINFO, CertPath.data()); /* "/etc/certs/cabundle.pem" */
			}
			
			/* init to NULL is important */
			pHeader_ = curl_slist_append(pHeader_, "Accept: application/json");
			//pHeader_ = curl_slist_append(pHeader_, "Content-Type: application/json");
			pHeader_ = curl_slist_append(pHeader_, "Content-Type: application/x-www-form-urlencoded");
			
			/* pass our list of custom made headers */
			curl_easy_setopt(pHandle_, CURLOPT_HTTPHEADER, pHeader_);
			isInit_.store(true);
		}

		Request::Request(std::string Url){
			this->Url = Url;
		}
		Request::Request(std::string Url, METHODS Method){
			this->Url = Url;
			this->Method = Method;
		}
		Request::Request(std::string Url, std::string CertPath){
			this->Url = Url;
			this->CertPath = CertPath;
		}
		Request::Request(std::string Url, unsigned short TimeOut){
			this->Url = Url;
			this->TimeOut = TimeOut;
		}
		Request::Request(std::string Url, METHODS Method, unsigned short TimeOut){
			this->Url = Url;
			this->Method = Method;
			this->TimeOut = TimeOut;
		}
		Request::Request(std::string Url, std::string CertPath, unsigned short TimeOut){
			this->Url = Url;
			this->CertPath = CertPath;
			this->TimeOut = TimeOut;
		}
		Request::Request(std::string Url, std::string CertPath, METHODS Method, unsigned short TimeOut){
			this->Url = Url;
			this->CertPath = CertPath;
			this->Method = Method;
			this->TimeOut = TimeOut;
		}

		Request::~Request() {
			if (nullptr != pHeader_) {
				/* free the header list */
				curl_slist_free_all(pHeader_), pHeader_ = nullptr;
			}
			if (nullptr != pHandle_) {
				curl_easy_cleanup(pHandle_), pHandle_ = nullptr;
			}
			/*free up*/
			if (nullptr != dataChunk_.memory) {
				free(dataChunk_.memory), dataChunk_.memory = nullptr;
			}
			if (nullptr != headerChunk_.memory) {
				free(headerChunk_.memory), headerChunk_.memory = nullptr;
			}
		}

		Request& Request::HeaderAdd(std::string key, std::string value) {
			this->pHeader_ = curl_slist_append(this->pHeader_, (key + ": " + value).c_str());
			return *this;
		}
		

		/*curl call back*/
		unsigned int user_header_write(char* pContent, size_t size, size_t nmemb, void* pUsrData){
			if(nullptr==pUsrData)return -1;
			((std::string*)pUsrData)->append(pContent, size * nmemb);
			return size * nmemb;
			/*
			size_t iRealSize = size * nmemb;
			std::string* pResp = NULL;
			pResp = dynamic_cast<std::string*>((std::string*)pUsrData); //in IoT System , probably crashed
			if (!pResp)
			{
			   return -1;
			}
			*pResp += std::string(pContent, iRealSize);
			return iRealSize;
			*/
		}

		unsigned int user_data_write(char* pContent, size_t size, size_t nmemb, void* pUsrData){
			if(nullptr==pUsrData)return -1;
			_MemoryStruct* mem = (_MemoryStruct*)pUsrData;
			size_t realsize = size * nmemb;
			if(nullptr != mem->memory){
				if ((realsize + mem->size) > mem->cap) {
					char* ptr = (char*)realloc(mem->memory, mem->size + realsize + 1);
					if (nullptr == ptr) { return 0;  /* out of memory! */ }
					mem->memory = ptr;
					mem->cap = mem->size + realsize + 1;
				}
				memcpy(&(mem->memory[mem->size]), pContent, realsize);
				mem->size += realsize;
				mem->memory[mem->size] = 0;
				return realsize;
			}else	if(nullptr != mem->fd){
				mem->size += fwrite(pContent, size, nmemb, mem->fd);
				fflush(mem->fd);
				return realsize;
			}
			return -2;
		}
		
		/*tools function*/
		bool beautify_header_ptrdata(std::shared_ptr<Response> resp, char* buf){
			if(nullptr==buf)return false;
			char* index = nullptr;
			for(char* subStr= Split(buf, "\r\n", &index); subStr!=nullptr; subStr= Split(nullptr, "\r\n", &index)){
				std::string tmp(subStr);
				int pos = tmp.find_first_of(':');
				if(std::string::npos==pos){
					continue;
				}
				resp->Headers[tmp.substr(0, pos)] = tmp.substr(pos+2);
			}
			return true;
		}
		bool beautify_header_data(Response& resp, char* buf){
			return true;//add by daiping
			if(nullptr==buf)return false;
			char* index = nullptr;
			for(char* subStr= Split(buf, "\r\n", &index); subStr!=nullptr; subStr= Split(nullptr, "\r\n", &index)){
				std::string tmp(subStr);
				int pos = tmp.find_first_of(':');
				if(std::string::npos==pos){
					continue;
				}
				resp.Headers[tmp.substr(0, pos)] = tmp.substr(pos+2);
			}
			return true;
		}

		char* Split(char *s, const char *delim, char **save_ptr)
		{
		   char *token;
			if (s == NULL) s = *save_ptr;
			/* Scan leading delimiters.  */
			s += strspn(s, delim);
			if (*s == '\0')
				return NULL;
			/* Find the end of the token.  */
			token = s;
			s = strpbrk(token, delim);
			if (s == NULL){
				/* This token finishes the string.  */
				*save_ptr = strchr(token, '\0');
			} else {
				/* Terminate the token and make *SAVE_PTR point past it.  */
				*s = '\0';
				*save_ptr = s + 1;
			}
			return token;
		}



		/* ------------------HTTP Client------------------ */

		std::shared_ptr<Response> Client::Do(const Request& req) {
			auto& rreq = const_cast<Request&>(req);
			rreq._impl_init_();
			
			/* pass our list of custom made headers */
			curl_easy_setopt(req.pHandle_, CURLOPT_HTTPHEADER, req.pHeader_);
			
			/* send request*/
			CURLcode iCurlRet = curl_easy_perform(req.pHandle_);
			std::shared_ptr<Response> resp = std::make_shared<Response>();
			if (CURLE_OK == iCurlRet) {
				short int response_code;
				curl_easy_getinfo(req.pHandle_, CURLINFO_RESPONSE_CODE, &response_code);
				resp->Code = response_code;
				resp->Body = rreq.dataChunk_.memory;
			} else {
				resp->Code = iCurlRet;
				resp->Body = curl_easy_strerror(CURLcode(iCurlRet));
			}

			beautify_header_ptrdata(resp, rreq.headerChunk_.memory);
			return resp;
		}

		std::shared_ptr<Response> Client::Do(const Request& req,const std::string& body) {
			auto& rreq = const_cast<Request&>(req);
			rreq._impl_init_();
			
			/* pass our list of custom made headers */
			curl_easy_setopt(req.pHandle_, CURLOPT_HTTPHEADER, req.pHeader_);
			
			/* set body */
			if (!body.empty()) {
				curl_easy_setopt(req.pHandle_, CURLOPT_POSTFIELDS, body.c_str());
				curl_easy_setopt(req.pHandle_, CURLOPT_POSTFIELDSIZE, body.length());
			}
			
			/* send request*/
			CURLcode iCurlRet = curl_easy_perform(req.pHandle_);
			std::shared_ptr<Response> resp = std::make_shared<Response>();
			if (CURLE_OK == iCurlRet) {
				short int response_code;
				curl_easy_getinfo(req.pHandle_, CURLINFO_RESPONSE_CODE, &response_code);
				resp->Code = response_code;
				resp->Body = rreq.dataChunk_.memory;
			} else {
				resp->Code = iCurlRet;
				resp->Body = curl_easy_strerror(CURLcode(iCurlRet));
			}

			beautify_header_ptrdata(resp, rreq.headerChunk_.memory);
			return resp;
		}
		void Client::Do(std::shared_ptr<Request> req, RESP_CB cb) {
			auto thFunc = std::bind(&Client::_impl_do, this, placeholders::_1, placeholders::_2, placeholders::_3);
			std::thread thr(thFunc, req, cb, "");
			thr.detach();
		}
		void Client::Do(std::shared_ptr<Request> req, const std::string& body, RESP_CB cb) {
			auto thFunc = std::bind(&Client::_impl_do, this, placeholders::_1, placeholders::_2, placeholders::_3);
			std::thread thr(thFunc, req, cb, body);
			thr.detach();
		}

		void Client::_impl_do(std::shared_ptr<Request> req, RESP_CB cb, const std::string& body/* = ""*/) {
			req->_impl_init_();
			
			/* pass our list of custom made headers */
			curl_easy_setopt(req->pHandle_, CURLOPT_HTTPHEADER, req->pHeader_);
			
			/* set body */
			if (!body.empty()) {
				curl_easy_setopt(req->pHandle_, CURLOPT_POSTFIELDS, body.c_str());
				curl_easy_setopt(req->pHandle_, CURLOPT_POSTFIELDSIZE, body.length());
			}

			/* send request*/
			CURLcode iCurlRet = curl_easy_perform(req->pHandle_);
			//shared_ptr<Response> resp = std::make_shared<Response>();
			Response resp;
			if (CURLE_OK == iCurlRet) {
				short int response_code;
				curl_easy_getinfo(req->pHandle_, CURLINFO_RESPONSE_CODE, &response_code);
				resp.Code = response_code;
				if(nullptr!=req->dataChunk_.memory){
					resp.Body = req->dataChunk_.memory;
				}
			} else {
				resp.Code = iCurlRet;
				resp.Body = curl_easy_strerror(CURLcode(iCurlRet));
			}
			
			beautify_header_data(resp, req->headerChunk_.memory);

			/* return to the user namesapce */
			cb(resp);
		}

		/* Candy */
		std::shared_ptr<Response> Client::Get(const std::string Url){
			return this->Do(Request{Url: Url});
		}
		std::shared_ptr<Response> Client::Get(const std::string Url, const std::string CertPath){
			return this->Do(Request{Url: Url, CertPath: CertPath});
		}
		std::shared_ptr<Response> Client::Post(const std::string Url, const std::string body){
			return this->Do(Request{Url: Url, METHODS: POST}, body);
		}
		std::shared_ptr<Response> Client::Post(const std::string Url, const std::string CertPath, const std::string body){
			return this->Do(Request{Url: Url, CertPath: CertPath, METHODS: POST, TimeOut: 6}, body);
		}

		void Client::Download(std::shared_ptr<Request> req, const std::string FilePath, RESP_CB cb){
			req->filePath_ = FilePath;
			req->_impl_init_();
			std::thread thr(
				[](std::shared_ptr<Request> req, const std::string FilePath, RESP_CB cb){
					Response resp; 
					FILE* fhandle = fopen(FilePath.c_str(), "w");
					if(!fhandle){
						resp.Code = -1;
						resp.Body = FilePath + " -> File Open Failed, Please Check";
						cb(resp);
						return;
					}
			
					/* response buf */
					req->dataChunk_.fd = fhandle;
					curl_easy_setopt(req->pHandle_, CURLOPT_WRITEDATA, &req->dataChunk_);
					CURLcode iCurlRet = curl_easy_perform(req->pHandle_);
					
					short int response_code;
					curl_easy_getinfo(req->pHandle_, CURLINFO_RESPONSE_CODE, &response_code);
					resp.Code = response_code;
					if (CURLE_OK != iCurlRet) {
						resp.Body =  curl_easy_strerror(CURLcode(iCurlRet));
					}
					
					beautify_header_data(resp, req->headerChunk_.memory);
			
					/*close file handle*/
					req->dataChunk_.fd = nullptr;
					fclose(fhandle);

					/* return to the user namesapce */
					cb(resp);
				}, req, FilePath, cb);
			thr.detach();
		}
	}/* end of namespace  SimpleHttp */
}/* end of namespace Exercises */
