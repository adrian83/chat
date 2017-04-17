import 'dart:html';
import 'dart:convert';
import 'dart:async';

import 'messages.dart';

class WebSocketWrapper {
  WebSocket _socket;

  WebSocketWrapper(this._socket);

  void onOpen(fun) {
    this._socket.onOpen.listen((e) => fun(e));
  }

  void onClose(fun) {
    this._socket.onError.listen((e) => fun(e));
  }

  void onError(fun) {
    this._socket.onError.listen((e) => fun(e));
  }

  void onMessage(fun) {
    this._socket.onMessage.listen((MessageEvent e) => fun(e));
  }

  void close() {
    this._socket.close();
  }

  void send(String msg) {
    this._socket.send(msg);
  }
}

class WSClient {
  String sessionId;
  WebSocketWrapper _socket;
  bool _closedByMe;
  StreamController _onMsgCtrl = new StreamController.broadcast();
  StreamController _onOpenCtrl = new StreamController.broadcast();
  StreamController _onCloseCtrl = new StreamController.broadcast();
  StreamController _onErrorCtrl = new StreamController.broadcast();

  WSClient(this.sessionId, this._socket);

  Stream<Message> get messages => _onMsgCtrl.stream;
  Stream<bool> get open => _onOpenCtrl.stream;
  Stream<bool> get close => _onCloseCtrl.stream;
  Stream<bool> get errors => _onErrorCtrl.stream;

  void start() {
    this._socket.onOpen((e) {
      print("<----- connection opened ----->");
      _onOpenCtrl.add(true);
    });

    this._socket.onClose((e) {
      print("<----- connection closed ----->");
      if (!_closedByMe) {
        _onCloseCtrl.add(true);
      }
    });

    this._socket.onError((e) {
      print("<----- error ----->");
      _onErrorCtrl.add(true);
    });

    this._socket.onMessage((MessageEvent e) {
      Map parsedMap = JSON.decode(e.data);
      var msg = fromJSONMap(parsedMap);
      print("[ON MESSAGE] Message: " + e.data.toString());
      //send(e.data);
      _onMsgCtrl.add(msg);
    });
  }

  void closeClient() {
    _closedByMe = true;
    _socket.close();
  }

  void sendTextMessage(String text, String channel) {
    var msg = new TextMsg(sessionId, "x", text, channel);
    send(msg.toJSON());
  }

  void sendCreateChannelMessage(String channelName) {
    var msg = new ChannelAddedMsg(sessionId, channelName);
    send(msg.toJSON());
  }

  void sendUserLeftChannelMessage(String channelName) {
    var msg = new UserLeftChannelMsg(sessionId, channelName);
    send(msg.toJSON());
  }

  void sendJoinChannelMessage(String channelName) {
    var msg = new UserJoinedChannelMsg(sessionId, channelName);
    send(msg.toJSON());
  }

  void logout() {
    print("logout sent");
    var msg = new LogoutMsg(this.sessionId);
    send(msg.toJSON());
  }

  void send(String msg) {
    print("[SEND] JSON: " + msg);
    this._socket.send(msg);
  }
}
