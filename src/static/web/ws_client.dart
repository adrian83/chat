import 'dart:html';
import 'dart:async';

import 'package:logging/logging.dart';

import 'messages.dart';


class WSClient {
  final Logger logger = new Logger('ChannelList');

  MessageParser _parser = new MessageParser();
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
    this._socket.onOpen.listen((e) => _onOpenCtrl.add(true));

    this._socket.onClose.listen((e){
      if (!_closedByMe) {
         _onCloseCtrl.add(true);
       }
    });

    this._socket.onError.listen((e) => _onErrorCtrl.add(true));

    this._socket.onMessage.listen((e){
      var msg = _parser.parse(e.data.toString());
      print("[ON MESSAGE] Message: " + e.data.toString());
      //send(e.data);
      _onMsgCtrl.add(msg);
    });

  }

  void closeClient() {
    _closedByMe = true;
    this._socket.close();
  }

  void sendTextMessage(String text, String channel) {
    var msg = new TextMsg(sessionId, sessionId, text, channel);
    logger.info("Sending: $msg");
    sendMsg(msg);
  }

  void sendCreateChannelMessage(String channelName) {
    var msg = new ChannelAddedMsg(sessionId, channelName);
    sendMsg(msg);
  }

  void sendUserLeftChannelMessage(String channelName) {
    var msg = new UserLeftChannelMsg(sessionId, channelName);
    sendMsg(msg);
  }

  void sendJoinChannelMessage(String channelName) {
    var msg = new UserJoinedChannelMsg(sessionId, channelName);
    sendMsg(msg);
  }

  void logout() {
    var msg = new LogoutMsg(this.sessionId);
    sendMsg(msg);
  }

  void sendMsg(Message msg) {
    var str = _parser.stringify(msg);
    logger.info("Sending stringified msg: $str");
    this._socket.send(str);
  }

  void send(String msg) {
    print("[SEND] JSON: " + msg);
    this._socket.send(msg);
  }
}
