import 'dart:html';

import 'ws_client.dart';
import 'channel_tabs.dart';
import 'channel_list.dart';
import 'messages.dart';

main() {



  InputElement element = querySelector("#session-id");
  var sessionId = element.value;

  var wssocket = new WebSocket("ws://localhost:7070/talk");

  var client = new WSClient(sessionId, wssocket);
  var channelManager = new ChannelsManager(sessionId);
  var channelList = new NewChannelList();

  client.start();

  void sendTextMessage(TMessage tMsg) {
    var msg = new TextMsg(sessionId, "x", tMsg.text, tMsg.channel);
    client.send(msg.toJSON());
  }

  client.messages.listen((msg) => print("Data: " + msg.toString()));
  client.messages.listen((msg) => channelManager.onMessage(msg));
  client.messages.listen((msg) => channelList.onMessage(msg));

  channelManager.closedTabs.listen((name) => print("Tab closed: " + name));
  channelManager.loggedOut.listen((b) => print("User logged out"));
  channelManager.messages.listen((msg) => sendTextMessage(msg));


}
