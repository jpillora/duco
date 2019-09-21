function handler() {
  EVENT_DATA=$1
  echo "$EVENT_DATA" 1>&2
  RESPONSE="Echo: '$EVENT_DATA'"
  echo $RESPONSE
}
