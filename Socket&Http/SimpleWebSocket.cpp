#include "SimpleWebSocket.h"
#include <time.h>

namespace Inesa {
   namespace SimpleWebSocket {
		/* Candy */
		Client::Client(std::string ip, unsigned short port){
			ip_ = ip;
			port_ = port;
			_impl_init_();
		}
		Client::Client(std::string ip, unsigned short port, std::string path){
			ip_ = ip;
			port_ = port;
			uri_ = path;
			_impl_init_();
		}
		Client::Client(std::string ip, unsigned short port, std::string path, std:: string certPath){
			ip_ = ip;
			port_ = port;
			uri_ = path;
			certPath_ = certPath;
			ix::SocketTLSOptions tlsOpt;
			tlsOpt.caFile = certPath_;
			tlsOpt.tls = true;
			wsHandle_.setTLSOptions(tlsOpt);
			/*
			wsHandle_.setTLSOptions({
					.certFile = "",
					.keyFile = "",
					.caFile = certPath_, // as a file, or in memory buffer in PEM format
					.tls = true 			// required in server mode
				});
			*/
			_impl_init_();
		}
			
		bool Client::_impl_init_(){
			ix::initNetSystem();	//for windows
			/* header init */
			wsHandle_.addSubProtocol("appProtocol-v1");
			wsHandle_.addSubProtocol("appProtocol-v2");
			headers_["Origin"] = (certPath_.empty()?"http://":"https://") + ip_ + ":" + std::to_string(port_);
			address_ = (certPath_.empty()?"ws://":"wss://") +ip_ + ":" + std::to_string(port_) + uri_;
			headers_["User-Agent"] = "INESA";
			wsHandle_.setUrl(address_);

			/* call back*/
			wsHandle_.setOnMessageCallback([this](const ix::WebSocketMessagePtr& msg){	
					std::map<std::string, std::string> h;
					for(auto it: msg->openInfo.headers){
						h[it.first] = it.second;
					}
					switch(msg->type)
					{
						case ix::WebSocketMessageType::Open:
							isConnected_.store(true);
							cb_->OnOpen(h, msg->openInfo.uri);
							break;
						case ix::WebSocketMessageType::Close:
							isConnected_.store(false);
							cb_->OnClose(msg->closeInfo.code, msg->closeInfo.reason);
							break;
						case ix::WebSocketMessageType::Error:
							cb_->OnError(msg->errorInfo.http_status,  msg->errorInfo.reason);
							break;
						case ix::WebSocketMessageType::Message:
							cb_->OnMessage(h, msg->str);
							break;
						case ix::WebSocketMessageType::Ping:
						case ix::WebSocketMessageType::Pong:

							break;
						default:
							//we don't care about it.
							break;
					}
				}
			);
			return true;
		}
		Client::~Client(){
			 ix::uninitNetSystem();//for windows
		}
		Client& Client::SetCallBack(std::shared_ptr<WSCallBack> cb){
			cb_ = cb;
			return *this;
		}
		Client& Client::HeaderAdd(std::string key, std::string val){
			headers_[key] = val;
			return *this;
		}
		Client& Client::SetHeartBeatInterval(unsigned short interval){/* keep alive */
			wsHandle_.setPingInterval(interval);
			return *this;
		}
		void Client::Connect(unsigned timeout/* = 6*/){//unit: s
			if(isConnected_.load())return;
			wsHandle_.setHandshakeTimeout(timeout);
			wsHandle_.setExtraHeaders(headers_);
			wsHandle_.start();
		}
		bool Client::AutoReconnect(bool r, unsigned short interval){ /* unit: s, return the auto state after set */
			r?wsHandle_.enableAutomaticReconnection(),wsHandle_.setMaxWaitBetweenReconnectionRetries(interval*1000):wsHandle_.disableAutomaticReconnection();
		}
		size_t Client::SendText(std::string msg, bool binary/* = false*/){
			if(!isConnected_.load())return -1;
			const ix::WebSocketSendInfo sendInfo = std::move(wsHandle_.send(msg, binary));
			return sendInfo.wireSize;
		}
		void Client::DisConnect(){
			if(!isConnected_.load())return;
			wsHandle_.stop();
		}
		std::string const Client::GetServerAddress(){
			return address_;
		}
   }/* end of namespace SimpleWebSocket */
}/* end of namespace Inesa */
