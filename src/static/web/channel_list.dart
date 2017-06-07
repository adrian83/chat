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
    findOne("#ch-create").withOnClickListener((event) => createChannel());
    findOne("#ch-name").withOnKeyPressListener(handleEnter((event) => createChannel()));
  }

  void createChannel() => _onCreatedController.add(getChannelName());

  Element createLinkElement(String channel) {
    return link()
        .withId("ch-list-name-" + removeWhitespace(channel))
        .withText(channel)
        .withHref("#")
        .withOnClickListener((e) => _onSelectedController.add(channel))
        .withClass("list-group-item")
        .get();
  }

  String getChannelName() {
    InputElement element = findOne("#ch-name").get();
    var name = element.value;
    element.value = "";
    return name;
  }

  void onMessage(Message msg) {
    if (msg is ChannelAddedMsg) {
      findOne("#ch-list").withChild(createLinkElement(msg.channel));
    } else if (msg is ChannelsListMsg) {
      findOne("#ch-list").withChildren(msg.channels.map((ch) => createLinkElement(ch)));
    } else if (msg is ChannelRemovedMsg) {
      findOne("#ch-list-name-" + removeWhitespace(msg.channel)).get().remove();
    }
  }
}
