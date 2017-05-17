import "package:test/test.dart";
import "package:mockito/mockito.dart";

import 'dart:convert';
import 'dart:async';

class WebSocketMock extends Mock implements WebSocketWrapper {}

/*
void main() {

  var webSocketMock = new WebSocketMock();
  var sessionId = "abc-def-ghi";

  var wsClient = new WSClient(sessionId, webSocketMock);

  test("send(msg) should send given string message", () {
    // given
    var msg = "Hello Dart";

    // when
    wsClient.send(msg);

    // then
    verify(webSocketMock.send(msg));
  });

}
*/
