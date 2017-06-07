import 'dart:async';
import 'dart:html';

import 'messages.dart';
import 'html_utils.dart';
import 'utils.dart';

const String MAIN = "main";
const String EXIT = "exit";
const String LOGOUT = "logout";

class ChannelsManager {
  String _clientId;
  Map<String, ChannelTab> _channels = new Map<String, ChannelTab>();

  StreamController _onTabClosedCtrl = new StreamController.broadcast();
  StreamController _onLoggedOutCtrl = new StreamController.broadcast();
  StreamController _onMessageCtrl = new StreamController.broadcast();

  Stream<String> get closedTabs => _onTabClosedCtrl.stream;
  Stream<bool> get loggedOut => _onLoggedOutCtrl.stream;
  Stream<BaseMessage> get messages => _onMessageCtrl.stream;

  ChannelsManager(this._clientId);

  void newMessage(BaseMessage msg) => _onMessageCtrl.add(msg);
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
    print("visible "+name);
    var n = _channels.containsKey(name) ? name : MAIN;
    _channels[n].setVisible();
  }

  bool channelExists(String name) => _channels.containsKey(name);

  void onMessage(Message msg) {
    if (msg is ChannelsListMsg) {
      msg.channels.forEach((name){
        addChannel(name);
      });
      //addChannel(MAIN);
      setVisible(MAIN);
    } else if (_shouldDisplayMessage(msg)) {
      _channels[msg.channel]._displayMessage(msg.senderName, msg.content);
    } else if (_shouldAddChannel(msg)) {
      addChannel(msg.channel);
      setVisible(msg.channel);
    }
  }

  bool _shouldDisplayMessage(Message msg) => msg is TextMsg && channelExists(msg.channel);

  bool _shouldAddChannel(Message msg) => msg is UserJoinedChannelMsg || (msg is ChannelAddedMsg && msg.senderId == _clientId);
}

class BaseMessage {
  String _channel, _text;
  BaseMessage(this._channel, this._text);
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

    var tabTitleLink = link()
        .withHref("#")
        .withText(_name)
        .withOnClickListener((e) => setVisible())
        .get();

    var tabListElem = li()
        .withId("ch-$_escapedName")
        .withAttributes([strPair("role", "presentation")])
        .withChild(tabTitleLink)
        .get();

    findOne("#ch-tabs").withChild(tabListElem);


    var sendMsgButton = button()
        .withId("msg-send-$_escapedName")
        .withText("Send")
        .withOnClickListener(_onSent)
        .withClasses(["btn", "btn-default"]).get();

    var sendMegSpan = span()
        .withClass("input-group-btn")
        .withChild(sendMsgButton)
        .get();

    var msgTextInput = textInput()
        .withId("msg-content-$_escapedName")
        .withOnKeyPressListener(handleEnter(_onSent))
        .withClass("form-control")
        .get();

    var inputGroupDiv = div()
        .withClass("input-group")
        .withChildren([msgTextInput, sendMegSpan])
        .get();

    var conversationDiv = div()
        .withId("conversation-$_escapedName")
        .withStyle((style){
          style.maxHeight = "400px";
          style.overflowY = "scroll";
        })
        .get();

    var contentDiv = div()
        .withId("content-$_escapedName")
        .withChildren([
          brake().get(),
          inputGroupDiv,
          brake().get(),
          conversationDiv])
        .withStyle((style){
          style.display = "none";
        })
        .get();

    findOne("#ch-contents").withChild(contentDiv);
  }

  void _onSent(e) {
    var text = _getMessageText();
    if (text == EXIT) {
      if (_name != MAIN) {
        _manager.tabClosed(_name);
        _manager.setVisible(MAIN);
        close();
      }
    } else if (text == LOGOUT) {
      _manager.userLoggedOut();
    } else {
      var tMsg = new BaseMessage(_name, text);
      _manager.newMessage(tMsg);
    }
  }

  void setVisible() {
    findOne("#ch-tabs").forEachChild((e) => e.withoutClass("active"));
    findOne("#ch-$_escapedName").withClass("active");
    findOne("#ch-contents").forEachChild((ch){
      ch.withStyle((style){
        style.display = "none";
      });
    });
    findOne("#content-$_escapedName").show();
    findOne("#msg-content-$_escapedName").get().focus();
  }

  void close() {
    findOne("#ch-$_escapedName").remove();
    findOne("#content-$_escapedName").remove();
  }

  void _displayMessage(String author, String text) {
    var textParagraph = p().withText("$author: $text").get();
    findOne("#conversation-$_escapedName").withChildAtIndex(0, textParagraph);
  }

  String _getMessageText() {
    InputElement element = findOne("#msg-content-$_escapedName").get();
    var text = element.value;
    element.value = "";
    return text;
  }
}
