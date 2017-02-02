import 'dart:html';
import 'dart:convert';

import 'messages.dart';
import 'listeners.dart';

class WSClient {
  String sessionId;
  WebSocket _socket;

  List<MessageListener> _msgListeners = new List<MessageListener>();

  WSClient(this.sessionId, this._socket) {}

  void addMessageListener(MessageListener listener) {
    _msgListeners.add(listener);
  }

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
      print("[ON MESSAGE] Message: " + msg.toString());
      for (var listener in _msgListeners) {
        listener.onMessage(msg);
      }
    });
  }

  void send(Message msg) {
    var jsonStr = msg.toJSON();
    print("[SEND] JSON: " + jsonStr);
    this._socket.send(jsonStr);
  }
}
