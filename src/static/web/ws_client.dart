import 'dart:html';
import 'dart:convert';
import 'dart:async';

import 'messages.dart';

class WSClient {
  String sessionId;
  WebSocket _socket;
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
    this._socket.onOpen.listen((e) {
      print("<----- connection opened ----->");
      _onOpenCtrl.add(true);
    });

    this._socket.onClose.listen((e) {
      print("<----- connection closed ----->");
      if(!_closedByMe){
        _onCloseCtrl.add(true);
      }
    });

    this._socket.onError.listen((e) {
      print("<----- error ----->");
      _onErrorCtrl.add(true);
    });

    this._socket.onMessage.listen((MessageEvent e) {
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
    print("join");
    var msg = new ChannelAddedMsg(sessionId, channelName);
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
