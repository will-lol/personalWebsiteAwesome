if ('serviceWorker' in navigator) {
	try {
		const registration = await navigator.serviceWorker.register('assets/js/sw.js');
		document.getElementById("notifications").addEventListener('click', async (e) => {
			const subscription = await registration.pushManager.subscribe({
				userVisibleOnly: true,
			});
			console.log(subscription);
		})

	} catch (error) {
		console.error(`Registration failed with ${error}`);
	}
}

