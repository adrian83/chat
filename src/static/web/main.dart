import 'dart:html';

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

  var host = window.location.hostname +
      (window.location.port != null ? ':' + window.location.port : '');
  var wssocket = new WebSocket("ws://" + host + "/talk");

  var client = new WSClient(sessionId, wssocket);
  var channelManager = new ChannelsManager(sessionId);
  var channelList = new ChannelList();
  var errorsPanel = new ErrorsPanel();

  client.start();

  void onMsg(Message msg) {
    if (msg is LogoutMsg) {
      client.closeClient();
      window.location.assign('/logout');
    } else if (msg is UserJoinedChannelMsg) {
      if (channelManager.channelExists(msg.channel)) {
        channelManager.setVisible(msg.channel);
      } else {
        channelManager.addChannel(msg.channel);
        channelManager.setVisible(msg.channel);
      }
    }
  }

  void switchTab(String name) {
    if (channelManager.channelExists(name)) {
      channelManager.setVisible(name);
    } else {
      client.sendJoinChannelMessage(name);
    }
  }

  channelList.selectedChannel.listen((name) => switchTab(name));
  channelList.createdChannel
      .listen((name) => client.sendCreateChannelMessage(name));

  client.messages.listen((msg) => print("Data: " + msg.toString()));
  client.messages.listen((msg) => channelManager.onMessage(msg));
  client.messages.listen((msg) => channelList.onMessage(msg));
  client.messages.listen((msg) => onMsg(msg));

  void onSocketClose() {
    //hideElement("#panels");
    //hideElement("#logout-info");
    var errMsg = new ErrorMsg("system", "system",
        "Connection has been closed. Maybe server is down.");
    errorsPanel.onMessage(errMsg);
  }

  client.close.listen((b) => onSocketClose());
  client.open.listen((b){
    findOne("#connection-info").hide();
    findOne("#panels").show();
    findOne("#logout-info").show();
  });
  client.errors.listen((b) => onSocketClose());

  void onTabClosed(String channelName) {
    client.sendUserLeftChannelMessage(channelName);
  }

  channelManager.closedTabs.listen((name) => onTabClosed(name));
  channelManager.loggedOut.listen((b) => client.logout());
  channelManager.messages
      .listen((msg) => client.sendTextMessage(msg.text, msg.channel));
}
