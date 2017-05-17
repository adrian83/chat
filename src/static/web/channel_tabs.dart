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
    var attrs = new Map<String, String>();
    attrs['role'] = 'presentation';
    var liTab = createListElem("ch-" + _escapedName, attrs);

    void onSelect(e) => setVisible();

    //var tabTitle = createLink("#", _name, onSelect);

    var tabTitle = new Element.tag('a');
    tabTitle.href = "#";
    tabTitle.text = _name;
    tabTitle.addEventListener("click", onSelect);

    liTab.children.add(tabTitle);

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

    var sendMsgButton = createButton("msg-send-" + _escapedName, "Send",
        const ["btn", "btn-default"], onSent);

    var sp = new Element.tag('span');
    sp.classes.add("input-group-btn");
    sp.children.add(sendMsgButton);

    var textInput = createTextInput("msg-content-" + _escapedName,
        const ["form-control"], handleEnter(onSent));

    var inputGroupDiv = new Element.tag('div');
    inputGroupDiv.classes.add("input-group");
    inputGroupDiv.children.add(textInput);
    inputGroupDiv.children.add(sp);

    var conversationDiv = new Element.tag('div');
    conversationDiv.id = "conversation-" + _escapedName;
    conversationDiv.style.maxHeight = "400px";
    conversationDiv.style.overflowY = "scroll";

    var contentDiv = new Element.tag('div');
    contentDiv.children.add(new Element.tag('br'));
    contentDiv.children.add(inputGroupDiv);
    contentDiv.children.add(new Element.tag('br'));
    contentDiv.children.add(conversationDiv);
    contentDiv.id = "content-" + _escapedName;

    contentDiv.style.display = "none";

    contents().add(contentDiv);
  }

  void setVisible() {
    tabs().forEach((tab) => tab.classes.remove("active"));
    querySelector("#ch-" + _escapedName).classes.add("active");
    contents().forEach((div) => div.style.display = "none");
    querySelector("#content-" + _escapedName).style.display = "block";
    querySelector("#msg-content-" + _escapedName).focus();
  }

  void close() {
    print("close " + _escapedName);
    querySelector("#ch-" + _escapedName).remove();
    querySelector("#content-" + _escapedName).remove();
  }

  void displayMessage(String author, String text) {
    var textP = new Element.tag('p');
    textP.text = author + ": " + text;
    var elem = querySelector("#conversation-" + _escapedName);
    elem.children.insert(0, textP);
  }

  List<Element> tabs() => querySelector("#ch-tabs").children;

  List<Element> contents() => querySelector("#ch-contents").children;

  String getMessageText() {
    InputElement element = querySelector("#msg-content-" + _escapedName);
    var text = element.value;
    element.value = "";
    return text;
  }
}
