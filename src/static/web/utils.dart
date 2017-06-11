
String removeWhitespace(String text) {
  return text.replaceAll(new RegExp(r'\s'), '_');
}

class Pair<T,P> {
  T _first;
  P _second;

  Pair(this._first, this._second);

  T get fst => _first;
  P get snd  => _second;
}

Pair<String, String> strPair(String fst, String snd) {
  return new Pair<String, String>(fst, snd);
}
