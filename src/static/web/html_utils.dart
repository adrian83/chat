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

  MyElement withoutClass(String name) {
    _elem.classes.remove(name);
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

  MyElement hide() {
    _elem.style.display = "none";
    return this;
  }

  MyElement show() {
    _elem.style.display = "block";
    return this;
  }

  MyElement remove() {
    _elem.remove();
    return this;
  }

  Element get() {
    return _elem;
  }
}

MyElement button() => new MyElement(new ButtonElement());

MyElement textInput() => new MyElement(new InputElement(type: "text"));

MyElement link() => new MyElement(new Element.tag('a'));

MyElement span() => new MyElement(new Element.tag('span'));

MyElement div() => new MyElement(new Element.tag('div'));

MyElement p() => new MyElement(new Element.tag('p'));

MyElement li() => new MyElement(new LIElement());

MyElement brake() => new MyElement(new Element.tag('br'));

MyElement findOne(String id) => new MyElement(querySelector(id));


Function handleEnter(Function handler) {
  realHandler(KeyboardEvent event) {
    KeyEvent keyEvent = new KeyEvent.wrap(event);
    if (keyEvent.keyCode == KeyCode.ENTER) {
      handler(event);
    }
  }
  return realHandler;
}
