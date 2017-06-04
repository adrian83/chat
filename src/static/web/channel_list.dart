import 'dart:html';

import 'messages.dart';
import 'html_utils.dart';
import 'utils.dart';
import 'dart:async';

class ChannelList {
  StreamController _onCreatedController = new StreamController.broadcast();
  StreamController _onSelectedController = new StreamController.broadcast();

  Stream<String> get selectedChannel => _onSelectedController.stream;
  Stream<String> get createdChannel => _onCreatedController.stream;

  ChannelList() {
    querySelector("#ch-create")
        .addEventListener("click", (event) => createChannel());

    querySelector("#ch-name")
        .onKeyUp
        .listen(handleEnter((event) => createChannel()));
  }

  void createChannel() {
    _onCreatedController.add(getChannelName());
  }

  createLinkElement(String channel) {
    return link()
    .withId("ch-list-name-" + removeWhitespace(channel))
    .withText(channel)
    .withHref("#")
    .withOnClickListener((e) => _onSelectedController.add(channel))
    .withClass("list-group-item")
    .create();

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
    } else if (msg is ChannelsListMsg) {
      var list = querySelector("#ch-list");
      msg.channels.forEach((ch) => list.children.add(createLinkElement(ch)));
    } else if (msg is ChannelRemovedMsg) {
      querySelector("#ch-list-name-" + removeWhitespace(msg.channel)).remove();
    }
  }
}
