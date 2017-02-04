import 'dart:html';
import 'dart:convert';
import 'dart:async';

import 'messages.dart';

class WSClient {
  String sessionId;
  WebSocket _socket;
  StreamController _onExitCtrl = new StreamController.broadcast();

  WSClient(this.sessionId, this._socket);

  Stream<Message> get messages => _onExitCtrl.stream;

  void start() {
    this._socket.onOpen.listen((e) {
      print("<----- connection opened ----->");
    });

    this._socket.onClose.listen((e) {
      print("<----- connection closed ----->");
    });

    this._socket.onError.listen((e) {
      print("<----- error ----->");
    });

    this._socket.onMessage.listen((MessageEvent e) {
      Map parsedMap = JSON.decode(e.data);
      var msg = fromJSONMap(parsedMap);
      print("[ON MESSAGE] Message: " + e.data.toString());
      //send(e.data);
      _onExitCtrl.add(msg);
    });
  }

  void send(String msg) {
    print("[SEND] JSON: " + msg);
    this._socket.send(msg);
  }
}
