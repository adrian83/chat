import 'dart:html';

import 'ws_client.dart';
import 'channel_tabs.dart';
import 'channel_list.dart';
import 'messages.dart';
import 'errors.dart';
import 'html_utils.dart';

main() {
  InputElement element = querySelector("#session-id");
  var sessionId = element.value;

  var host = window.location.hostname +
      (window.location.port != null ? ':' + window.location.port : '');
  var wssocket = new WebSocket("ws://" + host + "/talk");

  var client = new WSClient(sessionId, wssocket);
  var channelManager = new ChannelsManager(sessionId);
  var channelList = new NewChannelList();
  var errorsPanel = new ErrorsPanel();

  client.start();

  void onMsg(Message msg) {
    if (msg is LogoutMsg) {
      client.closeClient();
      window.location.assign('/?reason=logout');
    }
  }

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
  client.open.listen((b) => hideElement("#connection-info"));
  client.errors.listen((b) => onSocketClose());

  channelManager.closedTabs.listen((name) => print("Tab closed: " + name));
  channelManager.loggedOut.listen((b) => client.logout());
  channelManager.messages
      .listen((msg) => client.sendTextMessage(msg.text, msg.channel));
}
