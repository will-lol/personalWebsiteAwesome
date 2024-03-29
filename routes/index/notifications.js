import { registration } from '/assets/js/lib/components/sw-register.js'

document.getElementById("subscribe").addEventListener('click', async (e) => {
  const subscription = await registration.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: await fetch("api/notifications/publickey").then((res) => res.text())
  });
  await fetch("api/notifications/subscribe", {
    method: "POST",
    body: JSON.stringify(subscription.toJSON()),
  });
})
