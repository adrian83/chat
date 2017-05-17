import 'dart:html';

import 'messages.dart';
import 'html_utils.dart';

class ErrorsPanel {
  int _id = 0;

  void onMessage(Message msg) {
    if (msg is ErrorMsg) {
      var errorsList = querySelector("#errors-list");

      var spanE2 = new Element.tag('span');
      spanE2.text = msg.content;

      var spanE = new Element.tag('span');
      spanE.attributes["aria-hidden"] = "true";
      spanE.text = "x";

      var buttonE = new Element.tag('button');
      buttonE.classes.add("close");
      buttonE.attributes["data-dismiss"] = "alert";
      buttonE.attributes["aria-label"] = "Close";
      buttonE.attributes["onclick"] = _closeErrorScript(_id);


      buttonE.children.add(spanE);

      var divE = new Element.tag('div');
      divE.id = "error-" + _id.toString();
      //divE.classes.add("alert");
      //divE.classes.add("alert-danger");
      //divE.classes.add("alert-dismissible");
      withClasses(divE, const["alert", "alert-danger", "alert-dismissible"]);
      divE.attributes["role"] = "alert";
      divE.children.add(buttonE);
      divE.children.add(spanE2);

      errorsList.children.add(divE);

      _id += 1;
    }
  }

  String _closeErrorScript(int id) {
    return "var element = document.getElementById(" +
    "'error-" + id.toString() + "');" +
    "element.parentNode.removeChild(element);";
  }

}
