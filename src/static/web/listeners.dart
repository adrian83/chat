import 'messages.dart';

abstract class MessageListener {
  void onMessage(Message msg);
}

abstract class ChannelTabClosedListener {
  void onClose(String name);
}

abstract class MessageSentListener {
  void onSent(String channel, String text);
}

abstract class ChannelSelectedListener {
  void onSelect(String channel);
}

abstract class ChannelCreatedListener {
  void onCreated(String channel);
}

abstract class LogoutListener {
  void onLogout();
}
