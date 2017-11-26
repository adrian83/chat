import 'dart:html';
import 'dart:convert';

import 'ws_client.dart';
import 'room_tabs.dart';
import 'room_list.dart';
import 'messages.dart';
import 'errors.dart';
import 'html_utils.dart';

import 'package:logging/logging.dart';


main() {

  Logger.root.level = Level.ALL;
  Logger.root.onRecord.listen((LogRecord rec) {
    print('${rec.level.name}: ${rec.time}: ${rec.message}');
  });

  Logger logger = new Logger('Main');
  logger.info("Application is starting...");

  InputElement element = querySelector("#session-id");
  var sessionId = element.value;

  var host = window.location.hostname + (window.location.port != null ? ':' + window.location.port : '');
  var wssocket = new WebSocket("ws://$host/talk");

  var messageParser = new MessageParser(new JsonEncoder(), new JsonDecoder());
  var client = new WSClient(sessionId, wssocket, messageParser);
  var roomManager = new RoomsManager(sessionId);
  var roomList = new RoomList();
  var errorsPanel = new ErrorsPanel();
  var messageFactory = new MessageFactory(sessionId);


  client.start();


  void switchTab(String roomName) {
    if (roomManager.roomExists(roomName)) {
      roomManager.setVisible(roomName);
    } else {
      var msg = messageFactory.newJoinRoomMessage(roomName);
      client.send(msg);
    }
  }

  roomList.selectedRoom.listen((name) => switchTab(name));
  roomList.createdRoom.listen((name) => client.send(messageFactory.newCreateRoomMessage(name)));

  var msgConsumers = [roomManager, roomList, errorsPanel, client];

  client.messages.listen((msg){
    msgConsumers.forEach((c){
      c.onMessage(msg);
    });
  });

  void onSocketClose() {
    var errMsg = new ErrorMsg("system", "Connection has been closed. Maybe server is down.");
    errorsPanel.onMessage(errMsg);
  }

  client.close.listen((b) => onSocketClose());
  client.open.listen((b){
    findOne("#connection-info").hide();
    findOne("#panels").show();
    findOne("#logout-info").show();
  });
  client.errors.listen((b) => onSocketClose());


  roomManager.closedTabs.listen((name) => client.send(messageFactory.newUserLeftRoomMessage(name)));
  roomManager.loggedOut.listen((b) => client.send(messageFactory.newLogoutMessage()));
  roomManager.messages.listen((msg) => client.send(messageFactory.newTextMessage(msg.text, msg.room)));
}
