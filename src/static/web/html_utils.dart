import 'dart:html';

import 'utils.dart';

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

  MyElement withChild(Element child) {
    _elem.children.add(child);
    return this;
  }

  MyElement withChildAtIndex(int index, Element child) {
    _elem.children.insert(index, child);
    return this;
  }

  MyElement withAttributes(List<Pair<String,String>> attrs) {
    attrs.forEach((p) => _elem.attributes[p.fst] = p.snd);
    return this;
  }

  MyElement withChildren(List<Element> children) {
    _elem.children.addAll(children);
    return this;
  }

  MyElement forEachChild(Function f) {
    _elem.children.forEach((e) => f(new MyElement(e)));
    return this;
  }

  MyElement withStyle(Function setStylesFunction) {
    setStylesFunction(_elem.style);
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
  return new MyElement(new Element.tag('a'));
}

MyElement span() {
  return new MyElement(new Element.tag('span'));
}

MyElement div() {
  return new MyElement(new Element.tag('div'));
}

MyElement p() {
  return new MyElement(new Element.tag('p'));
}

MyElement li() {
  return new MyElement(new LIElement());
}

MyElement findOne(String id) {
  return new MyElement(querySelector(id));
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

Element brake() => new Element.tag('br');
