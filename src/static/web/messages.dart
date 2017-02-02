const String UNKNOWN = "UNKNOWN";
const String TEXT_MSG = "TEXT_MSG";
const String DEL_CHANNEL = "REM_CH";
const String ADD_CHANNEL = "ADD_CH";
const String CHANNELS_LIST = "CHAN_LIST_MSG";
const String USER_LEFT_CHANNEL = "USER_LEFT_CH";
const String USER_JOINED_CHANNEL = "USER_JOINED_CH";
const String ERROR_MSG = "ERROR";
const String LOGOUT_MSG = "LOGOUT_USER";

Message fromJSONMap(Map json) {
  var msgType = json["msgType"];
  print("msgType: " + msgType);

  switch (msgType) {
    case DEL_CHANNEL:
      return new ChannelRemovedMsg(json["senderId"], json["channel"]);
    case ADD_CHANNEL:
      return new ChannelAddedMsg(json["senderId"], json["channel"]);
    case TEXT_MSG:
      return new TextMsg(json["senderId"], json["senderName"], json["content"],
          json["channel"]);
    case ERROR_MSG:
      return new ErrorMsg(
          json["senderId"], json["senderName"], json["content"]);
    case CHANNELS_LIST:
      return new ChannelsListMsg(json["senderId"], json["senderName"],
          json["receiver"], json["channels"]);
    case USER_JOINED_CHANNEL:
      return new UserJoinedChannelMsg(json["senderId"], json["channel"]);
    case LOGOUT_MSG:
      return new LogoutMsg(json["senderId"]);
    default:
      var msg = new UnknownOpMsg();
      print("fromJSONMap: " + msg.toString());
      return msg;
  }
}

class Message {
  String _msgType;
  String _senderId;

  String get senderId => _senderId;

  String toJSON() {
    return "{ \"msgType\":\"" +
        this._msgType +
        "\", \"senderId\":\"" +
        this._senderId +
        "\" }";
  }
}

class LogoutMsg extends Message {
  LogoutMsg(String senderId) {
    this._msgType = LOGOUT_MSG;
    this._senderId = senderId;
  }

  String toString() {
    return "LogoutMsg {_msgType: " +
        this._msgType +
        ", senderId: " +
        this._senderId +
        " }";
  }

  String toJSON() {
    return "{ \"msgType\":\"" +
        this._msgType +
        "\", \"senderId\":\"" +
        this._senderId +
        "\" }";
  }
}

class UnknownOpMsg extends Message {
  UnknownOpMsg() {
    this._msgType = UNKNOWN;
    this._senderId = UNKNOWN;
  }

  String toString() {
    return "UnknownOpMsg {_msgType: " +
        this._msgType +
        ", senderId: " +
        this._senderId +
        " }";
  }

  String toJSON() {
    return "{ \"msgType\":\"" +
        this._msgType +
        "\", \"senderId\":\"" +
        this._senderId +
        "\" }";
  }
}

class UserLeftChannelMsg extends Message {
  String channel;

  UserLeftChannelMsg(String senderId, this.channel) {
    this._msgType = USER_LEFT_CHANNEL;
    this._senderId = senderId;
  }

  String toString() {
    return "UserLeftChannelMsg {_msgType: " +
        this._msgType +
        ", senderId: " +
        this._senderId +
        ", channel: " +
        this.channel +
        " }";
  }

  String toJSON() {
    return "{ \"msgType\":\"" +
        this._msgType +
        "\", \"senderId\":\"" +
        this._senderId +
        "\", \"channel\":\"" +
        this.channel +
        "\" }";
  }
}

class UserJoinedChannelMsg extends Message {
  String channel;

  UserJoinedChannelMsg(String senderId, this.channel) {
    this._msgType = USER_JOINED_CHANNEL;
    this._senderId = senderId;
  }

  String toString() {
    return "UserJoinedChannelMsg {_msgType: " +
        this._msgType +
        ", senderId: " +
        this._senderId +
        ", channel: " +
        this.channel +
        " }";
  }

  String toJSON() {
    return "{ \"msgType\":\"" +
        this._msgType +
        "\", \"senderId\":\"" +
        this._senderId +
        "\", \"channel\":\"" +
        this.channel +
        "\" }";
  }
}

class ChannelRemovedMsg extends Message {
  String channel;

  ChannelRemovedMsg(String senderId, this.channel) {
    this._msgType = DEL_CHANNEL;
    this._senderId = senderId;
  }

  String toString() {
    return "ChannelRemovedMsg {_msgType: " +
        this._msgType +
        ", senderId: " +
        this._senderId +
        ", channel: " +
        this.channel +
        " }";
  }

  String toJSON() {
    return "{ \"msgType\":\"" +
        this._msgType +
        "\", \"senderId\":\"" +
        this._senderId +
        "\", \"channel\":\"" +
        this.channel +
        "\" }";
  }
}

class ChannelAddedMsg extends Message {
  String channel;

  ChannelAddedMsg(String senderId, this.channel) {
    this._msgType = ADD_CHANNEL;
    this._senderId = senderId;
  }

  String toString() {
    return "ChannelAddedMsg {_msgType: " +
        this._msgType +
        ", senderId: " +
        this._senderId +
        ", channel: " +
        this.channel +
        " }";
  }

  String toJSON() {
    return "{ \"msgType\":\"" +
        this._msgType +
        "\", \"senderId\":\"" +
        this._senderId +
        "\", \"channel\":\"" +
        this.channel +
        "\" }";
  }
}

class ErrorMsg extends Message {
  String senderName;
  String content;

  ErrorMsg(String senderId, this.senderName, this.content) {
    this._msgType = ERROR_MSG;
    this._senderId = senderId;
  }

  String toString() {
    return "ErrorMsg {_msgType: " +
        this._msgType +
        ", senderId: " +
        this._senderId +
        ", senderName: " +
        this.senderName +
        ", content: " +
        this.content +
        " }";
  }

  String toJSON() {
    return "{ \"msgType\":\"" +
        this._msgType +
        "\", \"senderId\":\"" +
        this._senderId +
        "\", \"senderName\":\"" +
        this.senderName +
        "\", \"content\":\"" +
        this.content +
        "\" }";
  }
}

class TextMsg extends Message {
  String senderName;
  String content;
  String channel;

  TextMsg(String senderId, this.senderName, this.content, this.channel) {
    this._msgType = TEXT_MSG;
    this._senderId = senderId;
  }

  String toString() {
    return "TextMsg {_msgType: " +
        this._msgType +
        ", senderId: " +
        this._senderId +
        ", senderName: " +
        this.senderName +
        ", content: " +
        this.content +
        ", channel: " +
        this.channel +
        " }";
  }

  String toJSON() {
    return "{ \"msgType\":\"" +
        this._msgType +
        "\", \"senderId\":\"" +
        this._senderId +
        "\", \"senderName\":\"" +
        this.senderName +
        "\", \"content\":\"" +
        this.content +
        "\", \"channel\":\"" +
        this.channel +
        "\" }";
  }
}

class ChannelsListMsg extends Message {
  String senderName;
  String receiver;
  List<String> channels;

  ChannelsListMsg(
      String senderId, this.senderName, this.receiver, this.channels) {
    this._msgType = CHANNELS_LIST;
    this._senderId = senderId;
  }

  String toString() {
    return "ChannelsListMsg {_msgType: " +
        this._msgType +
        ", senderId: " +
        this._senderId +
        ", senderName: " +
        this.senderName +
        ", receiver: " +
        this.receiver +
        ", channels: [" +
        this.channels.join(",") +
        "] }";
  }

  String toJSON() {
    return "{ \"msgType\":\"" +
        this._msgType +
        "\", \"senderId\":\"" +
        this._senderId +
        "\", \"senderName\":\"" +
        this.senderName +
        "\", \"receiver\":\"" +
        this.receiver +
        "\", \"channels\":\"" +
        this.channels.join(",") +
        "\" }";
  }
}
