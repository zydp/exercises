#pragma once
#include <map>
#include <atomic>
#include <string>
#include <memory>
#include <iostream>
#include <functional>
#include <curl/curl.h>
using namespace std;

namespace Inesa {
	namespace SimpleHttp {
		/*----sweet------*/
		using HEADER_PTR = struct curl_slist*;
		using FORM = std::map<string, string>;
		using HEADLIST = std::map<string, string>;
		typedef struct {
			size_t size = 0;
			size_t cap = 0;
			char* memory = nullptr;
			FILE* fd = nullptr;
		}_MemoryStruct;
		/*========================================================================*/
		/*
		 *  Method of the HTTP request
		 */
		enum METHODS {
			GET,	/* http get */
			PUT,	/* http put */
			POST,	/* http post */
			DEL	/* http del */
		};

		/*========================================================================*/
		/*
		 * HTTP Reponse
		 */
		class Response{
		public:
			short int	Code = -1;
			std::string Body = "";
			HEADLIST	Headers;
		};

		/*========================================================================*/
		/*
		 * Callback Function
		 */
		using RESP_CB = std::function<void(const Response&)>;

		/*========================================================================*/
		/*
		 *  To build a http request.
		 */
		class Request
		{
			friend class Client;			/* let them be lovers. */
		public:
			std::string Url;
			std::string CertPath;
			METHODS Method = GET;
			unsigned short TimeOut = 6;
		public:
			Request(std::string Url);
			Request(std::string Url, METHODS Method);
			Request(std::string Url, std::string CertPath);
			Request(std::string Url, unsigned short TimeOut);
			Request(std::string Url, METHODS Method, unsigned short TimeOut);
			Request(std::string Url, std::string CertPath, unsigned short TimeOut);
			Request(std::string Url, std::string CertPath, METHODS Method, unsigned short TimeOut);
			virtual ~Request();
		public:
			Request& HeaderAdd(std::string key, std::string value);
		private:
			void _impl_init_();
			void* pHandle_ = NULL;			/* inside http handle from libcurl */
			std::string filePath_;
			HEADER_PTR pHeader_ = NULL;		/* header list */
			std::atomic<bool> isInit_{false};
			_MemoryStruct dataChunk_;	
			_MemoryStruct headerChunk_;	
		};

		/*========================================================================*/
		/*
		 * The Simple Client
		 */
		class Client
		{
		public:
			Client(){}
			virtual ~Client(){}
		public:
			std::shared_ptr<Response> Do(const Request& req);
			std::shared_ptr<Response> Do(const Request& req, const std::string& body);
			void Do(std::shared_ptr<Request> req, RESP_CB cb);
			void Do(std::shared_ptr<Request> req, const std::string& body, RESP_CB cb);
			void Download(std::shared_ptr<Request> req, const std::string FilePath, RESP_CB cb);
			/* Candy */
			std::shared_ptr<Response> Get(const std::string Url);
			std::shared_ptr<Response> Get(const std::string Url, const std::string CertPath);
			std::shared_ptr<Response> Post(const std::string Url, const std::string body);
			std::shared_ptr<Response> Post(const std::string Url, const std::string CertPath, const std::string body);
		private:
			void _impl_do(std::shared_ptr<Request> req, RESP_CB cb, const std::string& body = "" );
		};
	}/* end of namespace  SimpleHttp */
}/* end of namespace Inesa */
