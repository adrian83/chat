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
  String msgType;
  String senderId;

  Message(this.msgType, this.senderId);

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map["msgType"] = msgType;
    map["senderId"] = senderId;
    return map;
  }
}

class TextMsg extends Message {
  String senderName;
  String content;
  String channel;

  TextMsg(String senderId, this.senderName, this.content, this.channel)
      : super(TEXT_MSG, senderId);

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["senderName"] = senderName;
    map["content"] = content;
    map["channel"] = channel;
    return map;
  }
}

class LogoutMsg extends Message {
  LogoutMsg(String senderId) : super(LOGOUT_MSG, senderId);
}

class UnknownOpMsg extends Message {
  UnknownOpMsg() : super(UNKNOWN, UNKNOWN);
}

class UserLeftChannelMsg extends Message {
  String channel;

  UserLeftChannelMsg(String senderId, this.channel)
      : super(USER_LEFT_CHANNEL, senderId);

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["channel"] = channel;
    return map;
  }
}

class UserJoinedChannelMsg extends Message {
  String channel;

  UserJoinedChannelMsg(String senderId, this.channel)
      : super(USER_JOINED_CHANNEL, senderId);

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["channel"] = channel;
    return map;
  }
}

class ChannelRemovedMsg extends Message {
  String channel;

  ChannelRemovedMsg(String senderId, this.channel)
      : super(DEL_CHANNEL, senderId);

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["channel"] = channel;
    return map;
  }
}

class ChannelAddedMsg extends Message {
  String channel;

  ChannelAddedMsg(String senderId, this.channel) : super(ADD_CHANNEL, senderId);

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["channel"] = channel;
    return map;
  }
}

class ErrorMsg extends Message {
  String senderName;
  String content;

  ErrorMsg(String senderId, this.senderName, this.content)
      : super(ERROR_MSG, senderId);

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["senderName"] = senderName;
    map["content"] = content;
    return map;
  }
}

class ChannelsListMsg extends Message {
  String senderName;
  List<String> channels;

  ChannelsListMsg(String senderId, this.senderName, this.channels)
      : super(CHANNELS_LIST, senderId);

  Map<String, Object> toJson() {
    var map = new Map<String, Object>();
    map.addAll(super.toJson());
    map["senderName"] = senderName;
    map["channels"] = channels;
    return map;
  }
}
