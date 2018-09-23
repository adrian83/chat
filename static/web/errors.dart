import 'dart:html';

import 'package:logging/logging.dart';

import 'messages.dart';
import 'utils.dart';
import 'html_utils.dart';

class ErrorsPanel implements MessageConsumer {
  final Logger logger = new Logger('ErrorsPanel');

  int _id = 0;

  void onMessage(Message msg) {
    if (msg is ErrorMsg) {
      var errorElem = _create(this._id, msg.content);
      findOne("#errors-list").withChild(errorElem);
      _id += 1;
    }
  }

  Element _create(int id, String text) {
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

  String _closeErrorScript(int id) {
    return "var element = document.getElementById('error-$id'); element.parentNode.removeChild(element);";
  }
}
