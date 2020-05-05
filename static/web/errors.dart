import 'dart:html';

import 'package:logging/logging.dart';
import 'package:uuid/uuid.dart';

import 'messages.dart';
import 'utils.dart';
import 'html_utils.dart';

class ErrorsPanel implements MessageConsumer {
  final Logger logger = new Logger('ErrorsPanel');

  void onMessage(Message msg) {
    if (msg is ErrorMsg) {
      var errorElem = _create(Uuid().v4(), msg.content);
      findOne("#errors-list").withChild(errorElem);
    }
  }

  Element _create(String id, String text) {
    logger.info("Displaing error with id '$id' and text '$text'");

    var spanText = span().withText(text).get();

    var spanClose = span()
        .withText("x")
        .withAttributes([strPair("aria-hidden", "true")])
        .get();

    var buttonClose = button()
        .withClass("close")
        .withAttributes([
          strPair("data-dismiss", "alert"),
          strPair("aria-label", "Close"),
          strPair("onclick", _closeErrorScript(id))
        ])
        .withChild(spanClose)
        .get();

    return div()
        .withId("error-$id")
        .withClasses([
          "alert",
          "alert-danger",
          "alert-dismissible"
        ])
        .withAttributes([strPair("role", "alert")])
        .withChildren([buttonClose, spanText])
        .get();
  }

  String _closeErrorScript(String id) {
    return "var element = document.getElementById('error-$id'); element.parentNode.removeChild(element);";
  }
}
