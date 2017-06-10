//import 'dart:html';
import 'dart:async';

import 'package:logging/logging.dart';

import 'messages.dart';
import 'html_utils.dart';


class WSClient implements MessageConsumer {
  final Logger logger = new Logger('ChannelList');

  MessageParser _parser;
  String sessionId;
  var _socket;
  bool _closedByMe;

  StreamController _onMsgCtrl = new StreamController.broadcast();
  StreamController _onOpenCtrl = new StreamController.broadcast();
  StreamController _onCloseCtrl = new StreamController.broadcast();
  StreamController _onErrorCtrl = new StreamController.broadcast();

  WSClient(this.sessionId, this._socket, this._parser);

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
      var msgStr = e.data.toString();
      logger.info("Received message string: $msgStr");
      var msg = _parser.decode(msgStr);
      _onMsgCtrl.add(msg);
    });

  }

  void onMessage(Message msg) {
    if (msg is LogoutMsg) {
      _closeClient();
      changeLocation('/logout');
    }
  }

  void _closeClient() {
    _closedByMe = true;
    this._socket.close();
  }

  void send(Message msg) {
    var str = _parser.encode(msg);
    logger.info("Sending stringified msg: $str");
    this._socket.send(str);
  }

}
