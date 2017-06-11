import "package:test/test.dart";

import "messages.dart";

var messageParser = new MessageParser();

void main() {
  test("parser should parse ", () {
    // given
    var senderId = "test_client";
    var msgContent = "This is error msg";
    var str = "{\"msgType\":\"ERROR\",\"senderId\":\"$senderId\",\"content\":\"$msgContent\"}";

    // when
    ErrorMsg msg = messageParser.decode(str);

    // then
    expect(msg is ErrorMsg, isTrue);
    expect(msg.senderId, equals(senderId));
    expect(msg.content, equals(msgContent));
  });
}
