import 'dart:html';

import 'messages.dart';
import 'utils.dart';
import 'html_utils.dart';

class ErrorElement {
  int _id;
  String _text;

  ErrorElement(this._id, this._text);

  Element create() {
    var spanText = span()
    .withText(this._text)
    .create();

    var spanClose = span()
    .withText("x")
    .withAttributes([strPair("aria-hidden", "true")])
    .create();

    var buttonClose = button()
    .withClass("close")
    .withAttributes([
      strPair("data-dismiss", "alert"),
      strPair("aria-label", "Close"),
      strPair("onclick", _closeErrorScript())
    ])
    .withChild(spanClose)
    .create();

    return div()
    .withId("error-" + this._id.toString())
    .withClasses(["alert", "alert-danger", "alert-dismissible"])
    .withAttributes([strPair("role", "alert")])
    .withChildren([buttonClose, spanText])
    .create();
  }

  String _closeErrorScript() {
    return "var element = document.getElementById(" +
    "'error-" + this._id.toString() + "');" +
    "element.parentNode.removeChild(element);";
  }

}

class ErrorsPanel {
  int _id = 0;

  void onMessage(Message msg) {
    if (msg is ErrorMsg) {
      var errorElem = new ErrorElement(this._id, msg.content).create();
      findOne("#errors-list").withChild(errorElem);

      _id += 1;
    }
  }

}
