import polka from 'polka';
import { parse } from 'url';

// import sveltekit handler from build output
import { handler } from './build/handler.js';

const app = polka();

const PORT = process.env.PORT || 3000;

app.use(async (req, res) => {
  const parsedUrl = parse(req.url!, true);
  const { pathname } = parsedUrl;
  console.log('Request:', pathname);

  function handleNextRequest() {
    console.log('Handling next request');
  }

  const result = handler(req, res, handleNextRequest);
  console.log('Result:', result);

  res.send(`Hello, world!`);
});

app.listen(PORT, (err: Error) => {
  if (err) {
    console.error('Error starting server:', err);
  } else {
    console.log(`Server running on http://localhost:${PORT}`);
  }
});