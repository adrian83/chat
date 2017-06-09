import 'dart:async';
import 'dart:html';

import 'package:logging/logging.dart';

import 'messages.dart';
import 'html_utils.dart';
import 'utils.dart';

const String MAIN = "main";
const String EXIT = "exit";
const String LOGOUT = "logout";

class BaseMessage {
  String _channel, _text;
  BaseMessage(this._channel, this._text);
  String get channel => _channel;
  String get text => _text;
}

class ChannelsManager implements MessageConsumer {
  final Logger logger = new Logger('ChannelsManager');

  String _clientId;
  Map<String, ChannelTab> _channels = new Map<String, ChannelTab>();

  StreamController _onTabClosedCtrl = new StreamController.broadcast();
  StreamController _onLoggedOutCtrl = new StreamController.broadcast();
  StreamController _onMessageCtrl = new StreamController.broadcast();

  ChannelsManager(this._clientId);

  Stream<String> get closedTabs => _onTabClosedCtrl.stream;
  Stream<bool> get loggedOut => _onLoggedOutCtrl.stream;
  Stream<BaseMessage> get messages => _onMessageCtrl.stream;

  void newMessage(BaseMessage msg) => _onMessageCtrl.add(msg);
  void userLoggedOut() => _onLoggedOutCtrl.add(true);
  void tabClosed(String name) {
    _onTabClosedCtrl.add(name);
    _channels.remove(name);
    logger.info("Tab with name '$name' should be closed");
  }

  void hideAll() {
    findOne("#ch-tabs").forEachChild((e) => e.withoutClass("active"));
    findOne("#ch-contents").forEachChild((ch) {
      ch.withStyle((style) {
        style.display = "none";
      });
    });
  }

  void addChannel(String name) {
    if (!channelExists(name)) {
      var channel = new ChannelTab(name, this);
      channel.show();
      _channels[name] = channel;
      logger.info("Tab with name '$name' should be added");
    }
  }

  void setVisible(String name) {
    var n = _channels.containsKey(name) ? name : MAIN;
    _channels[n].setVisible();
  }

  bool channelExists(String name) => _channels.containsKey(name);

  void onMessage(Message msg) {
    if (msg is ChannelsListMsg) {
      _handleChannelsListMsg(msg);
    } else if (msg is UserJoinedChannelMsg) {
      _handleUserJoinedChannelMsg(msg);
    } else if (_shouldDisplayMessage(msg)) {
      _handleTextMsg(msg);
    } else if (_shouldAddChannel(msg)) {
      _handleChannelAddedMsg(msg);
    }
  }

  void _handleChannelsListMsg(ChannelsListMsg msg) {
    _addAndShowChannel(MAIN);
    msg.channels.forEach((name) => addChannel(name));
  }

  void _handleTextMsg(TextMsg msg) {
    _channels[msg.channel]._displayMessage(msg.senderName, msg.content);
  }

  void _handleUserJoinedChannelMsg(UserJoinedChannelMsg msg) {
    _addAndShowChannel(msg.channel);
  }

  void _handleChannelAddedMsg(ChannelAddedMsg msg) {
    _addAndShowChannel(msg.channel);
  }

  void _addAndShowChannel(String name) {
    addChannel(name);
    setVisible(name);
  }

  bool _shouldDisplayMessage(Message msg) =>
      msg is TextMsg && channelExists(msg.channel);

  bool _shouldAddChannel(Message msg) =>
      msg is ChannelAddedMsg && msg.senderId == _clientId;
}

class ChannelTab {
  final Logger logger = new Logger('ChannelTab');

  ChannelsManager _manager;
  String _name;
  String _escapedName;

  ChannelTab(this._name, this._manager) {
    this._escapedName = removeWhitespace(_name);
    logger.info("Created tab with name '$_name'");
  }

  void setVisible() {
    _manager.hideAll();
    findOne("#ch-$_escapedName").withClass("active");
    findOne("#content-$_escapedName").show();
    findOne("#msg-content-$_escapedName").get().focus();
    logger.info("Tab with name '$_name' should be visible");
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
        .withClasses(["btn", "btn-default"])
        .get();

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
        .withStyle((style) {
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
          conversationDiv
        ])
        .withStyle((style) {
          style.display = "none";
        })
        .get();

    findOne("#ch-contents").withChild(contentDiv);
  }

  void _onSent(e) {
    var text = _getMessageText();
    logger.info("Sending text '$text' from tab with name '$_name'");
    if (text == EXIT) {
      if (_name != MAIN) {
        _manager.tabClosed(_name);
        _manager.setVisible(MAIN);
        _close();
        logger.info("Exiting...");
      }
    } else if (text == LOGOUT) {
      _manager.userLoggedOut();
      logger.info("Logging out...");
    } else {
      _manager.newMessage(new BaseMessage(_name, text));
    }
  }

  void _close() {
    findOne("#ch-$_escapedName").remove();
    findOne("#content-$_escapedName").remove();
  }

  void _displayMessage(String author, String text) {
    var textParagraph = paragraph().withText("$author: $text").get();
    findOne("#conversation-$_escapedName").withChildAtIndex(0, textParagraph);
  }

  String _getMessageText() {
    InputElement element = findOne("#msg-content-$_escapedName").get();
    var text = element.value;
    element.value = "";
    return text;
  }
}
