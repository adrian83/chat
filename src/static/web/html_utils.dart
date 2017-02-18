import 'dart:html';

ButtonElement createButton(
    String id, String text, List<String> cssClasses, Function onClickListener) {
  var button = new ButtonElement();
  button.id = id;
  button.text = text;
  button.addEventListener("click", onClickListener);

  cssClasses.forEach((cssClass) => button.classes.add(cssClass));
  return button;
}

InputElement createTextInput(String id, List<String> cssClasses, Function onKeyPressListener) {
  var textInput = new InputElement(type: "text");
  cssClasses.forEach((cssClass) => textInput.classes.add(cssClass));
  textInput.id = id;
  textInput.onKeyPress.listen(onKeyPressListener);
  return textInput;
}

LinkElement createLink(String href, String text, Function onClickListener) {
  var link = new LinkElement();
  link.href = href;
  link.text = text;
  link.addEventListener("click", onClickListener);
  return link;
}

LinkElement createLink2(String id, String href, String text, List<String> cssClasses, Function onClickListener) {
  var link = new LinkElement();
  link.id = id;
  link.href = href;
  link.text = text;
  link.addEventListener("click", onClickListener);
  cssClasses.forEach((cssClass) => link.classes.add(cssClass));
  return link;
}

LIElement createListElem(String id, Map<String, String> attrs) {
  var listElem = new LIElement();
  listElem.id = id;
  attrs.forEach((attrName, attrValue) => listElem.attributes[attrName] = attrValue);
  return listElem;
}

Function handleEnter(Function handler) {
  realHandler(KeyboardEvent event) {
    KeyEvent keyEvent = new KeyEvent.wrap(event);
    if (keyEvent.keyCode == KeyCode.ENTER) {
      handler(event);
    }
  }
  return realHandler;
}

void hideElement(String elemId) {
  querySelector(elemId).style.display = "none";
}
