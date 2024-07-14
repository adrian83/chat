
String removeWhitespace(String text) {
  return text.replaceAll(new RegExp(r'\s'), '_');
}

class Pair<T,P> {
  T fst;
  P snd;

  Pair(this.fst, this.snd);
}

Pair<String, String> strPair(String fst, String snd) {
  return new Pair<String, String>(fst, snd);
}
