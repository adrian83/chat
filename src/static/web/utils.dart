

String removeWhitespace(String text) {
  return text.replaceAll(new RegExp(r'\s'), '_');
}
