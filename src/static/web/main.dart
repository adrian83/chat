import 'dart:html';
import 'dart:convert';

import 'ws_client.dart';
import 'channel_tabs.dart';
import 'channel_list.dart';
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
  var channelManager = new ChannelsManager(sessionId);
  var channelList = new ChannelList();
  var errorsPanel = new ErrorsPanel();
  var messageFactory = new MessageFactory(sessionId);


  client.start();


  void switchTab(String channelName) {
    if (channelManager.channelExists(channelName)) {
      channelManager.setVisible(channelName);
    } else {
      var msg = messageFactory.newJoinChannelMessage(channelName);
      client.send(msg);
    }
  }

  channelList.selectedChannel.listen((name) => switchTab(name));
  channelList.createdChannel.listen((name) => client.send(messageFactory.newCreateChannelMessage(name)));

  var msgConsumers = [channelManager, channelList, errorsPanel, client];

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


  channelManager.closedTabs.listen((name) => client.send(messageFactory.newUserLeftChannelMessage(name)));
  channelManager.loggedOut.listen((b) => client.send(messageFactory.newLogoutMessage()));
  channelManager.messages.listen((msg) => client.send(messageFactory.newTextMessage(msg.text, msg.channel)));
}
