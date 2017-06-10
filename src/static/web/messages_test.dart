import "package:test/test.dart";
import "package:mockito/mockito.dart";

import "messages.dart";


var messageParser = new MessageParser();


void main() {
  test("parser should parse ", () {
    // given

    var str = "{\"msgType\":\"ERROR\",\"senderId\":\"test_client\",\"content\":\"This is error msg\"}";

    // when
    Message msg = messageParser.decode(str);

    // then
    expect(msg.senderId, equals("test_client"));
  });
}
