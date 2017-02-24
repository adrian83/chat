import 'dart:html';

import 'messages.dart';

class ErrorsPanel {
  int _id = 0;

  void onMessage(Message msg) {
    print("[ErrorsPanel] onMessage: " + msg.toString());
    if (msg is ErrorMsg) {
      var l = querySelector("#errors-list");

      var spanE = new Element.tag('span');
      spanE.attributes["aria-hidden"] = "true";
      spanE.text = "x";

      var buttonE = new Element.tag('button');
      buttonE.classes.add("close");
      buttonE.attributes["data-dismiss"] = "alert";
      buttonE.attributes["aria-label"] = "Close";
      buttonE.attributes["onclick"] =
          "var element = document.getElementById('error-" +
              _id.toString() +
              "');element.parentNode.removeChild(element);";

      buttonE.children.add(spanE);

      var divE = new Element.tag('div');
      divE.id = "error-" + _id.toString();
      divE.classes.add("alert");
      divE.classes.add("alert-danger");
      divE.classes.add("alert-dismissible");
      divE.attributes["role"] = "alert";

      divE.children.add(buttonE);

      var spanE2 = new Element.tag('span');
      spanE2.text = msg.content;

      divE.children.add(spanE2);
      //alert alert-danger alert-dismissible

      l.children.add(divE);

      _id += 1;
    }
  }
}
