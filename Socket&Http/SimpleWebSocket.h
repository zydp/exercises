#pragma once
#include <map>
#include <atomic>
#include <string>
#include <memory>
#include <string.h>
#include <signal.h>
#include <iostream>
#include <functional>
#include <ixwebsocket/IXNetSystem.h>
#include <ixwebsocket/IXWebSocket.h>
#include <ixwebsocket/IXWebSocketHttpHeaders.h>

namespace Inesa {
	namespace SimpleWebSocket {
		
		/*
		 * Client Call Back
		 * Implement the relevant callback function
		 */
		class WSCallBack{
		public:
			WSCallBack(){}
			virtual ~WSCallBack(){}
		public:
			virtual void OnOpen(const std::map<std::string, std::string>& headers, const std::string& uri) = 0;
			virtual void OnMessage(const std::map<std::string, std::string>& headers, const std::string& msg) = 0;
			virtual void OnError(const int status, const std::string& reason) = 0;
			virtual void OnClose(const int& code, const std::string& reason) = 0;
		};


		/* Simple WebSocket Client */
		class Client
		{
		public:
			/* Candy */
			Client(std::string ip, unsigned short port);
			Client(std::string ip, unsigned short port, std::string path);
			Client(std::string ip, unsigned short port, std::string path, std:: string certPath);
			virtual ~Client();
		public:
			Client& SetCallBack(std::shared_ptr<WSCallBack> cb);
			Client& HeaderAdd(std::string key, std::string val);
			Client& SetHeartBeatInterval(unsigned short interval);
			void Connect(unsigned timeout = 6);	//unit: s
			bool AutoReconnect(bool r, unsigned short interval); /* unit: s, return the auto state after set */
			size_t SendText(std::string msg, bool binary = false);
			void DisConnect();
			std::string const GetServerAddress();
		private:
			bool _impl_init_();
			std::shared_ptr<WSCallBack> cb_;
			std::string address_, ip_, uri_, certPath_;
			unsigned short port_;
			ix::WebSocketHttpHeaders headers_;
			ix::WebSocket wsHandle_;
			std::atomic<bool> isInit_{false};
			std::atomic<bool> isConnected_{false};
		};
	}/* end of namespace SimpleWebSocket */
}/* end of namespace Inesa */

