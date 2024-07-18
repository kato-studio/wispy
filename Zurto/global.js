"use strict";

export function global() {
  // stores contains key value states
  // > ./state.js
  let _store = {};

  return {
    init: function () {
      _store = {
        __config: {
          //
        },
      };
    },
    set: function (key, value) {
      _store[key] = value;
    },
    get: function (key) {
      return _store[key];
    },
    remove: function (key) {
      delete _store[key];
    },
  };
}
