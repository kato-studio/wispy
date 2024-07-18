"use strict";

export function manyListeners(cb, states) {
  _states = states;

  states.forEach((state) => {
    state.subscribe(cb);
  });

  return {
    unsubscribe: () => {
      _states.forEach((state) => {
        state.unsubscribe(cb);
      });
    },
    add: (state) => {
      _states.push(state);
      state.subscribe(cb);
    },
    remove: (state) => {
      _states = _states.filter((s) => s !== state);
      state.unsubscribe(cb);
    },
  };
}
