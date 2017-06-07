import 'dart:html';

import 'messages.dart';
import 'utils.dart';
import 'html_utils.dart';

class ErrorsPanel {
  int _id = 0;

  void onMessage(Message msg) {
    if (msg is ErrorMsg) {
      var errorElem = create(this._id, msg.content);
      findOne("#errors-list").withChild(errorElem);
      _id += 1;
    }
  }

  Element create(int id, String text) {
    var spanText = span().withText(text).get();

    var spanClose = span()
        .withText("x")
        .withAttributes([strPair("aria-hidden", "true")]).get();

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
