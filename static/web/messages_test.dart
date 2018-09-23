import "package:test/test.dart";
import 'dart:convert';

import "package:mockito/mockito.dart";

import "messages.dart";


class JsonEncoderMock extends Mock implements JsonEncoder {}
class JsonDecoderMock extends Mock implements JsonDecoder {}

var jsonEncoderMock = new JsonEncoderMock();
var jsonDecoderMock = new JsonDecoderMock();

var messageParser = new MessageParser(new JsonEncoder(), new JsonDecoder());
var messageParserWithMocks = new MessageParser(jsonEncoderMock, jsonDecoderMock);

var senderId = "test_client";
var msgContent = "This is error msg";
var errorMsgType = "ERROR";

void main() {

  test("parser properly create ErrorMsg", () {
    // given
    var jsonMap = new Map<String, Object>();
    jsonMap["senderId"] = senderId;
    jsonMap["msgType"] = errorMsgType;
    jsonMap["content"] = msgContent;

    when(jsonDecoderMock.convert("{}")).thenReturn(jsonMap);

    // when
    ErrorMsg msg = messageParserWithMocks.decode("{}");

    // then
    expect(msg is ErrorMsg, isTrue);
    expect(msg.senderId, equals(senderId));
    expect(msg.content, equals(msgContent));
  });

  test("parser should parse JSON string and properly create ErrorMsg", () {
    // given
    var str = "{\"msgType\":\"ERROR\",\"senderId\":\"$senderId\",\"content\":\"$msgContent\"}";

    // when
    ErrorMsg msg = messageParser.decode(str);

    // then
    expect(msg is ErrorMsg, isTrue);
    expect(msg.senderId, equals(senderId));
    expect(msg.content, equals(msgContent));
  });
}
