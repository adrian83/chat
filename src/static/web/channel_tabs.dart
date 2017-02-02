import 'dart:html';

import 'messages.dart';
import 'listeners.dart';
import 'html_utils.dart';
import 'utils.dart';

const String MAIN = "main";

class ChannelsManager extends MessageListener {
  String _clientId;
  Map<String, NewChannel> _channels = new Map<String, NewChannel>();
  List<ChannelTabClosedListener> _onTabClosedListener =
      new List<ChannelTabClosedListener>();
  List<LogoutListener> _onLogoutListeners = new List<LogoutListener>();
  List<MessageSentListener> _onSentListeners = new List<MessageSentListener>();

  ChannelsManager(this._clientId);

  void addChannelTabClosedListener(ChannelTabClosedListener listener) {
    _onTabClosedListener.add(listener);
    _channels.forEach((n, ch) => ch.addChannelTabClosedListener(listener));
  }

  void addMessageSentListener(MessageSentListener listener) {
    _onSentListeners.add(listener);
    _channels.forEach((n, ch) => ch.addMessageSentListener(listener));
  }

  void addLogoutListener(LogoutListener listener) {
    _onLogoutListeners.add(listener);
    _channels.forEach((n, ch) => ch.addLogoutListener(listener));
  }

  void addChannel(String name) {
    var channel = new NewChannel(name);
    channel.show();

    _onTabClosedListener.forEach((l) => channel.addChannelTabClosedListener(l));
    _onSentListeners.forEach((l) => channel.addMessageSentListener(l));
    _onLogoutListeners.forEach((l) => channel.addLogoutListener(l));

    _channels[name] = channel;
  }

  void setVisible(String name) {
    var n = _channels.containsKey(name) ? name : MAIN;
    _channels[n].setVisible();
  }

  bool channelExists(String name) {
    return _channels.containsKey(name);
  }

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

class NewChannel {
  List<ChannelTabClosedListener> _onCloseListeners = new List<ChannelTabClosedListener>();
  List<LogoutListener> _onLogoutListeners = new List<LogoutListener>();
  List<MessageSentListener> _onSentListeners = new List<MessageSentListener>();

  String _name;

  NewChannel(this._name);

  void addChannelTabClosedListener(ChannelTabClosedListener listener) {
    this._onCloseListeners.add(listener);
  }

  void addMessageSentListener(MessageSentListener listener) {
    this._onSentListeners.add(listener);
  }

  void addLogoutListener(LogoutListener listener) {
    this._onLogoutListeners.add(listener);
  }

  void show() {
    var attrs = new Map<String, String>();
    attrs['role'] = 'presentation';
    var liTab = createListElem("ch-" + removeWhitespace(_name), attrs);

    void onSelect(e) {
      setVisible();
    }

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
          _onCloseListeners.forEach((listener) => listener.onClose(_name));
          close();
        }
      } else if (text == "logout") {
        _onLogoutListeners.forEach((listener) => listener.onLogout());
      } else {
        _onSentListeners.forEach((listener) => listener.onSent(_name, text));
      }
    }

    var sendMsgButton = createButton(
        "msg-send-" + removeWhitespace(_name), "Send", const ["btn", "btn-default"], onSent);

    var sp = new Element.tag('span');
    sp.classes.add("input-group-btn");
    sp.children.add(sendMsgButton);

    var textInput = createTextInput(
        "msg-content-" + removeWhitespace(_name), const ["form-control"], handleEnter(onSent));

    var inputGroupDiv = new Element.tag('div');
    inputGroupDiv.classes.add("input-group");
    inputGroupDiv.children.add(textInput);
    inputGroupDiv.children.add(sp);

    var conversationDiv = new Element.tag('div');
    conversationDiv.id = "conversation-" + removeWhitespace(this._name);
    conversationDiv.style.maxHeight = "400px";
    conversationDiv.style.overflowY = "scroll";

    var contentDiv = new Element.tag('div');
    contentDiv.children.add(new Element.tag('br'));
    contentDiv.children.add(inputGroupDiv);
    contentDiv.children.add(new Element.tag('br'));
    contentDiv.children.add(conversationDiv);
    contentDiv.id = "content-" + removeWhitespace(this._name);

    contentDiv.style.display = "none";

    contents().add(contentDiv);
  }

  void setVisible() {
    tabs().forEach((tab) => tab.classes.remove("active"));
    querySelector("#ch-" + removeWhitespace(_name)).classes.add("active");
    contents().forEach((div) => div.style.display = "none");
    querySelector("#content-" + removeWhitespace(_name)).style.display = "block";
    querySelector("#msg-content-" + removeWhitespace(_name)).focus();
  }

  void close() {
    print("close " + removeWhitespace(_name));
    querySelector("#ch-" + removeWhitespace(_name)).remove();
    querySelector("#content-" + removeWhitespace(_name)).remove();
  }

  void displayMessage(String author, String text) {
    var textP = new Element.tag('p');
    textP.text = author + ": " + text;
    var elem = querySelector("#conversation-" + removeWhitespace(this._name));
    elem.children.insert(0, textP);
  }

  List<Element> tabs() {
    return querySelector("#ch-tabs").children;
  }

  List<Element> contents() {
    return querySelector("#ch-contents").children;
  }

  String getMessageText(){
    InputElement element = querySelector("#msg-content-" + removeWhitespace(_name));
    var text = element.value;
    element.value = "";
    return text;
  }
}
