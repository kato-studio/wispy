import polka from 'polka';
import { parse } from 'url';

export default function adapter(options = {}) {
  const { port = 3000 } = options;

  return {
    name: 'polka-adapter',

    async adapt(builder) {
      console.log('Adapting to Polka');
      console.log('builder', builder);

      return "polka goes boop";
    }
  };
}
