import 'dart:html';
import 'dart:async';

import 'package:logging/logging.dart';

import 'messages.dart';
import 'html_utils.dart';
import 'utils.dart';

class RoomList implements MessageConsumer {
  final Logger logger = new Logger('RoomList');

  StreamController _onCreatedController = new StreamController<String>();
  StreamController _onSelectedController = new StreamController<String>();

  Stream<String> get selectedRoom => _onSelectedController.stream;
  Stream<String> get createdRoom => _onCreatedController.stream;

  RoomList() {
    findOne("#ch-create").withOnClickListener((event) => _createRoom());
    findOne("#ch-name")..withOnKeyPressListener(handleEnter((event) => _createRoom()));
    logger.info("RoomList created");
  }

  void onMessage(Message msg) {
    if (msg is RoomAddedMsg) {
      _handleRoomAddedMsg(msg);
    } else if (msg is RoomsListMsg) {
      _handleRoomsListMsg(msg);
    } else if (msg is RoomRemovedMsg) {
      _handleRoomRemovedMsg(msg);
    }
  }

  void _handleRoomAddedMsg(RoomAddedMsg msg) {
    findOne("#ch-list").withChild(_createLinkElement(msg.room));
  }

  void _handleRoomsListMsg(RoomsListMsg msg) {
    findOne("#ch-list").withChildren(msg.rooms.map((ch) => _createLinkElement(ch)).toList());
  }

  void _handleRoomRemovedMsg(RoomRemovedMsg msg) {
    findOne("#ch-list-name-" + removeWhitespace(msg.room)).remove();
  }

  void _createRoom() => _onCreatedController.add(_getRoomName());

  Element _createLinkElement(String room) {
    return link()
        .withId("ch-list-name-" + removeWhitespace(room))
        .withText(room)
        .withHref("#")
        .withOnClickListener((e) => _onSelectedController.add(room))
        .withClass("list-group-item")
        .get();
  }

  String _getRoomName() {
    InputElement element = findOne("#ch-name").get();
    var name = element.value;
    element.value = "";
    return name;
  }
}
