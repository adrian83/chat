import 'dart:html';

import 'messages.dart';
import 'listeners.dart';
import 'html_utils.dart';
import 'utils.dart';
import 'dart:async';

class NewChannelList extends MessageListener {
  StreamController _onCreatedController = new StreamController.broadcast();
  StreamController _onSelectedController = new StreamController.broadcast();

  Stream<String> get selectedChannel => _onSelectedController.stream;
  Stream<String> get createdChannel => _onCreatedController.stream;

  void createChannel() {
    _onCreatedController.add(getChannelName());
  }

  void show() {
    querySelector("#ch-create")
        .addEventListener("click", (event) => createChannel());

    querySelector("#ch-name")
        .onKeyUp
        .listen(handleEnter((event) => createChannel()));
  }

  createLinkElement(String channel) {
    var link = createLink2(
        "ch-list-name-" + removeWhitespace(channel),
        "#",
        channel,
        const ["list-group-item"],
        (e) => _onSelectedController.add(channel));
    return link;
  }

  String getChannelName() {
    InputElement element = querySelector("#ch-name");
    var name = element.value;
    element.value = "";
    return name;
  }

  void onMessage(Message msg) {
    if (msg is ChannelAddedMsg) {
      querySelector("#ch-list").children.add(createLinkElement(msg.channel));
    }

    if (msg is ChannelsListMsg) {
      var list = querySelector("#ch-list");
      msg.channels.forEach((ch) => list.children.add(createLinkElement(ch)));
    }

    if (msg is ChannelRemovedMsg) {
      querySelector("#ch-list-name-" + removeWhitespace(msg.channel)).remove();
    }
  }
}
