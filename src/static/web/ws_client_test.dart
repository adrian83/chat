import "package:test/test.dart";
import "package:mockito/mockito.dart";

import 'dart:html';

import 'ws_client.dart';
import 'messages.dart';

class WebSocketMock extends Mock implements WebSocket {}


void main() {

  var messageFactory = new MessageFactory("123-456");

  var webSocketMock = new WebSocketMock();
  var sessionId = "abc-def-ghi";

  var wsClient = new WSClient(sessionId, webSocketMock);

  test("send(msg) should send given string message", () {
    // given
    var msg = messageFactory.newLogoutMessage();

    // when
    wsClient.send(msg);

    // then
    verify(webSocketMock.send(msg));
  });

}
