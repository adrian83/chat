import 'dart:convert';

import 'package:logging/logging.dart';

const String UNKNOWN = "UNKNOWN";
const String TEXT_MSG_MT = "TEXT_MSG";
const String REMOVE_ROOM_MT = "REMOVE_ROOM";
const String CREATE_ROOM_MT = "CREATE_ROOM";
const String ROOMS_NAMES_MT = "ROOMS_LIST";
const String USER_LEFT_ROOM_MT = "USER_LEFT_ROOM";
const String USER_JOINED_ROOM_MT = "USER_JOINED_ROOM";
const String ERROR_MSG_MT = "ERROR";
const String LOGOUT_MT = "LOGOUT_USER";


abstract class MessageConsumer {
  void onMessage(Message msg);
}

class MessageFactory {
  String _senderId;

  MessageFactory(this._senderId);

  Message newTextMessage(String text, String roomName) => new TextMsg(_senderId, roomName, _senderId, text);
  Message newCreateRoomMessage(String roomName) => new RoomAddedMsg(_senderId, roomName);
  Message newUserLeftRoomMessage(String roomName) => new UserLeftRoomMsg(_senderId, roomName);
  Message newJoinRoomMessage(String roomName) => new UserJoinedRoomMsg(_senderId, roomName);
  Message newLogoutMessage() => new LogoutMsg(_senderId);
}

class MessageParser {
  final Logger logger = new Logger('MessageParser');

  final JsonEncoder _encoder;
  final JsonDecoder _decoder;

  MessageParser(this._encoder, this._decoder);

  String encode(Message msg) {
    return _encoder.convert(msg);
  }

  Message decode(String jsonStr) {
    logger.info("Parsing message: $jsonStr");

    Map json = _decoder.convert(jsonStr);

    var msgType = json["msgType"];
    var senderId = json["senderId"];
    var senderName = json["senderName"];
    var room = json["room"];
    var rooms = json["rooms"];
    var content = json["content"];

    switch (msgType) {
      case REMOVE_ROOM_MT:
        return new RoomRemovedMsg(senderId, room);
      case CREATE_ROOM_MT:
        return new RoomAddedMsg(senderId, room);
      case TEXT_MSG_MT:
        return new TextMsg(senderId, room, senderName, content);
      case ERROR_MSG_MT:
        return new ErrorMsg(senderId, content);
      case ROOMS_NAMES_MT:
        return new RoomsListMsg(senderId, senderName, rooms);
      case USER_JOINED_ROOM_MT:
        return new UserJoinedRoomMsg(senderId, room);
      case LOGOUT_MT:
        return new LogoutMsg(senderId);
      default:
        return new UnknownOpMsg();
    }
  }
}

class Message {
  String _msgType;
  String _senderId;

  Message(this._msgType, this._senderId);

  String get senderId => _senderId;

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map["msgType"] = _msgType;
    map["senderId"] = _senderId;
    return map;
  }
}

class RoomMessage extends Message {
  String _room;

  RoomMessage(String senderId, String msgType, this._room)
      : super(msgType, senderId);

  String get room => _room;

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["room"] = _room;
    return map;
  }
}

class TextMsg extends RoomMessage {
  String _senderName;
  String _content;


  TextMsg(String senderId, String room, this._senderName, this._content)
      : super(senderId, TEXT_MSG_MT, room);

  String get content => _content;
  String get senderName => _senderName;

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["senderName"] = _senderName;
    map["content"] = _content;
    return map;
  }
}

class LogoutMsg extends Message {
  LogoutMsg(String senderId) : super(LOGOUT_MT, senderId);
}

class UnknownOpMsg extends Message {
  UnknownOpMsg() : super(UNKNOWN, UNKNOWN);
}

class UserLeftRoomMsg extends RoomMessage {

  UserLeftRoomMsg(String senderId, String room)
      : super(senderId, USER_LEFT_ROOM_MT, room);
}

class UserJoinedRoomMsg extends RoomMessage {

  UserJoinedRoomMsg(String senderId, String room)
      : super(senderId, USER_JOINED_ROOM_MT, room);
}

class RoomRemovedMsg extends RoomMessage {

  RoomRemovedMsg(String senderId, String room)
      : super(senderId, REMOVE_ROOM_MT, room);
}

class RoomAddedMsg extends RoomMessage {

  RoomAddedMsg(String senderId, String room)
      : super(senderId, CREATE_ROOM_MT, room);
}

class ErrorMsg extends Message {
  String _content;

  ErrorMsg(String senderId, this._content)
      : super(ERROR_MSG_MT, senderId);

  String get content => _content;

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["content"] = _content;
    return map;
  }
}

class RoomsListMsg extends Message {
  String _senderName;
  List<String> _rooms;

  RoomsListMsg(String senderId, this._senderName, this._rooms)
      : super(ROOMS_NAMES_MT, senderId);

  List<String> get rooms => _rooms;

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["senderName"] = _senderName;
    map["rooms"] = _rooms;
    return map;
  }
}
