import 'dart:async';
import 'dart:html';

import 'messages.dart';
import 'html_utils.dart';
import 'utils.dart';

const String MAIN = "main";

class ChannelsManager {
  String _clientId;
  Map<String, ChannelTab> _channels = new Map<String, ChannelTab>();

  StreamController _onTabClosedCtrl = new StreamController.broadcast();
  StreamController _onLoggedOutCtrl = new StreamController.broadcast();
  StreamController _onMessageCtrl = new StreamController.broadcast();

  Stream<String> get closedTabs => _onTabClosedCtrl.stream;
  Stream<bool> get loggedOut => _onLoggedOutCtrl.stream;
  Stream<TMessage> get messages => _onMessageCtrl.stream;

  ChannelsManager(this._clientId);

  void newMessage(TMessage msg) => _onMessageCtrl.add(msg);
  void userLoggedOut() => _onLoggedOutCtrl.add(true);
  void tabClosed(String name) {
    _onTabClosedCtrl.add(name);
    _channels.remove(name);
  }

  void addChannel(String name) {
    var channel = new ChannelTab(name, this);
    channel.show();
    _channels[name] = channel;
  }

  void setVisible(String name) {
    var n = _channels.containsKey(name) ? name : MAIN;
    _channels[n].setVisible();
  }

  bool channelExists(String name) => _channels.containsKey(name);

  void onMessage(Message msg) {
    if (msg is ChannelsListMsg) {
      addChannel(MAIN);
      setVisible(MAIN);
    }
    if (msg is TextMsg && channelExists(msg.channel)) {
      _channels[msg.channel].displayMessage(msg.senderName, msg.content);
    }
    if (msg is ChannelAddedMsg && msg.senderId == _clientId) {
      addChannel(msg.channel);
      setVisible(msg.channel);
    }
    if (msg is UserJoinedChannelMsg) {
      addChannel(msg.channel);
      setVisible(msg.channel);
    }
  }
}

class TMessage {
  String _channel, _text;
  TMessage(this._channel, this._text);
  String get channel => _channel;
  String get text => _text;
}

class ChannelTab {
  ChannelsManager _manager;
  String _name;
  String _escapedName;

  ChannelTab(this._name, this._manager) {
    this._escapedName = removeWhitespace(_name);
  }

  void show() {


    var tabTitle = link()
        .withHref("#")
        .withText(_name)
        .withOnClickListener((e) => setVisible())
        .create();

    var liTab = li()
        .withId("ch-" + _escapedName)
        .withAttributes([strPair("role", "presentation")])
        .withChild(tabTitle).create();


    tabs().add(liTab);

    void onSent(e) {
      var text = getMessageText();
      if (text == "exit") {
        if (_name != MAIN) {
          _manager.tabClosed(_name);
          _manager.setVisible(MAIN);
          close();
        }
      } else if (text == "logout") {
        _manager.userLoggedOut();
      } else {
        var tMsg = new TMessage(_name, text);
        _manager.newMessage(tMsg);
      }
    }

    var sendMsgButton = button()
        .withId("msg-send-" + _escapedName)
        .withText("Send")
        .withOnClickListener(onSent)
        .withClasses(["btn", "btn-default"]).create();

    var sp = span()
        .withClass("input-group-btn")
        .withChild(sendMsgButton)
        .create();

    var textIn = textInput()
        .withId("msg-content-" + _escapedName)
        .withOnKeyPressListener(handleEnter(onSent))
        .withClass("form-control")
        .create();

    var inputGroupDiv =
        div().withClass("input-group").withChildren([textIn, sp]).create();


    var conversationDiv = div()
        .withId("conversation-" + _escapedName)
        .withStyle((style){
          style.maxHeight = "400px";
          style.overflowY = "scroll";
        })
        .create();

    var contentDiv = div()
        .withId("content-" + _escapedName)
        .withChildren([brake(), inputGroupDiv, brake(), conversationDiv])
        .withStyle((style){
          style.display = "none";
        })
        .create();

    findOne("#ch-contents").withChild(contentDiv);
  }

  void setVisible() {
    tabs().forEach((tab) => tab.classes.remove("active"));
    querySelector("#ch-" + _escapedName).classes.add("active");
    findOne("#ch-contents").forEachChild((ch){
      ch.withStyle((style){
        style.display = "none";
      });
    });

    querySelector("#content-" + _escapedName).style.display = "block";
    findOne("#msg-content-" + _escapedName).create().focus();
  }

  void close() {
    print("close " + _escapedName);
    querySelector("#ch-" + _escapedName).remove();
    querySelector("#content-" + _escapedName).remove();
  }

  void displayMessage(String author, String text) {
    var textP = p().withText(author + ": " + text).create();
    findOne("#conversation-" + _escapedName).withChildAtIndex(0, textP);
  }

  List<Element> tabs() => querySelector("#ch-tabs").children;


  String getMessageText() {
    InputElement element = findOne("#msg-content-" + _escapedName).create();
    var text = element.value;
    element.value = "";
    return text;
  }
}
