self.addEventListener('push', event => {
  console.log('[Service Worker] Push Received.');
  console.log(`[Service Worker] Push had this data: "${event.data.text()}"`);

  const title = 'will.forsale';
  const options = {
    body: event.data.text(),
  };

  event.waitUntil(self.registration.showNotification(title, options));
});
