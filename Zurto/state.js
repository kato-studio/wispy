"use strict";
// References/Inspirations:
// - https://github.com/nanostores/nanostores#vanilla-js
// - https://gist.github.com/developit/a0430c500f5559b715c2dddf9c40948d

export function state(initialValue = null) {
  let _value = initialValue;
  let _subscribers = [];

  const _state = {
    set: (new_value) => {
      _value = new_value;
      _subscribers.forEach((cb) => cb(_value));
    },
    subscribe: (cb) => {
      _subscribers.push(cb);
      cb(_value);
      return () => {
        _subscribers = _subscribers.filter((s) => s !== cb);
      };
    },
    unsubscribe: (cb) => {
      return () => {
        _subscribers = _subscribers.filter((s) => s !== cb);
      };
    },
  };

  return _state;
}
