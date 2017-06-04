import 'dart:html';


class MyElement {
  Element _elem;

  MyElement(this._elem);

  MyElement withId(String id) {
    _elem.id = id;
    return this;
  }

  MyElement withText(String text) {
    _elem.text = text;
    return this;
  }

  MyElement withHref(String href) {
    _elem.href = href;
    return this;
  }

  MyElement withOnClickListener(Function onClickListener) {
    _elem.addEventListener("click", onClickListener);
    return this;
  }

  MyElement withOnKeyPressListener(Function onKeyPressListener) {
    _elem.onKeyPress.listen(onKeyPressListener);
    return this;
  }

  MyElement withClass(String name) {
    _elem.classes.add(name);
    return this;
  }

  MyElement withClasses(List<String> names) {
    _elem.classes.addAll(names);
    return this;
  }

  Element create() {
    return _elem;
  }
}

MyElement button() {
  return new MyElement(new ButtonElement());
}

MyElement textInput() {
  return new MyElement(new InputElement(type: "text"));
}

MyElement link() {
  return new MyElement(new LinkElement());
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

void withClasses(Element elem, List<String> classes) {
  classes.forEach((cls) => elem.classes.add(cls));
}

Element brake() => new Element.tag('br');
