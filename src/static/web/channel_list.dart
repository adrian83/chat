import 'dart:html';
import 'dart:async';

import 'package:logging/logging.dart';

import 'messages.dart';
import 'html_utils.dart';
import 'utils.dart';

class ChannelList implements MessageConsumer {
  final Logger logger = new Logger('ChannelList');

  StreamController _onCreatedController = new StreamController.broadcast();
  StreamController _onSelectedController = new StreamController.broadcast();

  Stream<String> get selectedChannel => _onSelectedController.stream;
  Stream<String> get createdChannel => _onCreatedController.stream;

  ChannelList() {
    findOne("#ch-create").withOnClickListener((event) => _createChannel());
    findOne("#ch-name")..withOnKeyPressListener(handleEnter((event) => _createChannel()));
    logger.info("Created");
  }

  void onMessage(Message msg) {
    if (msg is ChannelAddedMsg) {
      _handleChannelAddedMsg(msg);
    } else if (msg is ChannelsListMsg) {
      _handleChannelsListMsg(msg);
    } else if (msg is ChannelRemovedMsg) {
      _handleChannelRemovedMsg(msg);
    }
  }

  void _handleChannelAddedMsg(ChannelAddedMsg msg) {
    findOne("#ch-list").withChild(_createLinkElement(msg.channel));
  }

  void _handleChannelsListMsg(ChannelsListMsg msg) {
    findOne("#ch-list").withChildren(msg.channels.map((ch) => _createLinkElement(ch)));
  }

  void _handleChannelRemovedMsg(ChannelRemovedMsg msg) {
    findOne("#ch-list-name-" + removeWhitespace(msg.channel)).remove();
  }

  void _createChannel() => _onCreatedController.add(_getChannelName());

  Element _createLinkElement(String channel) {
    return link()
        .withId("ch-list-name-" + removeWhitespace(channel))
        .withText(channel)
        .withHref("#")
        .withOnClickListener((e) => _onSelectedController.add(channel))
        .withClass("list-group-item")
        .get();
  }

  String _getChannelName() {
    InputElement element = findOne("#ch-name").get();
    var name = element.value;
    element.value = "";
    return name;
  }
}
