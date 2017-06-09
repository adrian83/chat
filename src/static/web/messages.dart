import 'dart:convert';

import 'package:logging/logging.dart';

const String UNKNOWN = "UNKNOWN";
const String TEXT_MSG = "TEXT_MSG";
const String DEL_CHANNEL = "REM_CH";
const String ADD_CHANNEL = "ADD_CH";
const String CHANNELS_LIST = "CHAN_LIST_MSG";
const String USER_LEFT_CHANNEL = "USER_LEFT_CH";
const String USER_JOINED_CHANNEL = "USER_JOINED_CH";
const String ERROR_MSG = "ERROR";
const String LOGOUT_MSG = "LOGOUT_USER";

abstract class MessageConsumer {
  void onMessage(Message msg);
}

class MessageParser {
  final Logger logger = new Logger('ChannelList');

  final JsonEncoder encoder = new JsonEncoder();
  final JsonDecoder decoder = new JsonDecoder();

  String stringify(Message msg) {
    return encoder.convert(msg);
  }

  Message parse(String jsonStr) {
    logger.info("Parsing message: $jsonStr");

    Map json = decoder.convert(jsonStr);

    var msgType = json["msgType"];
    var senderId = json["senderId"];
    var senderName = json["senderName"];
    var channel = json["channel"];
    var channels = json["channels"];
    var content = json["content"];

    switch (msgType) {
      case DEL_CHANNEL:
        return new ChannelRemovedMsg(senderId, channel);
      case ADD_CHANNEL:
        return new ChannelAddedMsg(senderId, channel);
      case TEXT_MSG:
        return new TextMsg(senderId, senderName, content, channel);
      case ERROR_MSG:
        return new ErrorMsg(senderId, senderName, content);
      case CHANNELS_LIST:
        return new ChannelsListMsg(senderId, senderName, channels);
      case USER_JOINED_CHANNEL:
        return new UserJoinedChannelMsg(senderId, channel);
      case LOGOUT_MSG:
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

class ChannelMessage extends Message {
  String _channel;

  ChannelMessage(String senderId, String msgType, this._channel)
      : super(msgType, senderId);

  String get channel => _channel;

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["channel"] = _channel;
    return map;
  }
}

class TextMsg extends ChannelMessage {
  String _senderName;
  String _content;


  TextMsg(String senderId, String channel, this._senderName, this._content)
      : super(senderId, TEXT_MSG, channel);

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
  LogoutMsg(String senderId) : super(LOGOUT_MSG, senderId);
}

class UnknownOpMsg extends Message {
  UnknownOpMsg() : super(UNKNOWN, UNKNOWN);
}

class UserLeftChannelMsg extends ChannelMessage {

  UserLeftChannelMsg(String senderId, String channel)
      : super(senderId, USER_LEFT_CHANNEL, channel);
}

class UserJoinedChannelMsg extends ChannelMessage {

  UserJoinedChannelMsg(String senderId, String channel)
      : super(senderId, USER_JOINED_CHANNEL, channel);
}

class ChannelRemovedMsg extends ChannelMessage {

  ChannelRemovedMsg(String senderId, String channel)
      : super(senderId, DEL_CHANNEL, channel);
}

class ChannelAddedMsg extends ChannelMessage {

  ChannelAddedMsg(String senderId, String channel)
      : super(senderId, ADD_CHANNEL, channel);
}

class ErrorMsg extends Message {
  String _senderName;
  String _content;

  ErrorMsg(String senderId, this._senderName, this._content)
      : super(ERROR_MSG, senderId);

  String get content => _content;

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["senderName"] = _senderName;
    map["content"] = _content;
    return map;
  }
}

class ChannelsListMsg extends Message {
  String _senderName;
  List<String> _channels;

  ChannelsListMsg(String senderId, this._senderName, this._channels)
      : super(CHANNELS_LIST, senderId);

  List<String> get channels => _channels;

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["senderName"] = _senderName;
    map["channels"] = _channels;
    return map;
  }
}
