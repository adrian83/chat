import 'dart:html';

import 'messages.dart';
import 'channel_tabs.dart';
import 'listeners.dart';
import 'ws_client.dart';
import 'channel_list.dart';
import 'errors.dart';

class RealChannelCreatedListener extends ChannelCreatedListener {
  WSClient _client;

  RealChannelCreatedListener(this._client) {}

  void onCreated(String channel) {
    var msg = new ChannelAddedMsg(_client.sessionId, channel);
    _client.send(msg);
  }
}

class RealChannelSelectedListener extends ChannelSelectedListener {
  WSClient _client;
  ChannelsManager _channelsManager;

  RealChannelSelectedListener(this._client, this._channelsManager);

  void onSelect(String channel) {
    if (!_channelsManager.channelExists(channel)) {
      var msg = new UserJoinedChannelMsg(_client.sessionId, channel);
      _client.send(msg);
    } else {
      _channelsManager.setVisible(channel);
    }
  }
}

class RealChannelTabClosedListener extends ChannelTabClosedListener {
  WSClient _client;
  ChannelsManager _channelsManager;

  RealChannelTabClosedListener(this._client, this._channelsManager);

  void onClose(String name) {
    var msg = new UserLeftChannelMsg(_client.sessionId, name);
    _client.send(msg);
    _channelsManager.setVisible("main");
  }
}

class RealMessageSentListener extends MessageSentListener {
  WSClient _client;

  RealMessageSentListener(this._client);

  void onSent(String channel, String text) {
    var msg = new TextMsg(_client.sessionId, _client.sessionId, text, channel);
    _client.send(msg);
  }
}

class RealLogoutListener extends LogoutListener {
  WSClient _client;

  RealLogoutListener(this._client);

  void onLogout() {
    var msg = new LogoutMsg(_client.sessionId);
    _client.send(msg);
  }
}

class LogoutMessageListener extends MessageListener {
  void onMessage(Message msg) {
    if (msg is LogoutMsg) {
      // cleanup
      window.location.assign('/logout');
    }
  }
}

main() {
  InputElement element = querySelector("#session-id");
  var sessionId = element.value;

  var channelsManager = new ChannelsManager(sessionId);

  var channelList = new NewChannelList();
  channelList.show();

  var errorsPanel = new ErrorsPanel();

  var logoutMsgListener = new LogoutMessageListener();

  var wssocket = new WebSocket("ws://localhost:7070/talk");

  var client = new WSClient(sessionId, wssocket);
  client.addMessageListener(channelList);
  client.addMessageListener(channelsManager);
  client.addMessageListener(errorsPanel);
  client.addMessageListener(logoutMsgListener);

  var realChannelCreatedListener = new RealChannelCreatedListener(client);
  channelList.addChannelCreatedListener(realChannelCreatedListener);

  var realChannelSelectedListener =
      new RealChannelSelectedListener(client, channelsManager);
  channelList.addChannelSelectedListener(realChannelSelectedListener);

  var realChannelTabClosedListener =
      new RealChannelTabClosedListener(client, channelsManager);
  channelsManager.addChannelTabClosedListener(realChannelTabClosedListener);

  var realMessageSentListener = new RealMessageSentListener(client);
  channelsManager.addMessageSentListener(realMessageSentListener);

  var realLogoutListener = new RealLogoutListener(client);
  channelsManager.addLogoutListener(realLogoutListener);

  client.start();
}
